System Design: Twitter with 30 Files and 10,000 Lines of Code
In this system design exercise, we will create a scalable Twitter-like platform, implemented across 30 files with a total of approximately 10,000 lines of code. Twitter is a free microblogging social network where users post short messages called "tweets," engage with others through likes, replies, and retweets, and follow accounts to curate their timelines. With around 397 million users as of 2021, Twitter handles massive scale, requiring a robust and efficient architecture. This design will focus on functional and non-functional requirements, a high-level architecture with client-side load balancing, and a detailed file structure to meet the specified constraints.
Requirements
Functional Requirements
User Management: Register accounts, authenticate users, manage profiles, and handle follow/unfollow actions.

Tweet Management: Post tweets (up to 280 characters), retrieve tweets, and support likes and retweets.

Timeline Generation: Generate personalized timelines based on followed users’ tweets.

Search: Enable searching for tweets and users by keywords or hashtags.

Notifications: Send real-time notifications for mentions, likes, and retweets.

Non-Functional Requirements
Scalability: Handle millions of users and billions of tweets, with capacity to scale horizontally.

Performance: Low-latency timeline retrieval and tweet posting (e.g., <1s response time).

Availability: Achieve 99.9% uptime through redundancy and fault tolerance.

Consistency: Ensure eventual consistency for timelines and notifications in a distributed system.

Estimations
Storage: Assuming 1KB per tweet (including metadata) and 1 billion tweets daily, storage is ~1TB/day.

Bandwidth: With 397 million users averaging 10 timeline refreshes daily (1KB each), bandwidth is ~4TB/day.

Requests: 10,000 requests/second for tweet posting and timeline retrieval combined.

High-Level Design
The system adopts a microservices architecture with client-side load balancing, distributing functionality across independent services that communicate via HTTP/REST APIs. Key components include:
Major Components
User Service: Manages user accounts, authentication (e.g., JWT), and relationships (following).

Tweet Service: Handles tweet creation, storage, and retrieval.

Timeline Service: Precomputes and serves user timelines using caching.

Search Service: Indexes and searches tweets/users with a simplified in-memory store.

Notification Service: Delivers real-time updates via a message queue.

Load Balancer: Embedded in each service to select healthy instances of other services.

Supporting Systems:
Database: SQL (e.g., SQLite) for user data; NoSQL (e.g., key-value store) for tweets.

Cache: In-memory store (e.g., Redis) for timelines and hot tweets.

Message Queue: Asynchronous communication (e.g., RabbitMQ) for timeline updates and notifications.

Data Flow Example: Posting a Tweet
Client sends a POST request to the Tweet Service.

Tweet Service authenticates the user (via User Service), stores the tweet in the NoSQL database, and publishes a message to the queue.

Timeline Service consumes the message, updates followers’ cached timelines.

Notification Service notifies mentioned users or followers asynchronously.

API Design (Simplified)
POST /tweets: Create a tweet ({user_id, content}).

GET /timeline/{user_id}: Fetch a user’s timeline.

GET /users/{user_id}: Retrieve user profile.

GET /search?q={query}: Search tweets/users.

Top-K Problem
For tweets with millions of likes (e.g., trending posts), we use caching to store precomputed "hot" tweet lists, refreshed periodically by a background job in the Tweet Service.

