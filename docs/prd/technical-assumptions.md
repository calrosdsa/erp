# Technical Assumptions

## Repository Structure: Monorepo

The system will maintain the existing monolithic repository structure with modular loading pattern via `cmd/all/main.go`. This approach provides the benefits of simplified dependency management, atomic commits across modules, and unified CI/CD while preserving the modular architecture through the `/project/<module>` organization. The monorepo enables efficient code sharing between modules while maintaining clear domain boundaries.

## Service Architecture

**CRITICAL DECISION:** The system will continue the current domain-driven microservices within monolith approach. Each business module (Core, Stock, Accounting, Sales/Buying, Project Management) operates as an independent service within the monolithic deployment, communicating through NATS JetStream events and shared database access. This architecture provides modularity benefits without the operational complexity of distributed microservices, enabling rapid development while maintaining the flexibility to extract services later if needed.

## Testing Requirements

**CRITICAL DECISION:** The system will implement a comprehensive testing pyramid with emphasis on integration testing due to the event-driven architecture:

- **Unit Testing:** Table-driven tests with GoMock for interface mocking, achieving 90% coverage for business logic
- **Integration Testing:** Testcontainers-go with real PostgreSQL and NATS instances for realistic testing scenarios
- **BDD Testing:** Godog (Cucumber for Go) for stakeholder-readable business workflow specifications
- **Contract Testing:** Pact framework ensuring API compatibility between modules and external integrations
- **End-to-End Testing:** Playwright for complete user journey validation across web interface
- **Performance Testing:** k6 for load testing with automated baseline comparisons
- **Chaos Testing:** Simulated failures for network partitions, module crashes, and database issues

## Additional Technical Assumptions and Requests


### Backend Technology Choices
- **Language:** Go 1.21+ leveraging existing team expertise and performance characteristics
- **Web Framework:** Continue with Echo v4 + Huma v2 for OpenAPI-first development
- **Database:** PostgreSQL 14+ with GORM ORM, maintaining existing code generation pipeline
- **Message Queue:** NATS JetStream for event-driven communication between modules
- **Caching:** Redis for session management and frequently accessed data

### DevOps & Infrastructure
- **Containerization:** Docker with multi-stage builds for optimized production images
- **Orchestration:** Kubernetes with Helm charts for both cloud and on-premises deployment
- **CI/CD:** GitHub Actions with automated testing, security scanning, and deployment pipelines
- **Monitoring:** OpenTelemetry with Prometheus metrics and Jaeger distributed tracing
- **Database Migrations:** Custom migration system integrated with existing GORM Gen pipeline

### Security Architecture
- **Authentication:** JWT tokens with refresh token rotation and secure storage
- **Authorization:** Permify for fine-grained, policy-based access control
- **API Security:** Rate limiting, request validation, and automated security testing in CI/CD
- **Data Protection:** Encryption at rest using database-level encryption, TLS 1.3 for all communications

### Integration Strategy
- **API Design:** REST-first with GraphQL consideration for complex data fetching scenarios
- **Webhook Support:** Outbound webhooks for real-time integration with customer systems
- **Batch Processing:** Background job processing using NATS JetStream for heavy operations
- **File Storage:** S3-compatible object storage for documents and file attachments
