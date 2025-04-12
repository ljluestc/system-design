package cache

import (
    "crypto/md5"
    "crypto/sha256"
    "encoding/binary"
    "errors"
    "fmt"
    "sort"
    "strconv"
    "sync"
    "time"

    "distributed-cache/internal/utils"
)

// Constants for HashRing configuration
const (
    DefaultReplicas     = 3               // Default number of virtual nodes per physical node
    MinReplicas         = 1               // Minimum allowed replicas
    MaxReplicas         = 1000            // Maximum allowed replicas to prevent memory explosion
    HashAlgorithmMD5    = "md5"           // MD5 hash algorithm identifier
    HashAlgorithmSHA256 = "sha256"        // SHA-256 hash algorithm identifier
)

// ErrInvalidReplicas is returned when the number of replicas is invalid
var ErrInvalidReplicas = errors.New("number of replicas must be between MinReplicas and MaxReplicas")

// ErrNodeAlreadyExists is returned when attempting to add a duplicate node
var ErrNodeAlreadyExists = errors.New("node already exists in the hash ring")

// ErrEmptyRing is returned when operations are attempted on an empty hash ring
var ErrEmptyRing = errors.New("hash ring is empty")

// HashRing represents a consistent hashing ring for distributing keys across nodes
type HashRing struct {
    replicas    int                  // Number of virtual nodes per physical node
    nodes       []string             // List of physical node identifiers
    ring        map[uint32]string    // Mapping of hash values to node identifiers
    sorted      []uint32             // Sorted list of hash values for binary search
    mu          sync.RWMutex         // Mutex for thread-safe operations
    hashAlgo    string               // Hash algorithm to use (md5 or sha256)
    initialized time.Time            // Timestamp of ring initialization
}

// NewHashRing creates a new HashRing instance with the specified number of replicas
// and a default hash algorithm (MD5). It ensures the replicas count is valid.
// 
// Parameters:
// - replicas: Number of virtual nodes per physical node
// 
// Returns:
// - *HashRing: A pointer to the initialized HashRing
// - error: An error if the replicas count is invalid
func NewHashRing(replicas int) (*HashRing, error) {
    if replicas < MinReplicas || replicas > MaxReplicas {
        utils.Logger.Printf("ERROR: Invalid replicas count %d; must be between %d and %d", 
            replicas, MinReplicas, MaxReplicas)
        return nil, ErrInvalidReplicas
    }

    hr := &HashRing{
        replicas:    replicas,
        nodes:       make([]string, 0),
        ring:        make(map[uint32]string),
        sorted:      make([]uint32, 0),
        hashAlgo:    HashAlgorithmMD5,
        initialized: time.Now(),
    }
    utils.Logger.Printf("INFO: HashRing initialized with %d replicas using %s algorithm", 
        replicas, hr.hashAlgo)
    return hr, nil
}

// SetHashAlgorithm sets the hash algorithm to use (md5 or sha256)
// 
// Parameters:
// - algo: The algorithm identifier (HashAlgorithmMD5 or HashAlgorithmSHA256)
// 
// Returns:
// - error: An error if the algorithm is invalid or the ring already has nodes
func (h *HashRing) SetHashAlgorithm(algo string) error {
    h.mu.Lock()
    defer h.mu.Unlock()

    if len(h.nodes) > 0 {
        utils.Logger.Printf("ERROR: Cannot change hash algorithm after nodes are added")
        return errors.New("hash algorithm cannot be changed after nodes are added")
    }

    if algo != HashAlgorithmMD5 && algo != HashAlgorithmSHA256 {
        utils.Logger.Printf("ERROR: Invalid hash algorithm %s; supported: %s, %s", 
            algo, HashAlgorithmMD5, HashAlgorithmSHA256)
        return fmt.Errorf("invalid hash algorithm: %s", algo)
    }

    h.hashAlgo = algo
    utils.Logger.Printf("INFO: Hash algorithm set to %s", algo)
    return nil
}

// AddNode adds a new node to the hash ring with the specified number of virtual nodes
// 
// Parameters:
// - node: The identifier of the node (e.g., "node1:8080")
// 
// Returns:
// - error: An error if the node already exists or another issue occurs
func (h *HashRing) AddNode(node string) error {
    h.mu.Lock()
    defer h.mu.Unlock()

    // Check for duplicate nodes
    for _, existing := range h.nodes {
        if existing == node {
            utils.Logger.Printf("ERROR: Node %s already exists in the hash ring", node)
            return ErrNodeAlreadyExists
        }
    }

    // Add virtual nodes
    startTime := time.Now()
    for i := 0; i < h.replicas; i++ {
        hash := h.hash(node + strconv.Itoa(i))
        if _, exists := h.ring[hash]; exists {
            utils.Logger.Printf("WARN: Hash collision detected for %s replica %d; hash: %d", 
                node, i, hash)
        }
        h.ring[hash] = node
        h.sorted = append(h.sorted, hash)
    }

    // Sort the hash values for efficient lookup
    sort.Slice(h.sorted, func(i, j int) bool {
        return h.sorted[i] < h.sorted[j]
    })
    h.nodes = append(h.nodes, node)

    duration := time.Since(startTime)
    utils.Logger.Printf("INFO: Added node %s with %d replicas in %v", 
        node, h.replicas, duration)
    return nil
}

