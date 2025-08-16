# Stack Choices

This repo does not yet include implementation code or dependencies beyond `go.mod`. Below are recommended options to choose from before coding.

## GraphQL Server (chosen: gqlgen)

- Schema-first workflow: design schema → generate code → implement services/resolvers.
- Keep schema split by concern under `api/graphql/`:
  - `query.graphqls`, `mutation.graphqls`
  - domain models by file, e.g., `user.graphqls`, `post.graphqls`
- Generated resolver files should be thin: delegate to application services only. Avoid business logic in generated files as they may be overwritten.

## Web Framework (chosen: Gin)

- Use `gin-gonic/gin` for HTTP server, routing, and middleware.
- Mount the GraphQL handler (from gqlgen) on a Gin route (e.g., `POST /query`).
- Keep transport concerns (middleware, authn/z, request IDs) in the interfaces layer.

## Persistence (chosen: database/sql; no ORM)

- Use `database/sql` with a driver (e.g., `pgx` for Postgres) and repository implementations in `internal/infrastructure/persistence/sql`.
- No ORM. Prefer prepared statements and transactions encapsulated in repositories.
- Start with in-memory repos for early development if needed, then switch to SQL.

## Configuration (chosen: Viper)

- Use `spf13/viper` for environment/file-based configuration.
- Provide defaults, load env vars, and validate at startup.
- Keep config types in `internal/infrastructure/config` and inject them at composition root.

## Observability

- Logging: structured logs (`zap`, `zerolog`, or stdlib `log/slog` in modern Go).
- Tracing/Metrics: OpenTelemetry SDK (`go.opentelemetry.io/otel`).

## Validation & Errors

- Validation: `go-playground/validator` for rich validation if needed at the edges.
- Errors: Wrap with context using `%w` and expose typed errors from domain/application.

## Testing

- Unit tests with `testing`.
- Mocks for interfaces using `gomock` + `mockgen`.
- Directive on each interface file to generate mocks for all interfaces in that file into a sibling `mocks/` directory using package `mocks`:

  ```go
  //go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks
  ```

- Run generation across the module:

  ```bash
  go generate ./...
  ```

- Table-driven tests and golden files where helpful.
- See `docs/testing.md` for full conventions.

## Suggested Default Setup

- GraphQL: **gqlgen** (schema-first, codegen)
- Web: **Gin** (routing/middleware)
- Config: **Viper** (env-first, typed config)
- Persistence: **database/sql** (+ driver, no ORM)
- Logging: **slog** or **zap**
- Observability: OpenTelemetry (later phases)
