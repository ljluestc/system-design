Part 1: System Design for Distributed LRU Cache
A Distributed LRU Cache stores key-value pairs across multiple nodes, evicting the least recently used items when capacity is reached, and serves high-throughput, low-latency requests for distributed applications. I’ll design it in four stages, mirroring the Google Docs approach.
Step 1: Requirements of Distributed LRU Cache
Functional Requirements
Get/Put Operations: Retrieve (Get) and store (Put) key-value pairs with O(1) average time complexity.

LRU Eviction: Evict the least recently used key when a node’s cache reaches capacity.

Distributed Access: Route requests to the correct node based on keys.

Expiration: Support optional TTL (time-to-live) for cache entries.

Monitoring: Track hit/miss rates and latency.

Consistency: Ensure consistent views of cache across nodes for the same key.

Non-Functional Requirements
Latency: <10ms for 99th percentile Get/Put operations.

Consistency: Eventual consistency across nodes, with strong consistency within a node.

Scalability: Handle 1M requests/second across 100 nodes, with 1TB total cache capacity.

Availability: 99.99% uptime, resilient to node failures.

Fault Tolerance: Handle node crashes without data loss for non-evicted items.

Resource Estimations
Assumptions:
10M keys, average key-value size: 1KB.

100 nodes, each with 10GB cache (1TB total).

1M requests/second, 50% Get, 50% Put.

Cache hit rate: 90%.

Storage:
Per node: 10GB ÷ 1KB = 10M keys.

Total: 100 × 10GB = 1TB.

Bandwidth:
Requests: 1M × 1KB = 1GB/second.

Network: ~8Gbps across nodes (assuming replication).

Compute:
Requests/node: 1M ÷ 100 = 10K RPS.

Servers: 10K RPS ÷ 5K RPS/server = 2 servers/node, ~200 servers total.

Step 2: Design of Distributed LRU Cache
Components
API Gateway: Routes client requests (Get, Put) to appropriate nodes, handles load balancing.

Cache Node: Runs an in-memory LRU cache with local key-value storage.

Shard Manager: Determines which node owns a key using consistent hashing.

Replication Service: Replicates writes to backup nodes for fault tolerance.

Cache Store: In-memory key-value store with LRU eviction (e.g., hashmap + doubly-linked list).

Data Stores:
Cache (in-memory): Primary storage for key-value pairs.

Persistent Storage (optional): For recovery (mocked here).

Message Queue: Kafka for async replication and monitoring.

Monitoring: Prometheus for hit/miss rates, latency.

Configuration Service: Manages node membership and sharding.

Architecture

[Client] --> [API Gateway] --> [Load Balancer]
                                   |
        --------------------------------------
        |                |                   |
 [Cache Node 1]   [Cache Node 2] ...  [Cache Node N]
        |                |                   |
    [Shard Manager]  [Replication Service]  [Monitoring]
        |                                    |
       [Cache Store]                      [Message Queue]
                                             [Config Service]

Workflow:
Client sends Get/Put via API Gateway.

Shard Manager routes request to the primary node using consistent hashing.

Cache Node processes request, updates LRU order.

Replication Service async propagates writes to replicas.

Monitoring logs metrics (hits, latency).

Configuration Service updates node membership on failures.

Why These Components?
API Gateway: Centralizes routing, like Google Docs’ gateway.

Cache Node: Localized LRU logic, unlike Google Docs’ collaborative editing.

Shard Manager: Ensures scalability, similar to Google Docs’ queue partitioning.

Replication: Enhances availability, akin to Google Docs’ data replication.

Step 3: Concurrency in Distributed LRU Cache
Concurrency Challenges
Race Conditions: Multiple clients accessing the same key on a node.

Replication Lag: Inconsistent views across replicas.

Node Failures: Ensuring data availability during crashes.

Solutions
Race Conditions:
Mutex Locks: Use per-node mutexes for thread-safe Get/Put.

Unlike Google Docs’ OT/CRDT, LRU operations are simpler, requiring only local synchronization.

Replication Lag:
Async Replication: Write to primary, async replicate to backups, accepting eventual consistency.

Similar to Google Docs’ async notifications but critical for cache availability.

Node Failures:
Consistent Hashing: Redistribute keys to live nodes.

Differs from Google Docs’ document replication, focusing on key reassignment.

Step 4: Evaluation of Distributed LRU Cache
Non-Functional Requirements
Consistency:
Strong within a node (mutexes).

Eventual across nodes (async replication).

Gossip protocol for node state sync (like Google Docs).

Latency:
In-memory Get/Put: <1ms.

Network overhead: ~5ms, total <10ms.

Scalability:
Horizontal scaling: Add nodes, rebalance via consistent hashing.

Sharding reduces per-node load.

Availability:
Replication ensures data availability.

Configuration Service handles node failures.

99.99% uptime via redundancy.

Fault Tolerance:
Consistent hashing minimizes data loss.

Monitoring detects and reroutes around failures.

Trade-offs
Latency vs. Consistency: Async replication reduces latency but risks stale reads.

Scalability vs. Complexity: Sharding adds overhead but supports millions of RPS.

Availability vs. Cost: Replication increases storage but ensures uptime.

