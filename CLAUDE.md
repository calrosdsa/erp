# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a modular Go-based ERP system built with a monolithic architecture that loads individual business modules:

### Core System Structure
- **Main entry point**: `cmd/all/main.go` - loads all modules in a monolith pattern
- **System foundation**: `pkg/system/` - provides core infrastructure (DB, event bus, DI container, etc.)
- **Domain-driven design**: Each business domain is a separate module in `project/` directory

### Module Organization
Business modules follow a consistent pattern:
```
project/<module>/
├── handler/          # HTTP REST handlers and event handlers
│   ├── rest/        # REST API endpoints
│   └── event/       # Domain event handlers
├── repository/      # Data access layer (PostgreSQL)
├── usecase/         # Business logic layer
├── pkg/            # Module-specific utilities
└── module.go       # Module startup and DI configuration
```

### Key Business Modules
- **Core**: Users, companies, addresses, contacts, activities, notifications
- **Stock**: Items, inventory, warehouses, price lists, stock entries
- **Accounting**: Ledgers, journals, cost centers, financial reports
- **Sales/Buying**: Orders, invoices, customers, suppliers, quotations
- **Project Management**: Projects, tasks, time tracking
- **Document Management**: File handling, terms & conditions

### Technology Stack
- **Framework**: Echo v4 for HTTP API with Huma v2 for OpenAPI
- **Database**: PostgreSQL with GORM ORM and code generation (`gen/db/`)
- **Event System**: NATS JetStream for event streaming
- **Authentication**: JWT with Permify for authorization
- **Observability**: OpenTelemetry for tracing and logging
- **Code Generation**: Uses GORM Gen for models and queries

### Database and Code Generation
- **Generated models**: `gen/db/model/` - auto-generated from database schema
- **Generated queries**: `gen/db/query/` - type-safe database operations
- **Schema**: `db/schema.sql` and `db/init.sql` for database initialization

### Configuration
- **Main config**: `configs/config.json` - database, API, NATS, observability settings
- **Environment**: Uses `.env` files for sensitive configuration
- **Module configs**: Each module can have specific configuration in `internal/app/config/`

### Event-Driven Architecture
- **Domain events**: Defined in `internal/domain/event/` and `project/*/handler/event/`
- **Event topics**: Registered in `internal/domain/events_topics.go`
- **Event bus**: Uses NATS for pub/sub messaging between modules

### Development Patterns
- **Dependency Injection**: Uses custom DI container in `pkg/di/`
- **Repository Pattern**: Data access abstracted through interfaces
- **Use Case Pattern**: Business logic encapsulated in use case services
- **FSM (Finite State Machine)**: For complex business process flows in `pkg/fsm/`



## Common Development Tasks

### Adding a New Business Module
1. Create module structure in `project/<module-name>/`
2. Implement handler/repository/usecase layers
3. Add module registration in `cmd/all/main.go`
4. Define domain events in `internal/domain/event/`
5. Run `make generate` to update generated code

### Database Changes
1. Update schema in `db/schema.sql`
2. Run `make models` to regenerate GORM models
3. Update repository interfaces as needed
4. Test with `make test`

### Adding New API Endpoints
1. Define DTOs in `api/dto/`
2. Implement handlers in `project/<module>/handler/rest/`
3. Add path definitions in corresponding `paths.go`
4. Update middleware in `api/middlewares/` if needed