# ERP System

A modular Go-based Enterprise Resource Planning (ERP) system built with a monolithic architecture that loads individual business modules.

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL
- NATS JetStream
- Make

### Setup
1. Clone the repository
2. Copy environment configuration:
   ```bash
   cp .env.example .env
   ```
3. Update database connection in `configs/config.json`
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Generate models and queries:
   ```bash
   make generate
   ```
6. Run the application:
   ```bash
   make run
   ```

## Architecture

### System Structure
- **Main entry point**: `cmd/all/main.go` - loads all modules in a monolith pattern
- **System foundation**: `pkg/system/` - provides core infrastructure (DB, event bus, DI container)
- **Business modules**: `project/` directory contains domain-specific modules

### Module Organization
Each business module follows this structure:
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

### Business Modules
- **Core**: Users, companies, addresses, contacts, activities, notifications
- **Stock**: Items, inventory, warehouses, price lists, stock entries
- **Accounting**: Ledgers, journals, cost centers, financial reports
- **Sales/Buying**: Orders, invoices, customers, suppliers, quotations
- **Project Management**: Projects, tasks, time tracking
- **Document Management**: File handling, terms & conditions

## Technology Stack
- **Framework**: Echo v4 for HTTP API with Huma v2 for OpenAPI
- **Database**: PostgreSQL with GORM ORM
- **Event System**: NATS JetStream for event streaming
- **Authentication**: JWT with Permify for authorization
- **Observability**: OpenTelemetry for tracing and logging
- **Code Generation**: GORM Gen for models and queries

## Development

### Common Commands
```bash
# Generate models and queries from database schema
make generate

# Run tests
make test

# Build application
make build

# Run application
make run

# Generate database models
make models

# Generate query annotations
make annotations
```

### Adding a New Module
1. Create module structure in `project/<module-name>/`
2. Implement handler/repository/usecase layers
3. Add module registration in `cmd/all/main.go`
4. Define domain events in `internal/domain/event/`
5. Run `make generate` to update generated code

### Database Changes
1. Update schema in `db/schema.sql`
2. Run `make models` to regenerate GORM models
3. Run `make annotations` to regenerate model queries
4. Update repository interfaces as needed

### API Development
1. Define DTOs in `api/dto/`
2. Implement handlers in `project/<module>/handler/rest/`
3. Add path definitions in corresponding `paths.go`
4. Update middleware in `api/middlewares/` if needed

## Configuration
- **Main config**: `configs/config.json` - database, API, NATS, observability settings
- **Environment**: Uses `.env` files for sensitive configuration
- **Module configs**: Each module can have specific configuration

## Generated Code
- **Models**: `gen/db/model/` - auto-generated from database schema
- **Queries**: `gen/db/query/` - type-safe database operations
- **Proto**: `gen/proto/` - generated from protobuf definitions

## Event-Driven Architecture
The system uses NATS JetStream for event-driven communication between modules:
- **Domain events**: Defined in `internal/domain/event/`
- **Event topics**: Registered in `internal/domain/events_topics.go`
- **Event handlers**: Implemented in `project/*/handler/event/`

## Development Patterns
- **Dependency Injection**: Custom DI container in `pkg/di/`
- **Repository Pattern**: Data access abstracted through interfaces
- **Use Case Pattern**: Business logic encapsulated in use case services
- **FSM (Finite State Machine)**: For complex business process flows in `pkg/fsm/`