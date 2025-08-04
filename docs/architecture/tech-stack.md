# Tech Stack

This is the DEFINITIVE technology selection section that serves as the single source of truth for all development decisions.

## Cloud Infrastructure
- **Provider:** Multi-cloud (AWS/Azure/GCP) with Kubernetes
- **Key Services:** Container Registry, Load Balancers, PostgreSQL managed service, Redis managed service
- **Deployment Regions:** Configurable based on customer requirements (US, EU, APAC)

## Technology Stack Table

| Category | Technology | Version | Purpose | Rationale |
|----------|------------|---------|---------|-----------|
| **Language** | Go | 1.21+ | Primary backend development | Excellent performance, strong typing, team expertise, concurrent processing for ERP workloads |
| **Web Framework** | Echo | v4.11+ | HTTP API framework | Lightweight, fast, extensive middleware ecosystem, production proven |
| **API Documentation** | Huma | v2.0+ | OpenAPI-first development | Auto-generated docs, type-safe handlers, excellent DX for API-first architecture |
| **Database** | PostgreSQL | 14+ | Primary data store | ACID compliance, JSON support, excellent performance for ERP workloads, existing investment |
| **ORM** | GORM | v1.25+ | Database abstraction | Strong Go ecosystem support, code generation, migration management |
| **Code Generation** | GORM Gen | Latest | Type-safe queries | Eliminates runtime query errors, improves performance, maintains existing patterns |
| **Message Queue** | NATS JetStream | v2.10+ | Event-driven messaging | High performance, clustering support, persistence, stream processing |
| **Caching** | Redis | 7.0+ | Session storage & caching | In-memory performance, pub/sub capabilities, session management |
| **Authentication** | JWT + Permify | JWT latest, Permify v0.8+ | Auth & authorization | Stateless auth, fine-grained permissions, policy-based access control |
| **Frontend Framework** | React | 18+ | User interface | Large ecosystem, TypeScript support, component reusability |
| **Frontend Language** | TypeScript | 5.3+ | Type-safe frontend | Compile-time error detection, better developer experience, API type safety |
| **Containerization** | Docker | 24+ | Application packaging | Multi-stage builds, consistent environments, platform independence |
| **Orchestration** | Kubernetes | 1.28+ | Container orchestration | Scalability, rolling deployments, service discovery, production readiness |
| **CI/CD** | GitHub Actions | Latest | Automation pipeline | Integrated with repository, extensive marketplace, cost-effective |
| **Observability** | OpenTelemetry | v1.21+ | Distributed tracing | Vendor-neutral, comprehensive observability, performance monitoring |
