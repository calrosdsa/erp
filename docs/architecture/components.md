# Components

## System Foundation Layer (`pkg/system/`)
**Primary Responsibility**: Provides core infrastructure services and system initialization
- **Key Interfaces**: `system.Service`, `system.Module`
- **Dependencies**: None (foundational layer)
- **Technology Specifics**: Echo v4, Huma v2, OpenTelemetry, NATS JetStream, PostgreSQL with GORM

## Business Domain Components

### **Core Business Module** (`project/core/`)
**Primary Responsibility**: Manages foundational business entities (users, companies, contacts, addresses)
- **Key Interfaces**: Activity management, Contact management, Address management
- **Dependencies**: System foundation, Event bus
- **Sub-components**: Activity, Address, Contact, Module, Notification management

### **Stock Management Module** (`project/stock/`)
**Primary Responsibility**: Manages inventory, items, warehouses, and stock transactions
- **Key Interfaces**: Item catalog, Inventory tracking, Warehouse management, Price lists
- **Dependencies**: Core module, Event bus, Accounting integration
- **Sub-components**: Item, Inventory, Warehouse, Price lists, Stock entries

### **Accounting Module** (`project/accounting/`)
**Primary Responsibility**: Manages financial transactions, ledgers, and accounting processes
- **Key Interfaces**: Journal entries, Ledger management, Payment processing, Financial reporting
- **Dependencies**: Core module, Event bus
- **Sub-components**: Ledger, Journal, Payment, Cost centers, Financial reporting
