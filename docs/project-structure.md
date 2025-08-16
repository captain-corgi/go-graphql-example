# Project Structure

This layout follows Clean Architecture/Hexagonal principles. It is intentionally empty to start, providing a clear separation of concerns as you add features.

```text
/ (repo root)
├─ .git/
├─ api/
│  └─ graphql/              # GraphQL schema and API surface (SDL, directives, docs)
├─ cmd/
│  └─ server/               # Application entrypoint (e.g., main.go) for the HTTP/GraphQL server
├─ internal/
│  ├─ application/          # Use cases (application services), input/output DTOs
│  ├─ domain/               # Core domain model (entities, value objects, domain services, repository interfaces)
│  ├─ infrastructure/       # Adapters: DB implementations, config, cache, external clients
│  └─ interfaces/           # Delivery layer: GraphQL resolvers/handlers, HTTP routing, middleware
├─ pkg/                     # Optional: shared util packages safe for external consumption
├─ docs/                    # Project documentation (this folder)
└─ go.mod                   # Go module definition
```

## Directory Responsibilities

- **`api/graphql/`**
  - Owns the GraphQL schema (`*.graphqls`, `*.graphql`, `*.gql`).
  - Co-locate API-specific docs, custom directives, and GraphQL boundary concerns.
  - Avoid business logic here; keep it schema-first or schema-owned.
  - File layout conventions:
    - Split queries and mutations: `query.graphqls`, `mutation.graphqls`.
    - Model files per domain: e.g., `user.graphqls`, `post.graphqls`.
    - Keep types close to their domain; avoid monolithic schemas.
  - gqlgen config: `gqlgen.yml` at repo root maps schema to generated code.

- **`cmd/server/`**
  - Hosts the executable entrypoint (e.g., `main.go`).
  - Composition root: wire dependencies, build resolvers, start HTTP server.
  - Keep this thin; orchestration only.

- **`internal/domain/`**
  - The heart of the system (DDD):
    - Entities and Value Objects model the ubiquitous language.
    - Aggregates encapsulate invariants; define repository interfaces per aggregate root.
    - Domain Services capture domain operations not naturally belonging to an entity/VO.
  - Pure Go, no framework/transport details.
  - Must not depend on other `internal/*` layers.

- **`internal/application/`**
  - Use cases that orchestrate domain behavior and enforce application rules.
  - Depends on `domain` for types/interfaces.
  - Exposes ports used by the delivery layer.

- **`internal/infrastructure/`**
  - Implementations of secondary adapters (DB, queues, cache, external APIs), configuration, logging, migrations.
  - Implements `domain` repository interfaces.
  - Depends inward on `domain` (and optionally `application` for DTOs), never the other way around.

- **`internal/interfaces/`**
  - Primary adapters: GraphQL resolvers/handlers, HTTP routing (Gin), middlewares.
  - Translates transport concerns into application use case calls.
  - Resolver guidance: generated resolver stubs must remain thin; delegate to application services only (no business logic in generated files).

- **`pkg/`**
  - Optional packages intended for reuse outside `internal/` (be cautious; keep it small and generic).

- **`docs/`**
  - Living documentation for the project. Start here when onboarding.

- **`go.mod`**
  - Defines the module path and Go version for reproducibility.

## Import Rules (high-level)

- `domain` → imports nothing internal.
- `application` → may import `domain`.
- `interfaces` → may import `application` and `domain`.
- `infrastructure` → may import `domain` (and selectively `application` when necessary).
- `cmd/server` → may import all internal layers to assemble the app.

## Suggested Locations for Common Things

- Configuration loading: `internal/infrastructure/config`.
- Logging setup: `internal/infrastructure/logging`.
- Database adapters: `internal/infrastructure/persistence/sql` (raw `database/sql`, no ORM).
- GraphQL resolvers: `internal/interfaces/graphql` (wire to `api/graphql` schema).
- HTTP router/middleware (Gin): `internal/interfaces/http`.
- Cross-cutting utils (if any): `pkg/`.
