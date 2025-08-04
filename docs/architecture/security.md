# Security

## Existing Security Infrastructure

**Authentication System:**
- JWT-based authentication with multi-tenant scoping via company_id
- Refresh token rotation with secure storage and transmission
- Session management integrated with Redis for scalability

**Authorization Framework:**
- Permify integration providing fine-grained, policy-based access control
- Role-based access control (RBAC) with hierarchical role inheritance
- Resource-level permissions with dynamic policy evaluation

**Data Protection:**
- Multi-tenant data isolation through company_id foreign key constraints
- Encryption at rest using PostgreSQL built-in encryption capabilities
- TLS 1.3 for all network communications between services

## Security Requirements for New Features

**Authentication Integration:**
- All new endpoints MUST implement JWT token validation middleware
- Company-scoped authentication ensuring users only access authorized tenant data
- Consistent error handling for authentication failures (401 Unauthorized)

**Authorization Implementation:**
- Mandatory Permify permission checks for all business operations
- Resource-based authorization using entity IDs and user context
- Graceful degradation for authorization service unavailability

**Data Classification and Encryption:**
- PII data MUST be encrypted using application-level encryption before database storage
- Financial data requires additional audit logging with immutable transaction records
- Sensitive fields (passwords, tokens, API keys) MUST never appear in logs or error messages
