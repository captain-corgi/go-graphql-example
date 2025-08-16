# GraphQL Examples

This directory contains example GraphQL queries and mutations that demonstrate how to use the API.

## Files

- `queries.graphql` - Example queries for retrieving user data
- `mutations.graphql` - Example mutations for creating, updating, and deleting users

## Using the Examples

### GraphQL Playground

1. Start the server:

   ```bash
   go run cmd/server/main.go
   ```

2. Open GraphQL Playground in your browser:

   ```
   http://localhost:8080/playground
   ```

3. Copy and paste examples from the `.graphql` files into the playground
4. Use the variables panel for parameterized queries

### cURL Examples

You can also test the API using cURL:

```bash
# Query a user
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { user(id: \"550e8400-e29b-41d4-a716-446655440001\") { id email name } }"
  }'

# Create a user
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { createUser(input: { email: \"test@example.com\", name: \"Test User\" }) { user { id email name } errors { message } } }"
  }'
```

### Using Variables

Many examples include variable definitions. In GraphQL Playground:

1. Copy the query/mutation to the main panel
2. Copy the variables JSON to the variables panel (bottom left)
3. Execute the operation

Example with variables:

```graphql
# Query
query GetUser($userId: ID!) {
  user(id: $userId) {
    id
    email
    name
  }
}

# Variables (in variables panel)
{
  "userId": "550e8400-e29b-41d4-a716-446655440001"
}
```

## Sample Data

The examples reference sample user IDs from the development seed data:

- `550e8400-e29b-41d4-a716-446655440001` - John Doe
- `550e8400-e29b-41d4-a716-446655440002` - Jane Smith
- `550e8400-e29b-41d4-a716-446655440003` - Bob Wilson
- And more...

Make sure to seed your development database first:

```bash
psql -d graphql_service_dev -f scripts/seed-dev-data.sql
```

## Error Handling

The API returns errors in a structured format:

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

Common error codes:

- `USER_NOT_FOUND` - User with specified ID doesn't exist
- `INVALID_EMAIL` - Email format is invalid
- `DUPLICATE_EMAIL` - Email already exists
- `VALIDATION_ERROR` - General validation error

## Pagination

The `users` query supports cursor-based pagination:

1. First request: `users(first: 5)`
2. Get the `endCursor` from `pageInfo`
3. Next request: `users(first: 5, after: "cursor_value")`
4. Continue until `hasNextPage` is false

See `queries.graphql` for detailed pagination examples.
