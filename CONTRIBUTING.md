# Contributing to Go GraphQL Example

Thank you for your interest in contributing to this project! We welcome contributions from the community and are pleased to have you join us.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)

## Code of Conduct

This project adheres to a code of conduct that we expect all contributors to follow. Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

### Prerequisites

- Go 1.24.3 or later
- Docker and Docker Compose
- Make (recommended)
- Git

### Setting Up Your Development Environment

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:

   ```bash
   git clone https://github.com/YOUR_USERNAME/go-graphql-example.git
   cd go-graphql-example
   ```

3. **Add the original repository as upstream**:

   ```bash
   git remote add upstream https://github.com/captain-corgi/go-graphql-example.git
   ```

4. **Set up the development environment**:

   ```bash
   make setup
   ```

5. **Start the development services**:

   ```bash
   make docker-run
   make migrate-up
   ```

6. **Verify everything works**:

   ```bash
   make test
   make dev
   ```

## Development Workflow

### Branch Strategy

- `main` - Production-ready code
- `develop` - Integration branch for features
- `feature/*` - New features
- `bugfix/*` - Bug fixes
- `hotfix/*` - Critical production fixes

### Making Changes

1. **Create a feature branch**:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards

3. **Add tests** for new functionality

4. **Run the test suite**:

   ```bash
   make check
   ```

5. **Commit your changes**:

   ```bash
   git add .
   git commit -m "feat: add amazing new feature"
   ```

6. **Push to your fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request** on GitHub

### Keeping Your Fork Updated

```bash
git checkout main
git fetch upstream
git merge upstream/main
git push origin main
```

## Coding Standards

### Go Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` for formatting (run `make format`)
- Use `go vet` for static analysis (run `make vet`)
- Use `golangci-lint` for comprehensive linting (run `make lint`)

### Architecture Guidelines

#### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`)
   - Pure business logic, no external dependencies
   - Entities, Value Objects, Aggregates
   - Repository interfaces (ports)
   - Domain services

2. **Application Layer** (`internal/application/`)
   - Use cases and application services
   - DTOs for data transfer
   - Orchestrates domain operations
   - May only import `domain` layer

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - Database implementations
   - Configuration, logging, external services
   - Implements domain interfaces
   - May import `domain` and selectively `application`

4. **Interfaces Layer** (`internal/interfaces/`)
   - GraphQL resolvers, HTTP handlers
   - Transport-specific concerns
   - May import `application` and `domain`

#### Dependency Rules

```
interfaces → application → domain
     ↓            ↑
infrastructure ----┘
```

- Inner layers must not depend on outer layers
- Use dependency injection for testability
- Keep resolvers thin - delegate to application services

### Naming Conventions

#### Packages

- Short, descriptive names: `user`, `config`, `persistence`
- Avoid stuttering: `user.User` not `user.UserEntity`

#### Files

- Group related functionality: `user.go`, `user_test.go`
- Repository implementations: `user_repository.go`
- Service implementations: `user_service.go`

#### Interfaces

- Named by capability: `UserRepository`, `UserService`
- Keep interfaces small and focused

### GraphQL Guidelines

#### Schema Design

- Split schema by concern: `query.graphqls`, `mutation.graphqls`
- Domain-specific files: `user.graphqls`, `post.graphqls`
- Use meaningful names and descriptions
- Follow GraphQL best practices for pagination (cursor-based)

#### Resolver Implementation

- Keep resolvers thin - delegate to application services
- Handle errors gracefully with proper GraphQL error responses
- Use context for request-scoped data
- Implement proper authorization checks

### Error Handling

- Return errors rather than panics in application/domain layers
- Wrap errors with context using `fmt.Errorf` with `%w` verb
- Use typed errors in domain for business cases
- Log errors at appropriate levels

### Logging

- Use structured logging with `slog`
- Include correlation IDs/request IDs
- Log at appropriate levels (debug, info, warn, error)
- Avoid logging sensitive information

## Testing Guidelines

### Test Structure

- Test files alongside source: `*_test.go`
- Mocks in `mocks/` subdirectories with package name `mocks`
- Integration tests use build tags: `//go:build integration`

### Writing Tests

#### Unit Tests

```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    service := user.NewService(mockRepo, slog.Default())
    
    // Act & Assert
    // ... test implementation
}
```

#### Integration Tests

```go
//go:build integration

func TestUserRepository_Integration(t *testing.T) {
    // Use real database connection
    // Test actual database operations
}
```

#### Table-Driven Tests

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"invalid email", "invalid", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... test implementation
        })
    }
}
```

### Mock Generation

Add this directive to interface files:

```go
//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks
```

Regenerate mocks when interfaces change:

```bash
make generate-mocks
```

### Test Coverage

- Aim for high test coverage (>80%)
- Focus on critical business logic
- Test error paths and edge cases
- Use `make coverage` to generate coverage reports

## Documentation

### Code Documentation

- Document all public functions and types
- Use Go doc conventions
- Include examples for complex functionality
- Keep comments up to date with code changes

### API Documentation

- Update GraphQL schema descriptions
- Add examples to `examples/graphql/`
- Update API usage guide in `docs/api-usage.md`

### Architecture Documentation

- Update architecture decisions in `docs/architecture-decisions.md`
- Document significant design changes
- Update project structure documentation

## Pull Request Process

### Before Submitting

1. **Ensure all tests pass**:

   ```bash
   make check
   ```

2. **Update documentation** if needed

3. **Add or update tests** for new functionality

4. **Follow commit message conventions**:

   ```
   type(scope): description
   
   - feat: new feature
   - fix: bug fix
   - docs: documentation changes
   - style: formatting changes
   - refactor: code refactoring
   - test: adding tests
   - chore: maintenance tasks
   ```

### Pull Request Template

When creating a PR, please include:

- **Description**: What changes were made and why
- **Type of Change**: Feature, bug fix, documentation, etc.
- **Testing**: How the changes were tested
- **Checklist**: Confirm all requirements are met

### Review Process

1. **Automated Checks**: CI/CD pipeline runs tests and linting
2. **Code Review**: At least one maintainer reviews the code
3. **Discussion**: Address any feedback or questions
4. **Approval**: Maintainer approves the changes
5. **Merge**: Changes are merged into the target branch

## Issue Reporting

### Bug Reports

When reporting bugs, please include:

- **Description**: Clear description of the issue
- **Steps to Reproduce**: Detailed steps to reproduce the bug
- **Expected Behavior**: What should happen
- **Actual Behavior**: What actually happens
- **Environment**: Go version, OS, etc.
- **Logs**: Relevant log output or error messages

### Feature Requests

When requesting features, please include:

- **Description**: Clear description of the feature
- **Use Case**: Why this feature would be useful
- **Proposed Solution**: How you think it should work
- **Alternatives**: Other solutions you've considered

### Issue Labels

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `documentation` - Improvements to documentation
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention is needed
- `question` - Further information is requested

## Getting Help

If you need help or have questions:

- 📖 Check the [documentation](docs/)
- 🔍 Search existing [issues](https://github.com/captain-corgi/go-graphql-example/issues)
- 💬 Start a [discussion](https://github.com/captain-corgi/go-graphql-example/discussions)
- 📧 Contact the maintainers

## Recognition

Contributors will be recognized in:

- GitHub contributors list
- Release notes for significant contributions
- Special mentions for outstanding contributions

Thank you for contributing to Go GraphQL Example! 🚀
