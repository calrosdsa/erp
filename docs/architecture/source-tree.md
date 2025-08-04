# Source Tree

```
erp-project/
├── erp/                                    # Go Backend (Modular Monolith)
│   ├── cmd/                               # Application entry points
│   │   ├── all/                          # Main monolith entry (loads all modules)
│   │   │   ├── main.go                   # Primary application entry
│   │   │   └── internal/
│   │   ├── app/                          # Alternative single-app entry
│   │   └── nats/                         # NATS server entry
│   │
│   ├── project/                          # Business Domain Modules
│   │   ├── accounting/                   # Financial management
│   │   │   ├── ledger/
│   │   │   │   ├── handler/
│   │   │   │   │   ├── rest/             # REST API endpoints
│   │   │   │   │   └── event/            # Domain event handlers
│   │   │   │   ├── repository/           # Data access layer
│   │   │   │   ├── usecase/              # Business logic
│   │   │   │   └── module.go             # Module registration
│   │   │   └── module.go
│   │   │
│   │   ├── stock/                        # Inventory management
│   │   │   ├── item/
│   │   │   ├── warehouse/
│   │   │   ├── stock_entry/
│   │   │   └── module.go
│   │   │
│   │   ├── core/                         # Core system functionality
│   │   │   ├── activity/
│   │   │   ├── address/
│   │   │   ├── contact/
│   │   │   └── module.go
│   │   │
│   │   └── selling/                      # Sales operations
│   │       ├── customer/
│   │       └── module.go
│   │
│   ├── pkg/                              # Shared packages (reusable)
│   │   ├── system/                       # Core system infrastructure
│   │   ├── db/                           # Database utilities
│   │   ├── di/                           # Dependency injection
│   │   ├── bus/                          # Event bus
│   │   └── logger/                       # Logging utilities
│   │
│   ├── gen/                              # Generated code (DO NOT EDIT)
│   │   ├── db/
│   │   │   ├── model/                    # GORM generated models
│   │   │   └── query/                    # GORM generated queries
│   │   └── mocks/                        # Generated test mocks
│   │
│   ├── configs/                          # Configuration files
│   │   └── config.json                   # Main configuration
│   │
│   ├── docs/                             # Documentation
│   │   ├── prd.md                        # Product requirements
│   │   └── backend-architecture.md       # This document
│   │
│   ├── go.mod                            # Go module definition
│   └── Makefile                          # Build and development tasks
│
└── web/                                   # React Frontend (Future)
    ├── src/
    │   ├── modules/                      # Domain-specific modules
    │   │   ├── accounting/
    │   │   ├── stock/
    │   │   └── selling/
    │   └── shared/                       # Shared frontend utilities
    └── package.json
```
