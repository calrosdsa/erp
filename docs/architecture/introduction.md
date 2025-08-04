# Introduction

This document outlines the overall project architecture for Modular ERP System, including backend systems, shared services, and non-UI specific concerns. Its primary goal is to serve as the guiding architectural blueprint for AI-driven development, ensuring consistency and adherence to chosen patterns and technologies.

**Relationship to Frontend Architecture:**
If the project includes a significant user interface, a separate Frontend Architecture Document will detail the frontend-specific design and MUST be used in conjunction with this document. Core technology stack choices documented herein (see "Tech Stack") are definitive for the entire project, including any frontend components.

## Starter Template or Existing Project Analysis

This is **NOT** a greenfield project. You have an existing, sophisticated Go-based ERP system with:

**Existing Project Foundation:**
- **Architecture Pattern:** Domain-driven microservices within monolith
- **Main Entry Point:** `cmd/all/main.go` - monolithic loading of modular components
- **Module Structure:** `/project/<module>` organization pattern already established
- **Core Infrastructure:** `pkg/system/` provides database, event bus, DI container
- **Technology Stack:** Go + Echo v4 + Huma v2 + PostgreSQL + GORM + NATS JetStream
- **Code Generation:** Established pipeline with GORM Gen for models and queries
- **Event Architecture:** NATS JetStream for inter-module communication

**Pre-configured Capabilities:**
- JWT authentication with Permify authorization
- Multi-tenancy support through Core module
- Event-driven architecture with defined event topics
- Database schema and migration system
- OpenAPI documentation generation
- Observability with OpenTelemetry

## Change Log
| Date | Version | Description | Author |
|------|---------|-------------|---------|
| 2025-08-04 | 1.0 | Initial architecture document creation | Claude |
