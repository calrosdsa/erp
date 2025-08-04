# External APIs

## Payment Processing APIs

**Square API Integration** (Already partially implemented)
- **Purpose**: Credit card payment processing for sales orders
- **Authentication**: OAuth 2.0 with application credentials
- **Key Endpoints**: `/v2/payments`, `/v2/orders`, `/v2/payments/{payment_id}`

**Stripe API Integration** (Required by PRD)
- **Purpose**: Alternative payment processor with subscription support
- **Authentication**: API key-based authentication
- **Key Endpoints**: `/v1/payments/intents`, `/v1/customers`, `/v1/subscriptions`

## E-commerce Platform APIs

**Shopify Admin API** (Required by PRD)
- **Purpose**: Product catalog sync and order import from Shopify stores
- **Authentication**: OAuth 2.0 with shop-specific access tokens
- **Key Endpoints**: `/admin/api/2023-10/products.json`, `/admin/api/2023-10/orders.json`

**WooCommerce REST API** (Required by PRD)
- **Purpose**: WordPress e-commerce integration
- **Authentication**: OAuth 1.0a or Basic Auth with consumer key/secret
- **Key Endpoints**: `/wp-json/wc/v3/products`, `/wp-json/wc/v3/orders`
