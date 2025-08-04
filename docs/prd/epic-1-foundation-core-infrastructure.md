# Epic 1: Foundation & Core Infrastructure

**Epic Goal:** Establish a production-ready application foundation with authentication, deployment infrastructure, and basic user management capabilities while delivering a functional health-check endpoint that demonstrates system reliability and provides immediate operational value for monitoring and deployment validation.

## Story 1.1: Project Setup & Development Environment

As a **developer**,
I want **a complete development environment with hot-reload, testing, and code generation capabilities**,
so that **I can efficiently build and test ERP functionality with consistent tooling**.

### Acceptance Criteria
1. Go project initialized with proper module structure following existing `/project/<module>` pattern
2. React application scaffolded with TypeScript, Vite build system, and Tailwind CSS configuration
3. Docker Compose development environment including PostgreSQL, NATS JetStream, and Redis containers
4. Makefile commands for `make dev`, `make test`, `make build`, and `make generate` operations
5. Hot-reload functionality for both Go backend and React frontend during development
6. Code generation pipeline established for GORM models, OpenAPI specs, and TypeScript API clients
7. Git hooks configured for pre-commit linting, formatting, and basic test execution

## Story 1.2: Authentication & Authorization Foundation

As a **system administrator**,
I want **secure JWT-based authentication with role-based permissions**,
so that **I can control user access to different ERP modules and functions**.

### Acceptance Criteria
1. JWT authentication system implemented with access and refresh token rotation
2. Permify integration for policy-based authorization with predefined roles (Admin, User, Viewer)
3. User registration and login endpoints with proper password hashing and validation
4. Protected route middleware for Go API endpoints with role-based access control
5. React authentication context and protected route components for frontend
6. Token refresh mechanism handling expired tokens transparently for users
7. Logout functionality that properly invalidates tokens on both client and server
8. Basic user profile management (view/edit profile, change password)

## Story 1.3: Multi-Company Tenancy & User Management

As a **business owner**,
I want **to manage multiple companies and users within my ERP instance**,
so that **I can support multiple business entities or client organizations**.

### Acceptance Criteria
1. Company entity model with basic profile information (name, address, tax ID, settings)
2. User-company relationship management with role assignments per company
3. Company selection interface for users associated with multiple companies
4. Data isolation ensuring users only access data for their authorized companies
5. Company administrator role with user invitation and management capabilities
6. User invitation workflow with email notifications and secure registration links
7. Company switching functionality in the UI without requiring re-authentication
8. Audit logging for user management actions (create, update, delete, role changes)

## Story 1.4: Production Deployment & Health Monitoring

As a **system operator**,
I want **reliable deployment infrastructure with health monitoring and observability**,
so that **I can deploy and monitor the ERP system in production environments**.

### Acceptance Criteria
1. Kubernetes deployment manifests with proper resource limits and health checks
2. Health check endpoints returning system status, database connectivity, and NATS connectivity
3. OpenTelemetry integration for distributed tracing and metrics collection
4. Prometheus metrics exposed for system performance monitoring
5. Docker images optimized for production with security scanning in CI/CD pipeline
6. Database migration system that runs automatically during deployment
7. CI/CD pipeline with automated testing, security scanning, and deployment to staging
8. Rollback mechanism for failed deployments with database migration rollback support
