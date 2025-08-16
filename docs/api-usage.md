# API Usage Guide

This guide provides comprehensive documentation on how to use the GraphQL API, including authentication, error handling, and best practices.

## Getting Started

### Starting the Server

1. Ensure PostgreSQL is running and configured
2. Run database migrations:

   ```bash
   go run cmd/server/main.go -migrate
   ```

3. Start the server:

   ```bash
   go run cmd/server/main.go
   ```

4. Access GraphQL Playground at `http://localhost:8080/playground`

### Basic Configuration

The server can be configured using environment variables or configuration files. See `configs/README.md` for details.

## GraphQL Endpoint

- **URL**: `http://localhost:8080/query`
- **Method**: POST
- **Content-Type**: `application/json`

## Schema Overview

The API provides a User management system with the following operations:

### Queries

- `user(id: ID!)` - Get a single user by ID
- `users(first: Int, after: String)` - Get paginated list of users

### Mutations

- `createUser(input: CreateUserInput!)` - Create a new user
- `updateUser(id: ID!, input: UpdateUserInput!)` - Update an existing user
- `deleteUser(id: ID!)` - Delete a user

## Data Types

### User

```graphql
type User {
  id: ID!
  email: String!
  name: String!
  createdAt: String!
  updatedAt: String!
}
```

### Input Types

```graphql
input CreateUserInput {
  email: String!
  name: String!
}

input UpdateUserInput {
  email: String
  name: String
}
```

### Response Types

```graphql
type CreateUserPayload {
  user: User!
  errors: [Error!]
}

type UpdateUserPayload {
  user: User!
  errors: [Error!]
}

type DeleteUserPayload {
  success: Boolean!
  errors: [Error!]
}
```

## Pagination

The API uses cursor-based pagination for the `users` query:

```graphql
type UserConnection {
  edges: [UserEdge!]!
  pageInfo: PageInfo!
}

type UserEdge {
  node: User!
  cursor: String!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}
```

### Pagination Example

```graphql
# First page
query {
  users(first: 10) {
    edges {
      node { id email name }
      cursor
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}

# Next page using cursor
query {
  users(first: 10, after: "cursor_from_previous_query") {
    edges {
      node { id email name }
      cursor
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
  }
}
```

## Error Handling

### Error Structure

```graphql
type Error {
  message: String!
  field: String
  code: String
}
```

### Common Error Codes

| Code | Description | Field |
|------|-------------|-------|
| `USER_NOT_FOUND` | User with specified ID doesn't exist | - |
| `INVALID_EMAIL` | Email format is invalid | `email` |
| `DUPLICATE_EMAIL` | Email already exists | `email` |
| `VALIDATION_ERROR` | General validation error | varies |
| `INTERNAL_ERROR` | Server error | - |

### Error Response Example

```json
{
  "data": {
    "createUser": {
      "user": null,
      "errors": [
        {
          "message": "Invalid email format",
          "field": "email",
          "code": "INVALID_EMAIL"
        }
      ]
    }
  }
}
```

## Validation Rules

### Email Validation

- Must be a valid email format
- Must be unique across all users
- Maximum length: 255 characters

### Name Validation

- Required field
- Minimum length: 1 character
- Maximum length: 255 characters
- Cannot be only whitespace

## Rate Limiting

Currently, no rate limiting is implemented. In production environments, consider implementing rate limiting at the reverse proxy level or using middleware.

## Authentication

Authentication is not currently implemented. Future versions will include:

- JWT-based authentication
- Role-based access control
- API key authentication

## Best Practices

### Query Optimization

1. **Request only needed fields**: GraphQL allows you to specify exactly which fields you need
2. **Use pagination**: Always use pagination for list queries to avoid large responses
3. **Batch operations**: Use GraphQL's ability to send multiple operations in a single request

### Error Handling

1. **Check both data and errors**: GraphQL can return partial data with errors
2. **Handle network errors**: Implement proper error handling for network failures
3. **Validate input client-side**: Reduce server load by validating input before sending

### Performance

1. **Use variables**: Parameterize queries using variables instead of string concatenation
2. **Cache responses**: Implement caching strategies for frequently accessed data
3. **Monitor query complexity**: Be aware of query depth and complexity

## Example Implementations

### JavaScript/Node.js

```javascript
const query = `
  query GetUser($id: ID!) {
    user(id: $id) {
      id
      email
      name
    }
  }
`;

const response = await fetch('http://localhost:8080/query', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    query,
    variables: { id: 'user-id' }
  })
});

const result = await response.json();
```

### Python

```python
import requests

query = """
  query GetUser($id: ID!) {
    user(id: $id) {
      id
      email
      name
    }
  }
"""

response = requests.post(
    'http://localhost:8080/query',
    json={
        'query': query,
        'variables': {'id': 'user-id'}
    }
)

result = response.json()
```

### cURL

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query GetUser($id: ID!) { user(id: $id) { id email name } }",
    "variables": { "id": "user-id" }
  }'
```

## Testing

### Development Data

Use the provided seed data for testing:

```bash
psql -d graphql_service_dev -f scripts/seed-dev-data.sql
```

### Example Queries

See `examples/graphql/` directory for comprehensive query and mutation examples.

### Integration Testing

The API includes comprehensive integration tests. Run them with:

```bash
go test ./internal/interfaces/graphql/...
```

## Monitoring and Observability

### Logging

- All requests are logged with correlation IDs
- Errors include stack traces in development mode
- Structured JSON logging in production

### Metrics

Future versions will include:

- Request/response metrics
- Error rate monitoring
- Performance metrics
- Database query metrics

### Health Checks

A health check endpoint will be added in future versions at `/health`.

## Troubleshooting

### Common Issues

1. **Connection refused**: Ensure the server is running on the correct port
2. **Database connection errors**: Check database configuration and connectivity
3. **Invalid queries**: Use GraphQL Playground to validate query syntax
4. **CORS errors**: Configure CORS middleware for browser requests

### Debug Mode

Enable debug logging by setting:

```bash
GRAPHQL_SERVICE_LOGGING_LEVEL=debug
```

### Database Issues

Check database connectivity:

```bash
psql -d graphql_service_dev -c "SELECT version();"
```

For more detailed troubleshooting, check the application logs and ensure all dependencies are properly configured.
