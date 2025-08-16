# Conventions

These conventions help maintain consistency as the codebase grows.

## Package & Imports
- Respect import direction: outer layers depend inward only.
- Keep packages small and focused.
- `internal/` is not importable by other modules; `pkg/` is safe for published utilities (use sparingly).

## Naming
- Packages: short, lower_snake or lower (e.g., `persistence`, `graphql`, `user`).
- Constructors: `NewType(...) *Type`.
- Interfaces: named by role (`UserRepository`, `Clock`).

## Errors
- Return errors rather than panics in app/domain layers.
- Wrap with `%w`; add context at boundaries.
- Consider typed errors in domain for business cases.

## Context & Time
- Pass `context.Context` from handlers into use cases and infra.
- Time abstraction via an interface (e.g., `Clock`) to ease testing.

## Logging
- Structured logging. Avoid logging in domain; log at interfaces/infrastructure.
- Include correlation IDs/request IDs from the transport layer.

## Configuration
- Environment-driven configuration (12-factor style).
- Validate config at startup; fail fast.

## Testing
- Unit tests in all layers; isolate domain tests from infra.
- Use table-driven tests.
- Provide fakes/mocks for repositories.

## GraphQL
- Schema lives in `api/graphql/`.
- Keep resolvers thin; delegate to use cases.
- Use batch loading patterns (e.g., dataloader) when necessary.
