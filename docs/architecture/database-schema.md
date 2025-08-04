# Database Schema

## PostgreSQL Schema Design

```sql
-- Core System Tables
CREATE TABLE companies (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    is_group BOOLEAN DEFAULT FALSE,
    parent_id BIGINT REFERENCES companies(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_companies_parent_id ON companies(parent_id);
CREATE INDEX idx_companies_code ON companies(code);
CREATE INDEX idx_companies_deleted_at ON companies(deleted_at);

-- User Management
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    identifier VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    last_login TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_identifier ON users(identifier);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Party Management (Universal Customer/Supplier Model)
CREATE TABLE parties (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    party_type VARCHAR(20) NOT NULL CHECK (party_type IN ('CUSTOMER', 'SUPPLIER', 'BOTH')),
    company_id BIGINT NOT NULL REFERENCES companies(id),
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_parties_company_id ON parties(company_id);
CREATE INDEX idx_parties_type ON parties(party_type);
CREATE INDEX idx_parties_status ON parties(status);

-- Stock Management
CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NULL,
    company_id BIGINT NOT NULL REFERENCES companies(id),
    group_id BIGINT NULL REFERENCES groups(id),
    parent_id BIGINT NULL REFERENCES items(id),
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('ITEM', 'SERVICE', 'BUNDLE')),
    maintain_stock BOOLEAN DEFAULT TRUE,
    unit_of_measure_id BIGINT REFERENCES unit_of_measures(id),
    status VARCHAR(20) DEFAULT 'ENABLED',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_items_company_id ON items(company_id);
CREATE INDEX idx_items_code ON items(code);
CREATE INDEX idx_items_type ON items(item_type);
CREATE INDEX idx_items_parent_id ON items(parent_id);

-- Sales & Purchase Management
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL,
    company_id BIGINT NOT NULL REFERENCES companies(id),
    party_id BIGINT NOT NULL REFERENCES parties(id),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    posting_date DATE NOT NULL DEFAULT CURRENT_DATE,
    delivery_date DATE NULL,
    project_id BIGINT NULL REFERENCES projects(id),
    cost_center_id BIGINT NULL REFERENCES cost_centers(id),
    price_list_id BIGINT NULL REFERENCES price_lists(id),
    total_amount DECIMAL(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_orders_company_id ON orders(company_id);
CREATE INDEX idx_orders_party_id ON orders(party_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_posting_date ON orders(posting_date);

-- Accounting Schema
CREATE TABLE ledger_accounts (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    company_id BIGINT NOT NULL REFERENCES companies(id),
    parent_id BIGINT NULL REFERENCES ledger_accounts(id),
    account_type VARCHAR(50) NOT NULL,
    is_group BOOLEAN DEFAULT FALSE,
    balance DECIMAL(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(company_id, code)
);

CREATE INDEX idx_ledger_accounts_company_id ON ledger_accounts(company_id);
CREATE INDEX idx_ledger_accounts_parent_id ON ledger_accounts(parent_id);
CREATE INDEX idx_ledger_accounts_type ON ledger_accounts(account_type);
```
