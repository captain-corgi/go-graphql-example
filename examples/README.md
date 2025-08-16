# Examples

This directory contains practical examples and sample data for the GraphQL service.

## Directory Structure

```
examples/
├── graphql/           # GraphQL query and mutation examples
│   ├── queries.graphql    # Example queries for data retrieval
│   ├── mutations.graphql  # Example mutations for data modification
│   └── README.md         # GraphQL examples documentation
└── README.md         # This file
```

## Getting Started

### Prerequisites

1. **Database Setup**: Ensure PostgreSQL is running and configured
2. **Migrations**: Run database migrations to create tables
3. **Seed Data**: Load sample data for testing

```bash
# Create and setup database
createdb graphql_service_dev

# Run migrations
go run cmd/server/main.go -migrate

# Load seed data
psql -d graphql_service_dev -f scripts/seed-dev-data.sql

# Start the server
go run cmd/server/main.go
```

### Using the Examples

1. **GraphQL Playground**: Open `http://localhost:8080/playground` in your browser
2. **Copy Examples**: Copy queries/mutations from `graphql/` directory
3. **Execute**: Run the examples in the playground
4. **Modify**: Adapt examples for your specific use cases

## Sample Data

The examples use sample user data with the following IDs:

| ID | Email | Name |
|----|-------|------|
| `550e8400-e29b-41d4-a716-446655440001` | <john.doe@example.com> | John Doe |
| `550e8400-e29b-41d4-a716-446655440002` | <jane.smith@example.com> | Jane Smith |
| `550e8400-e29b-41d4-a716-446655440003` | <bob.wilson@example.com> | Bob Wilson |
| `550e8400-e29b-41d4-a716-446655440004` | <alice.johnson@example.com> | Alice Johnson |
| `550e8400-e29b-41d4-a716-446655440005` | <charlie.brown@example.com> | Charlie Brown |

*See `scripts/seed-dev-data.sql` for the complete list of sample users.*

## Example Categories

### Basic Operations

- Single user queries
- User creation, updates, and deletion
- Error handling examples

### Advanced Features

- Pagination with cursors
- Parameterized queries with variables
- Batch operations

### Error Scenarios

- Invalid input validation
- Non-existent resource handling
- Duplicate data conflicts

## Integration Examples

### cURL

```bash
# Query example
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "{ users(first: 5) { edges { node { id email name } } } }"}'
```

### JavaScript/Node.js

```javascript
const response = await fetch('http://localhost:8080/query', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    query: '{ users(first: 5) { edges { node { id email name } } } }'
  })
});
```

### Python

```python
import requests

response = requests.post(
    'http://localhost:8080/query',
    json={'query': '{ users(first: 5) { edges { node { id email name } } } }'}
)
```

## Best Practices Demonstrated

1. **Field Selection**: Examples show how to request only needed fields
2. **Variable Usage**: Parameterized queries for reusability
3. **Error Handling**: Proper error checking and response handling
4. **Pagination**: Cursor-based pagination for large datasets
5. **Input Validation**: Examples of both valid and invalid inputs

## Troubleshooting

### Common Issues

1. **No Data Returned**: Ensure seed data is loaded
2. **Connection Errors**: Verify server is running on correct port
3. **Database Errors**: Check database connectivity and migrations
4. **Invalid Queries**: Use GraphQL Playground for syntax validation

### Debug Tips

1. **Enable Debug Logging**: Set `GRAPHQL_SERVICE_LOGGING_LEVEL=debug`
2. **Check Server Logs**: Monitor console output for errors
3. **Validate Schema**: Use GraphQL Playground's schema explorer
4. **Test Connectivity**: Use simple queries first, then complex ones

## Contributing Examples

When adding new examples:

1. **Document Purpose**: Clearly explain what the example demonstrates
2. **Include Variables**: Provide variable examples where applicable
3. **Show Expected Output**: Include sample responses
4. **Handle Errors**: Show both success and error scenarios
5. **Update Documentation**: Keep README files current

## Related Documentation

- [API Usage Guide](../docs/api-usage.md) - Comprehensive API documentation
- [Development Guide](../docs/development.md) - Setup and development workflow
- [Architecture Decisions](../docs/architecture-decisions.md) - Design rationale
- [GraphQL Examples](./graphql/README.md) - Detailed GraphQL examples
