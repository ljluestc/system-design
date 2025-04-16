Part 1: System Design for Content Delivery Network (CDN)
Stage 1: Requirements of CDN Design
Functional Requirements
Based on the course content and real-world CDNs:
Content Delivery: Serve static (images, CSS, JS) and dynamic content (e.g., API responses, streaming) from edge servers.

Caching: Store content at edge servers with TTL-based or LRU eviction.

Request Routing: Direct users to the nearest edge server based on location or load.

Content Invalidation: Purge or update cached content when the origin changes.

Search: Query cached content availability across edge servers.

Analytics: Track cache hits, latency, and bandwidth usage.

Security: Protect against DDoS attacks and secure content with HTTPS.

Non-Functional Requirements
Latency: <100ms for 99th percentile content delivery.

Scalability: Support 10M requests/second, 1PB/day bandwidth across 100 edge servers.

Availability: 99.99% uptime, resilient to failures.

Consistency: Eventual consistency for cache updates; strong consistency for metadata (e.g., user sessions).

Security: Mitigate DDoS, ensure data integrity.

Optimization: Achieve >90% cache hit ratio, minimize origin load, and optimize for diverse devices (mobile, desktop).

Resource Estimations
Assumptions:
10M daily active users (DAU), 100 requests/user/day = 1B requests/day (~11,574 RPS).

Average content size: 200KB (mixed static/dynamic).

100 edge servers, each with 1TB cache, 10Gbps bandwidth.

Cache hit ratio: 90%, 10% requests hit origin.

Storage:
Cache per edge: 1TB ÷ 200KB = ~5M objects/server.

Total cache: 100 × 1TB = 100TB.

Origin storage: 10TB (compressed).

Bandwidth:
Daily traffic: 1B × 200KB = 200TB/day.

Edge-served (90%): 180TB/day.

Origin-served (10%): 20TB/day.

Network: 200TB ÷ 86,400s × 8 = ~18.5Gbps.

Compute:
Requests/server: 11,574 RPS ÷ 100 = ~116 RPS/edge.

Servers: 116 RPS ÷ 5K RPS/server = ~3 servers/edge, ~300 total.

Stage 2: Design of CDN
Components
API Gateway: Routes client requests, enforces DDoS protection, and load balances to edge servers.

Edge Server: Caches content (LRU-based), serves requests, applies compression (e.g., Brotli).

Origin Server: Stores original content, serves uncached/dynamic requests.

Routing Service: Uses Anycast or geo-DNS to select optimal edge servers.

Cache Manager: Handles caching, eviction, and invalidation (TTL, manual purges).

Data Stores:
Cache Store: In-memory (mocked as hashmap) for edge caching.

Metadata DB: Relational (mocked) for sessions, content metadata.

Blob Storage: Stores large files (e.g., videos) at origin.

Message Queue: Kafka (mocked) for async invalidation, analytics.

Monitoring: Prometheus (mocked) for latency, hit/miss rates.

Scrubber Service: Filters malicious traffic (mocked).

Architecture

[Client] --> [API Gateway/DNS] --> [Scrubber Service]
                                   |
        --------------------------------------
        |                |                   |
 [Edge Server 1]  [Edge Server 2] ...  [Edge Server N]
        |                |                   |
    [Cache Manager]  [Origin Server]   [Routing Service]
        |                |                   |
      [Cache Store]  [Blob Storage]     [Message Queue]
                        [Metadata DB]      [Monitoring]

Workflow:
Client requests content via DNS or API Gateway.

Scrubber Service filters malicious traffic.

Routing Service selects nearest edge server using Anycast.

Edge Server checks cache; serves if hit, else fetches from origin.

Cache Manager updates LRU, applies compression.

Origin Server provides uncached content.

Message Queue processes invalidations, analytics.

Monitoring tracks performance.

Optimization Strategies
Caching: LRU eviction, high TTL for static content, predictive caching for hot objects.

Routing: Anycast for low-latency server selection, geo-DNS fallback.

Compression: Brotli for text, WebP for images.

Edge Processing: Minify CSS/JS, resize images at edge.

Network: Pre-warmed connections to origin, optimized BGP routes.

Security: Rate-limiting, TLS termination at edge.

Stage 3: In-Depth Investigation of CDN
Push vs. Pull Models
Push CDN:
Origin proactively sends content to edge servers.

Ideal for static content (e.g., logos, CSS).

Optimizes availability but risks redundant pushes for dynamic data.

Pull CDN:
Edge pulls content on demand.

Suited for dynamic content (e.g., API responses).

Reduces storage but may increase origin load.

Implementation: Hybrid approach—push for static, pull for dynamic.

Dynamic Content Optimization
Edge Scripts: Run serverless scripts (e.g., image resizing) at edge to reduce origin calls.

Compression: Mock Brotli in edge/compressor.go (real systems use Railgun or ESI).

ESI (Edge Side Includes): Cache static page parts, fetch dynamic fragments.

Multi-Tier Architecture
Tiers: Edge servers → parent servers → origin.

Purpose: Reduces origin load, handles long-tail content.

Implementation: Edge servers (edge/server.go) query origin (origin/server.go) if cache misses.

Routing Mechanisms
Anycast: Single IP for all edges, BGP routes to nearest server (routing/anycast.go).

DNS Redirection: Geo-DNS selects edge based on client location (routing/geo.go).

HTTP Redirection: Fallback for legacy clients (not implemented).

Consistency Mechanisms
TTL: Content expires after a set time, edge re-fetches (edge/cache.go).

Periodic Polling: Mocked in queue/processor.go for updates.

Leases: Not implemented (complex for demo).

Deployment
Placement: Edge servers at ISPs or IXPs for low latency.

Count: 100 edges, scalable via config (config/config.go).

Stage 4: Evaluation of CDN Design
Non-Functional Requirements
Latency:
Cache hits: <10ms (in-memory).

Origin fetches: <100ms (optimized routes).

Compression reduces payload size.

Scalability:
Horizontal scaling: Add edge servers.

Anycast distributes load evenly.

Availability:
Redundant edges, failover via DNS.

99.99% uptime with monitoring.

Consistency:
Eventual for cache (TTL-based).

Strong for metadata (mocked DB).

Security:
Scrubber filters DDoS (mocked).

TLS assumed at edge.

Optimization:
90% hit ratio via LRU.

Compression and minification reduce bandwidth.

Trade-offs
Latency vs. Consistency: TTL risks stale content but lowers latency.

Cost vs. Performance: More edges improve speed but raise costs.

Complexity vs. Scalability: Anycast adds routing overhead but ensures global reach.

