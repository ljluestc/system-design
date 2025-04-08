System Design for a Simplified Load Balancer
This load balancer system is designed to handle millions of requests by distributing them across a pool of backend servers. It sits between clients and servers, providing scalability, availability, and performance. The system is containerized using Docker and includes a frontend for monitoring, a backend for request handling, and additional services like health checking and TLS termination.
High-Level Architecture
Frontend: A React-based web interface to monitor load balancer statistics and server health.

Backend: A Node.js/Express server that acts as the load balancer, distributing requests to backend servers.

Health Checking: Periodically verifies the operational status of backend servers.

TLS Termination: Manages HTTPS requests securely.

Service Discovery: Dynamically registers and discovers backend services.

Logging: Uses Winston for event and error logging.

Infrastructure: Docker for containerization and NGINX as a reverse proxy.