// RemoveNode removes a node and its virtual nodes from the hash ring
// 
// Parameters:
// - node: The identifier of the node to remove
// 
// Returns:
// - error: An error if the node doesn’t exist
func (h *HashRing) RemoveNode(node string) error {
    h.mu.Lock()
    defer h.mu.Unlock()

    // Check if node exists
    found := false
    for i, n := range h.nodes {
        if n == node {
            h.nodes = append(h.nodes[:i], h.nodes[i+1:]...)
            found = true
            break
        }
    }
    if !found {
        utils.Logger.Printf("ERROR: Node %s not found in hash ring", node)
        return fmt.Errorf("node %s not found", node)
    }

    // Remove all virtual nodes
    startTime := time.Now()
    for i := 0; i < h.replicas; i++ {
        hash := h.hash(node + strconv.Itoa(i))
        delete(h.ring, hash)
        // Rebuild sorted list
        for j, hVal := range h.sorted {
            if hVal == hash {
                h.sorted = append(h.sorted[:j], h.sorted[j+1:]...)
                break
            }
        }
    }

    duration := time.Since(startTime)
    utils.Logger.Printf("INFO: Removed node %s and %d replicas in %v", 
        node, h.replicas, duration)
    return nil
}

// GetNode returns the node responsible for the given key
// 
// Parameters:
// - key: The key to find the node for
// 
// Returns:
// - string: The node identifier
// - error: An error if the ring is empty
func (h *HashRing) GetNode(key string) (string, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    if len(h.sorted) == 0 {
        utils.Logger.Printf("ERROR: Attempted to get node for key %s on empty hash ring", key)
        return "", ErrEmptyRing
    }

    startTime := time.Now()
    hash := h.hash(key)
    idx := sort.Search(len(h.sorted), func(i int) bool {
        return h.sorted[i] >= hash
    })

    // If hash exceeds all values, wrap around to the first node
    if idx == len(h.sorted) {
        idx = 0
    }

    node := h.ring[h.sorted[idx]]
    duration := time.Since(startTime)
    utils.Logger.Printf("DEBUG: Key %s hashed to %d, assigned to node %s in %v", 
        key, hash, node, duration)
    return node, nil
}

// GetNodes returns all nodes in the hash ring
// 
// Returns:
// - []string: A copy of the nodes list
func (h *HashRing) GetNodes() []string {
    h.mu.RLock()
    defer h.mu.RUnlock()

    nodesCopy := make([]string, len(h.nodes))
    copy(nodesCopy, h.nodes)
    utils.Logger.Printf("DEBUG: Retrieved %d nodes from hash ring", len(nodesCopy))
    return nodesCopy
}

// hash computes the hash value for a given key using the configured algorithm
// 
// Parameters:
// - key: The key to hash
// 
// Returns:
// - uint32: The 32-bit hash value
func (h *HashRing) hash(key string) uint32 {
    switch h.hashAlgo {
    case HashAlgorithmMD5:
        hash := md5.Sum([]byte(key))
        // Use first 4 bytes of MD5 hash
        return binary.BigEndian.Uint32(hash[:4])
    case HashAlgorithmSHA256:
        hash := sha256.Sum256([]byte(key))
        // Use first 4 bytes of SHA-256 hash
        return binary.BigEndian.Uint32(hash[:4])
    default:
        // Fallback to MD5 (should never happen due to validation)
        utils.Logger.Printf("WARN: Unknown hash algorithm %s, falling back to %s", 
            h.hashAlgo, HashAlgorithmMD5)
        hash := md5.Sum([]byte(key))
        return binary.BigEndian.Uint32(hash[:4])
    }
}

// Size returns the number of nodes in the hash ring
// 
// Returns:
// - int: The number of physical nodes
func (h *HashRing) Size() int {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return len(h.nodes)
}

// VirtualNodeCount returns the total number of virtual nodes
// 
// Returns:
// - int: The total number of virtual nodes
func (h *HashRing) VirtualNodeCount() int {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return len(h.sorted)
}

// Clear removes all nodes from the hash ring
func (h *HashRing) Clear() {
    h.mu.Lock()
    defer h.mu.Unlock()

    h.nodes = make([]string, 0)
    h.ring = make(map[uint32]string)
    h.sorted = make([]uint32, 0)
    utils.Logger.Printf("INFO: Hash ring cleared")
}

// InitializationTime returns the time when the hash ring was created
// 
// Returns:
// - time.Time: The initialization timestamp
func (h *HashRing) InitializationTime() time.Time {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return h.initialized
}

// String provides a string representation of the hash ring
func (h *HashRing) String() string {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return fmt.Sprintf("HashRing{replicas=%d, nodes=%d, virtualNodes=%d, algo=%s}", 
        h.replicas, len(h.nodes), len(h.sorted), h.hashAlgo)
}

// Additional utility methods could be added here, such as:
// - BulkAddNodes for adding multiple nodes at once
// - Rebalance for redistributing keys after node addition/removal
// - ExportState for debugging or persistence
// - ImportState for restoring from a previous state

// Example expansion ideas:
// - Add metrics collection (e.g., Prometheus) for node distribution
// - Implement a health check mechanism for nodes
// - Add support for weighted nodes (different replica counts per node)