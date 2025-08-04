# Data Models

## Core Data Models

### **Company**
**Purpose:** Central entity representing business organizations in a multi-tenant hierarchy

**Key Attributes:**
- `id`: int64 - Primary key for internal references
- `uuid`: string - External identifier for API interactions
- `name`: string - Company display name
- `code`: string - Unique company identifier/code
- `is_group`: bool - Indicates if this is a parent company
- `parent_id`: *int64 - Self-referencing hierarchy support

**Relationships:**
- Self-referencing parent-child hierarchy for company groups
- One-to-many with Users, Items, Orders, and all business entities

### **User**
**Purpose:** Authentication and user management entity

**Key Attributes:**
- `id`: int64 - Primary key
- `uuid`: string - External user identifier
- `identifier`: string - Username/email for login
- `password_hash`: string - Encrypted password storage

**Relationships:**
- Many-to-many with Companies through workspace access
- One-to-many with Activities, Orders, and user-generated content

### **Item**
**Purpose:** Product/service catalog entity supporting complex inventory structures

**Key Attributes:**
- `id`: int64 - Primary key
- `uuid`: string - External item identifier  
- `name`: string - Item display name
- `code`: *string - SKU or item code
- `item_type`: string - ITEM/SERVICE/BUNDLE classification
- `maintain_stock`: bool - Inventory tracking flag

**Relationships:**
- Belongs to Company
- Self-referencing for variants and bundles
- One-to-many with ItemLines, StockEntries, Prices

### **Order**
**Purpose:** Sales and purchase order processing with full ERP workflow support

**Key Attributes:**
- `id`: int64 - Primary key
- `code`: string - Order number/reference
- `company_id`: int64 - Company ownership
- `party_id`: int64 - Customer/supplier reference
- `currency`: string - Transaction currency
- `status`: string - Order workflow state

**Relationships:**
- Belongs to Company and Party (customer/supplier)
- One-to-many with OrderLines (line items)
- Links to Invoices and StockEntries for fulfillment
