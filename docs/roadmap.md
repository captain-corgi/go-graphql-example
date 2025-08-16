# Roadmap

A suggested sequence to take this scaffold to a production-ready GraphQL service.

## Phase 1 — Hello GraphQL

- Choose stack (recommended: gqlgen).
- Add `api/graphql/schema.graphqls` with a minimal Query.
- Generate resolvers (gqlgen) and wire to `internal/application` use case that returns static data.
- Add `cmd/server/main.go` to start HTTP server and mount GraphQL endpoint and GraphiQL/Playground.

## Phase 2 — Domain Modeling

- Define core entities and value objects in `internal/domain`.
- Define repository interfaces in `domain`.
- Implement use cases in `internal/application`.

## Phase 3 — Persistence

- Implement in-memory repositories in `internal/infrastructure` to start.
- Introduce a backing store (Postgres or alternative) and migrate repository implementations.
- Add migrations and configuration management.

## Phase 4 — Interfaces & Performance

- Add batching/dataloaders where necessary.
- Introduce pagination, filtering, and error shaping for GraphQL.

## Phase 5 — Observability & Ops

- Structured logging; add request IDs.
- Metrics and tracing with OpenTelemetry.
- Health checks and readiness/liveness endpoints.

## Phase 6 — Packaging & Security

- Add containerization and runtime configuration.
- Input validation, authn/authz at the interfaces layer.

## Ongoing

- Testing strategy: unit tests across layers; integration tests around adapters.
- Documentation: keep `docs/` in sync with code changes.
