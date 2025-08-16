# Go GraphQL Example — Project Documentation

This repository is a scaffold for a GraphQL service written in Go, organized using Clean Architecture principles. It currently contains only the directory layout (no implementation code yet) and a Go module definition.

- Module: `github.com/captain-corgi/go-graphql-example`
- Go version: `1.24.3`

Use this documentation as the canonical guide for structure, architecture, conventions, and development workflow as the project grows.

## Table of Contents

- [Project Structure](./project-structure.md)
- [Architecture](./architecture.md)
- [Stack Choices](./stack.md)
- [Conventions](./conventions.md)
- [Development Guide](./development.md)
- [API Usage Guide](./api-usage.md)
- [Architecture Decisions](./architecture-decisions.md)
- [Roadmap](./roadmap.md)

## Current Status

- **Complete**: GraphQL service foundation with Clean Architecture implementation
- **Features**: User management API with CRUD operations, pagination, and error handling
- **Infrastructure**: Database connectivity, configuration management, structured logging
- **Testing**: Comprehensive unit and integration tests with mock generation
- **Documentation**: Complete API documentation and architecture decisions

## Quick Start

1. **Setup Database**: Create PostgreSQL database and run migrations
2. **Seed Data**: Load development data using `scripts/seed-dev-data.sql`
3. **Start Server**: Run `go run cmd/server/main.go`
4. **Explore API**: Visit `http://localhost:8080/playground` for GraphQL Playground
5. **Try Examples**: Use queries and mutations from `examples/graphql/`

## Key Resources

- [API Usage Guide](./api-usage.md) for comprehensive API documentation
- [Architecture Decisions](./architecture-decisions.md) for design rationale
- [Development Guide](./development.md) for setup and workflow
- [Project Structure](./project-structure.md) for codebase organization
