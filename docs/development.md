# Development Guide

This repository currently contains only structure and `go.mod`. Follow these steps to bootstrap and work efficiently.

## Prerequisites

- Go `1.24.3` (as defined in `go.mod`).

## Bootstrapping (Chosen: gqlgen, schema-first)

1. Create schema files under `api/graphql/` (split by concern):
   - `api/graphql/query.graphqls`
   - `api/graphql/mutation.graphqls`
   - Domain models per file, e.g., `api/graphql/user.graphqls`

   Example minimal schema:

   ```graphql
   # api/graphql/query.graphqls
   type Query {
     hello: String!
   }
   ```

2. Add gqlgen and generate code (see gqlgen docs for latest version):

   ```bash
   go run github.com/99designs/gqlgen init
   # or, if already set up
   go run github.com/99designs/gqlgen generate
   ```

3. Configure `gqlgen.yml` (schema glob and output locations):

   ```yaml
   # gqlgen.yml (example)
   schema:
     - api/graphql/**/*.graphqls
   resolver:
     layout: follow-schema
     dir: internal/interfaces/graphql/resolver
     package: resolver
   exec:
     filename: internal/interfaces/graphql/generated.go
     package: graphql
   model:
     filename: internal/interfaces/graphql/models_gen.go
     package: graphql
   ```

   Notes:
   - Keep generated files separate from hand-written code.
   - Do not put business logic in generated files; keep resolvers thin.

4. Implement resolvers that delegate to application services:

   - Place resolver implementations in `internal/interfaces/graphql/resolver/`.
   - Resolvers should call use cases in `internal/application/` only.
   - Keep resolver logic minimal to avoid regeneration conflicts.

5. Create `cmd/server/main.go` with Gin and mount the gqlgen handler:

   ```go
   package main

   import (
     "log"
     "github.com/gin-gonic/gin"
     "github.com/99designs/gqlgen/graphql/handler"
     "github.com/99designs/gqlgen/graphql/playground"
     graph "github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql"
   )

   func main() {
     r := gin.Default()

     srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver()}))
     r.POST("/query", gin.WrapH(srv))
     r.GET("/playground", gin.WrapH(playground.Handler("GraphQL", "/query")))

     if err := r.Run(); err != nil { // reads port from env if set
       log.Fatal(err)
     }
   }
   ```

6. Configuration via Viper (`spf13/viper`):

   - Create `internal/infrastructure/config/config.go` with a typed config struct (e.g., `Port`, `DatabaseURL`).
   - Load defaults, then env vars, then optional file. Validate at startup.
   - Inject config into composition root (in `cmd/server`).

7. Persistence with `database/sql` (no ORM):

   - Define repository interfaces in `internal/domain` per aggregate.
   - Implement SQL repositories in `internal/infrastructure/persistence/sql` using `database/sql` and a driver (e.g., `pgx`).
   - Use context-aware queries, prepared statements, and transactions.

## Workflow Summary

Design schema → generate code (gqlgen) → implement application services and thin resolvers → wire with Gin and Viper → add SQL repositories.

## Recommended Project Wiring

- `cmd/server/main.go`: composition root (construct repos, use cases, resolvers, router; start server).
- `internal/domain`: pure domain types and repository interfaces.
- `internal/application`: use cases invoked by resolvers.
- `internal/infrastructure`: repository implementations, config, logging.
- `internal/interfaces`: GraphQL resolvers/handlers, HTTP wiring.

## Running & Testing

- Build/run once `main.go` exists under `cmd/server/`:

  ```bash
  go run ./cmd/server
  ```

- Run tests when they are added:

  ```bash
  go test ./...
  ```

## Testing & Mocks

- Use `gomock` and `mockgen` for interaction tests.
- For each interface file, add a directive to generate mocks for all interfaces in that file into a sibling `mocks/` directory with package name `mocks`:

  ```go
  //go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks
  ```

- Generate or update all mocks whenever interfaces change:

  ```bash
  go generate ./...
  ```

- See `docs/testing.md` for full conventions and examples.

## Style & Checks

- Use `go fmt`, `go vet`.
- Consider adding linters (e.g., `golangci-lint`) later.

## Next Steps

- See the [Roadmap](./roadmap.md) to prioritize initial features and infrastructure.
