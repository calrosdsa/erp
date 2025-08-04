# Coding Standards

## Core Standards
- **Languages & Runtimes:** Go 1.21+, TypeScript 5.0+
- **Style & Linting:** Use `gofmt`, `golangci-lint` for Go; ESLint + Prettier for TypeScript
- **Test Organization:** Go tests in `*_test.go` files alongside source

## Critical Rules
- **Database Access:** All database queries must use the repository pattern with generated GORM models from `gen/db/`
- **API Response Format:** All REST endpoints must use the standard response wrapper
- **Event Publishing:** Domain events must be published through the event bus
- **Module Registration:** New modules must be registered in `cmd/all/main.go`
- **Generated Code:** Never manually edit files in `gen/db/` - always regenerate with `make models`
