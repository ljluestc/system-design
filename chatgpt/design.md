Part 1: System Design for ChatGPT
ChatGPT is a conversational AI service that processes user queries, generates human-like responses using a large language model (LLM), and handles millions of concurrent users. I’ll design it in four stages, mirroring the Google Docs approach: requirements, component design, concurrency, and evaluation.
Step 1: Requirements of ChatGPT’s Design
Functional Requirements
Based on ChatGPT’s real-world usage, users expect:
Query Processing: Users submit text prompts (questions, commands) and receive coherent, context-aware responses.

Context Retention: Maintains conversation history for multi-turn dialogues.

Content Moderation: Filters harmful or inappropriate inputs/outputs.

Personalization: Adapts responses based on user preferences (if authenticated).

Metrics Tracking: Tracks usage (e.g., query count, response time) for analytics.

Multilingual Support: Handles queries in multiple languages.

Non-Functional Requirements
Latency: Responses within 1–3 seconds for 99th percentile users.

Consistency: Eventual consistency for conversation history; strong consistency for user authentication.

Scalability: Supports 100M daily active users (DAU), ~1M queries/second at peak.

Availability: 99.99% uptime, resilient to node failures.

Cost Efficiency: Balances compute-intensive LLM inference with affordable infrastructure.

Resource Estimations
Assumptions:
100M DAU, 1 query/day average, peak 1M queries/second.

Average prompt: 50 tokens (200 bytes); response: 100 tokens (400 bytes).

Conversation history: 10KB per user (10 turns).

Model size: 175B parameters (like GPT-3), ~350GB in FP16.

Storage:
Queries/day: 100M × (200B + 400B) = 60TB/day (compressed to ~20TB with deduplication).

History: 100M × 10KB = 1TB/day.

Model weights: 350GB (static, replicated across regions).

Total (1 day): ~21TB (excluding history older than 1 day).

Bandwidth:
Incoming: 20TB/day ÷ 86,400s × 8 = ~2Gbps.

Outgoing (5x views): 100TB/day ÷ 86,400s × 8 = ~10Gbps.

Total: ~12Gbps.

Compute:
Inference: 1 query = 175B × 2 FLOPs × 100 tokens = 3.5T FLOPs.

At 100 TFLOPs/GPU (e.g., A100), 1 GPU handles ~30 queries/second.

Peak 1M queries/second ÷ 30 = ~33,333 GPUs.

Servers (8 GPUs/server): ~4,167 servers.

Servers:
API servers: 1M requests/second ÷ 50,000 RPS/server = ~20 servers.

Total: ~4,200 servers (rounded).

Step 2: Design of ChatGPT
Components
Inspired by Google Docs’ design, I’ll use modular components:
API Gateway: Routes user requests, handles authentication, rate-limiting, and load balancing. Uses REST/gRPC for queries, WebSockets for streaming responses.

Inference Service: Runs LLM inference (e.g., GPT-like model) on GPU clusters. Shards model across nodes for parallelism.

Context Store: NoSQL database (e.g., DynamoDB) for conversation history, enabling multi-turn dialogues.

Moderation Service: Filters inputs/outputs using rule-based and ML models.

Cache: Redis for frequent queries (e.g., FAQs) and session data.

Blob Storage: Stores model weights, logs, and training data snapshots.

Message Queue: Kafka for async tasks (e.g., logging, moderation).

Pub-Sub: For notifications (e.g., usage alerts).

Monitoring: Prometheus for latency, error rates, and GPU utilization.

CDN: Serves static assets (e.g., UI, model metadata).

Architecture

[Client] --> [API Gateway] --> [Load Balancer]
                                   |
        --------------------------------------
        |                |                   |
 [Inference Service]  [Context Store]  [Moderation Service]
        |                |                   |
     [Cache]        [Blob Storage]     [Message Queue]
        |                                    |
   [Monitoring]                        [Pub-Sub]

Workflow:
User sends prompt via API/WebSocket.

API Gateway authenticates, rate-limits, and routes to Inference Service.

Inference Service fetches context from Context Store, runs LLM, and caches results.

Moderation Service filters output.

Response returns via API Gateway; history updates in Context Store.

Async tasks (logs, metrics) go to Kafka/Pub-Sub.

Why These Components?
API Gateway: Centralizes request handling, like Google Docs’ gateway for edits.

Inference Service: GPU-heavy for LLM, unlike Google Docs’ CPU-based app servers.

Context Store: NoSQL for flexible history, similar to Google Docs’ time-series DB.

Moderation: Unique to ChatGPT for safety, absent in Google Docs.

Cache/CDN: Reduces latency, like Google Docs’ Redis/CDN for documents.

Step 3: Concurrency in ChatGPT
Concurrency Challenges
Inference Overload: Millions of concurrent queries overwhelm GPUs.

Context Conflicts: Multiple queries from one user updating history simultaneously.

Moderation Latency: Real-time filtering slows responses.

Solutions
Inference Concurrency:
Sharding: Partition model weights across GPUs (tensor parallelism).

Queueing: Use Kafka to batch inference requests, prioritizing low-latency queries.

Unlike Google Docs’ OT/CRDT for text edits, ChatGPT uses stateless inference per query, avoiding commutative/idempotent issues.

Context Store:
Eventual Consistency: Use DynamoDB with optimistic locking for history updates.

Similar to Google Docs’ time-series DB but simpler, as conflicts are rare (per-user history).

Moderation:
Async Moderation: Moderate outputs post-inference in parallel, caching safe results.

Differs from Google Docs, which needs real-time conflict resolution.

Why Not OT/CRDT?
ChatGPT’s queries are mostly independent, unlike collaborative editing’s shared document state, making OT/CRDT overkill.
Step 4: Evaluation of ChatGPT’s Design
Non-Functional Requirements
Consistency:
Eventual consistency for history via DynamoDB replication.

Strong consistency for authentication via relational DB.

Gossip protocol (like Google Docs) for intra-data-center sync.

Latency:
Cache hits (<100ms) for FAQs.

Inference (~1s) optimized by GPU sharding.

WebSockets for streaming reduce perceived latency.

Scalability:
Horizontal scaling: Add GPU nodes for inference, DynamoDB for history.

Microservices allow independent scaling (like Google Docs’ queues).

Availability:
Multi-region deployment with failover.

Cache/CDN for read-heavy queries.

Monitoring ensures 99.99% uptime.

Trade-offs
Latency vs. Consistency: Async moderation sacrifices strict safety for speed.

Cost vs. Performance: GPU clusters are expensive; caching offsets costs.

Scalability vs. Complexity: Sharding adds complexity but enables millions of queries/second.

