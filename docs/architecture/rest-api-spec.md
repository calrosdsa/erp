# REST API Spec

## OpenAPI 3.0 Specification

```yaml
openapi: 3.0.0
info:
  title: Modular ERP System API
  version: 1.0.0
  description: Comprehensive API for modular ERP system supporting multi-tenant operations
servers:
  - url: https://api.erp-system.com/v1
    description: Production server
  - url: https://staging-api.erp-system.com/v1
    description: Staging server

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    Company:
      type: object
      required:
        - name
        - code
      properties:
        id:
          type: integer
          format: int64
          readOnly: true
        uuid:
          type: string
          format: uuid
          readOnly: true
        name:
          type: string
          minLength: 1
          maxLength: 255
        code:
          type: string
          minLength: 1
          maxLength: 50
        isGroup:
          type: boolean
          default: false
        parentId:
          type: integer
          format: int64
          nullable: true

    Item:
      type: object
      required:
        - name
        - itemType
      properties:
        id:
          type: integer
          format: int64
          readOnly: true
        uuid:
          type: string
          format: uuid
          readOnly: true
        name:
          type: string
          minLength: 1
          maxLength: 255
        code:
          type: string
          maxLength: 50
          nullable: true
        itemType:
          type: string
          enum: [ITEM, SERVICE, BUNDLE]
        maintainStock:
          type: boolean
          default: true

security:
  - BearerAuth: []

paths:
  /auth/sign-in:
    post:
      summary: User authentication
      security: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - identifier
                - password
              properties:
                identifier:
                  type: string
                password:
                  type: string
      responses:
        '200':
          description: Authentication successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  accessToken:
                    type: string
                  refreshToken:
                    type: string
                  user:
                    $ref: '#/components/schemas/User'

  /companies:
    get:
      summary: List companies
      responses:
        '200':
          description: List of companies
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Company'
    post:
      summary: Create new company
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Company'
      responses:
        '201':
          description: Company created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Company'

  /items:
    get:
      summary: List items
      responses:
        '200':
          description: List of items
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Item'
    post:
      summary: Create new item
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Item'
      responses:
        '201':
          description: Item created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Item'
```
