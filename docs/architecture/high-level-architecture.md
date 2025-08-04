# High Level Architecture

## Technical Summary

The ERP system employs a domain-driven microservices-within-monolith architecture, combining the operational simplicity of monolithic deployment with the modularity benefits of microservices design. The system leverages Go's excellent concurrency model and the existing NATS JetStream event infrastructure to enable real-time, event-driven communication between business modules (Core, Stock, Accounting, Sales/Buying, Project Management, Document Management). This architecture directly supports the PRD goals of 99.5%+ event delivery reliability while maintaining the 30-day implementation timeline through simplified deployment and unified dependency management.

## High Level Overview

**Architectural Style:** Domain-Driven Microservices within Monolith
The system maintains separate business domains as independent modules within a single deployable unit, providing modularity without distributed system complexity.

**Repository Structure:** Monorepo (as specified in PRD)
All modules reside in the existing `/project/<module>` structure with unified CI/CD, dependency management, and atomic commits across modules.

**Service Architecture:** Modular Monolith with Event-Driven Communication
Each business module operates independently, communicating through NATS JetStream events and accessing shared PostgreSQL database with clear schema boundaries.

**Primary User Flow:** Dashboard-First with Contextual Workflows
Users land on role-specific dashboards showing key metrics and pending actions, then engage with guided business workflows (quote-to-cash, inventory management) optimized for SME operational patterns.

## High Level Project Diagram

```mermaid
graph TB
    subgraph "External Systems"
        E1[E-commerce Platforms<br/>Shopify, WooCommerce]
        E2[Payment Gateways<br/>Stripe, PayPal] 
        E3[Accounting Systems<br/>QuickBooks, Xero]
        E4[Shipping Providers<br/>UPS, FedEx, USPS]
    end

    subgraph "Frontend Layer"
        UI[React TypeScript UI<br/>Tailwind CSS + Headless UI]
    end

    subgraph "API Gateway Layer"
        API[Echo v4 + Huma v2<br/>OpenAPI REST API]
    end

    subgraph "ERP Monolith (Go)"
        subgraph "Business Modules"
            CORE[Core Module<br/>Users, Companies, Auth]
            STOCK[Stock Module<br/>Items, Inventory, Warehouses]
            SALES[Sales Module<br/>Orders, Quotes, Customers]
            ACCT[Accounting Module<br/>Ledgers, Journals, Reports]
            PROJ[Project Module<br/>Tasks, Time Tracking]
            DOCS[Document Module<br/>Files, Terms & Conditions]
        end
        
        subgraph "System Infrastructure"
            AUTH[JWT + Permify<br/>Authentication & Authorization]
            EVENTS[NATS JetStream<br/>Event Bus]
            DI[Dependency Injection<br/>Container]
        end
    end

    subgraph "Data Layer"
        DB[(PostgreSQL<br/>with GORM)]
        CACHE[(Redis<br/>Sessions & Cache)]
        FILES[S3-Compatible<br/>Object Storage]
    end

    %% External integrations
    E1 --> API
    E2 --> API
    E3 --> API
    E4 --> API

    %% Frontend to API
    UI --> API

    %% API to modules
    API --> CORE
    API --> STOCK
    API --> SALES
    API --> ACCT
    API --> PROJ
    API --> DOCS

    %% Module interactions via events
    CORE --> EVENTS
    STOCK --> EVENTS
    SALES --> EVENTS
    ACCT --> EVENTS
    PROJ --> EVENTS
    DOCS --> EVENTS
```

## Architectural and Design Patterns

- **Domain-Driven Design (DDD):** Clear business domain boundaries with separate modules for each domain context
- **Event-Driven Architecture:** NATS JetStream for asynchronous module communication with event sourcing
- **Repository Pattern:** Abstract data access through interfaces with GORM implementation
- **Dependency Injection:** Custom DI container managing service lifecycles and dependencies
- **CQRS (Command Query Responsibility Segregation):** Separate read and write models for complex business operations
- **Finite State Machine (FSM):** Business process flows using existing `pkg/fsm/` infrastructure
