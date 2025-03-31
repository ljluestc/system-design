package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "sync"
    "time"
)

// Backend represents a single server in the pool
type Backend struct {
    URL          *url.URL
    Alive        bool
    mu           sync.RWMutex
    ReverseProxy *httputil.ReverseProxy
}

// SetAlive updates the health status of the backend
func (b *Backend) SetAlive(alive bool) {
    b.mu.Lock()
    b.Alive = alive
    b.mu.Unlock()
}

// IsAlive returns the current health status
func (b *Backend) IsAlive() bool {
    b.mu.RLock()
    defer b.mu.RUnlock()
    return b.Alive
}

// LoadBalancer manages a pool of backend servers
type LoadBalancer struct {
    backends []*Backend
    current  int
    mu       sync.Mutex
}

// NewLoadBalancer creates a new LoadBalancer with the given backend URLs
func NewLoadBalancer(backendURLs []string) (*LoadBalancer, error) {
    var backends []*Backend
    for _, u := range backendURLs {
        parsedURL, err := url.Parse(u)
        if err != nil {
            return nil, fmt.Errorf("invalid backend URL %s: %v", u, err)
        }
        backend := &Backend{
            URL:          parsedURL,
            Alive:        true,
            ReverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
        }
        backends = append(backends, backend)
    }
    return &LoadBalancer{backends: backends}, nil
}

// NextBackend selects the next available backend using round-robin
func (lb *LoadBalancer) NextBackend() *Backend {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    for i := 0; i < len(lb.backends); i++ {
        lb.current = (lb.current + 1) % len(lb.backends)
        backend := lb.backends[lb.current]
        if backend.IsAlive() {
            return backend
        }
    }
    return nil // No healthy backends available
}

// ServeHTTP implements the HTTP handler interface
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    backend := lb.NextBackend()
    if backend == nil {
        http.Error(w, "Service Unavailable: No healthy backends", http.StatusServiceUnavailable)
        return
    }
    log.Printf("Forwarding request to %s", backend.URL.String())
    backend.ReverseProxy.ServeHTTP(w, r)
}

// HealthCheck periodically checks the health of backends
func (lb *LoadBalancer) HealthCheck() {
    for {
        for _, backend := range lb.backends {
            resp, err := http.Get(backend.URL.String() + "/health") // Assume /health endpoint
            alive := err == nil && resp != nil && resp.StatusCode == http.StatusOK
            if err == nil && resp != nil {
                resp.Body.Close()
            }
            backend.SetAlive(alive)
            log.Printf("Health check for %s: %v", backend.URL.String(), alive)
        }
        time.Sleep(10 * time.Second) // Check every 10 seconds
    }
}

// StartBackend starts a simple backend server for testing
func StartBackend(port string) {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from backend on port %s!", port)
    })
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "OK")
    })
    log.Printf("Starting backend server on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}

func main() {
    // Start backend servers for testing
    go StartBackend("8081")
    go StartBackend("8082")
    go StartBackend("8083")
    time.Sleep(1 * time.Second) // Give backends time to start

    // Initialize load balancer
    backendURLs := []string{
        "http://localhost:8081",
        "http://localhost:8082",
        "http://localhost:8083",
    }
    lb, err := NewLoadBalancer(backendURLs)
    if err != nil {
        log.Fatalf("Failed to create load balancer: %v", err)
    }

    // Start health checking
    go lb.HealthCheck()

    // Start load balancer server
    log.Println("Starting load balancer on :8080")
    log.Fatal(http.ListenAndServe(":8080", lb))
}