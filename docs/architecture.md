# Architecture

The project is structured around Clean Architecture to keep business logic independent from frameworks, databases, and transports. We also apply Domain-Driven Design (DDD) to model the domain using Entities, Value Objects, Aggregates, Repositories, and Domain Services.

## Layers

- **Domain (Core)**: Entities, value objects, domain services, and repository interfaces.
- **Application (Use Cases)**: Orchestrates domain operations, defines application services and DTOs.
- **Interfaces (Primary Adapters)**: GraphQL resolvers, HTTP handlers, routing, middleware.
- **Infrastructure (Secondary Adapters)**: DB implementations, caching, external APIs, configuration, logging.

## Dependency Direction

Inner layers must not depend on outer layers. Imports point inward.

```
[ Interfaces ]  →  [ Application ]  →  [ Domain ]
      ↓                           ↑
[ Infrastructure ]  --------------┘ (implements domain ports)
```

- `interfaces` calls `application` use cases.
- `application` uses `domain` types and repository interfaces.
- `infrastructure` implements repository interfaces defined in `domain`.
- `cmd/server` composes/wires dependencies and starts the server.

## Data Flow (GraphQL example)

1. GraphQL resolver (interfaces) receives a query.
2. Resolver calls a use case (application) with DTOs.
3. Use case executes domain logic and accesses repositories through interfaces defined in `domain`.
4. Repository is implemented in `infrastructure` and returns domain models.
5. Use case maps domain models to DTOs for the resolver, which returns GraphQL types.

## Technology Mapping

- Interfaces layer:
  - Web framework: Gin (`gin-gonic/gin`).
  - GraphQL: gqlgen resolvers and HTTP handler.
- Application layer:
  - Use case services invoked by resolvers; contains orchestration and application policies.
- Domain layer:
  - Entities, Value Objects, Aggregates, Repository interfaces, Domain Services.
- Infrastructure layer:
  - Persistence with `database/sql` (no ORM) and drivers (e.g., `pgx`).
  - Configuration via Viper (`spf13/viper`).
  - Logging, external clients, migrations.
- Composition root (`cmd/server`):
  - Build config, DB, repositories, services, and resolvers. Mount gqlgen handler on Gin routes.

## Design Principles

- Keep domain pure and framework-agnostic.
- Prefer constructor injection and small interfaces.
- Avoid cyclic dependencies by enforcing import rules.
- Keep composition at the edges (`cmd/server`).
- Schema-first workflow: design GraphQL schema → generate code (gqlgen) → implement services/resolvers.
- Resolver/service separation: generated resolvers remain thin and delegate to application services only (avoid business logic in generated files since they may be regenerated).

## Testing & Mocks

- Use `gomock` with `mockgen` for interaction tests around `application` and `interfaces` layers.
- For every file that declares one or more interfaces (commonly in `internal/domain` and `internal/application`), add a `go:generate` directive that generates mocks for all interfaces in that file into a sibling `mocks/` directory with package name `mocks`.

  Recommended pattern:

  ```go
  //go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks
  ```

- Regenerate mocks when interfaces change:

  ```bash
  go generate ./...
  ```

- Keep mocks out of production packages to avoid import cycles and maintain separation of concerns.
- No production code is implemented here; this documents conventions only.
