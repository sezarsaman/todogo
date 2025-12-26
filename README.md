# todogo - Task Manager API with Observability

A simple yet robust Task Management API built with Go, featuring comprehensive observability capabilities including Prometheus metrics and request tracing.

## Features

- **CRUD Operations**: Create, Read, Update, Delete tasks
- **Caching**: Redis-based caching for improved performance
- **Database**: PostgreSQL for persistent storage
- **Observability**:
  - Prometheus metrics (requests_total, request_latency_histogram, tasks_count)
  - Request tracing with unique trace IDs
  - Metrics endpoint for Prometheus scraping
- **API Documentation**: Swagger/OpenAPI documentation
- **Health Checks**: Built-in health check endpoint
- **CORS Support**: Cross-Origin Resource Sharing enabled

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.25+ (for local development)

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# Services will be available at:
# - API: http://localhost:8080
# - Prometheus: http://localhost:9090
# - PostgreSQL: localhost:5432
# - Redis: localhost:6379
```

### Local Development

```bash
# Load environment variables
source .env

# Build the application
go build ./cmd/api

# Run the application
./api
```

## API Endpoints

### Tasks

- `POST /tasks` - Create a new task
- `GET /tasks` - List all tasks
- `GET /tasks/:id` - Get a specific task
- `PUT /tasks/:id` - Update a task
- `DELETE /tasks/:id` - Delete a task

### Observability

- `GET /metrics` - Prometheus metrics endpoint
- `GET /health` - Health check endpoint

### Documentation

- `GET /swagger/index.html` - API documentation
- `GET /docs` - Redirect to API documentation

## Environment Variables

Create a `.env` file with the following variables:

```env
HTTP_PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskdb
REDIS_HOST=redis
REDIS_PORT=6379
```

## Observability Guide

For detailed information on metrics and tracing, see [OBSERVABILITY.md](./OBSERVABILITY.md).

### Quick Metrics Examples

View all metrics:
```bash
curl http://localhost:8080/metrics
```

Access Prometheus UI:
```bash
open http://localhost:9090
```

Get trace ID from response:
```bash
curl -i http://localhost:8080/tasks | grep X-Trace-ID
```

### Key Metrics

1. **task_manager_http_requests_total** - Total HTTP requests
2. **task_manager_http_request_latency_seconds** - Request latency histogram
3. **task_manager_tasks_count** - Current task count by status

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test ./... -v
```

Run specific test:
```bash
go test ./internal/task/handler -v
```

## Architecture

```
todogo/
├── cmd/
│   └── api/          # Application entrypoint
├── internal/
│   ├── config/       # Configuration management
│   ├── observability/
│   │   ├── metrics/  # Prometheus metrics
│   │   ├── middleware/ # HTTP middleware
│   │   └── tracing/  # Request tracing
│   ├── platform/
│   │   ├── cache/    # Redis cache layer
│   │   ├── database/ # PostgreSQL connection
│   │   └── http/     # HTTP routing & setup
│   └── task/
│       ├── handler/  # HTTP handlers
│       ├── service/  # Business logic
│       ├── repository/ # Data access
│       ├── cache/    # Task caching
│       └── model/    # Task models
├── docs/             # API documentation
├── docker-compose.yml # Docker services
├── Dockerfile        # Application Docker image
├── prometheus.yml    # Prometheus configuration
└── OBSERVABILITY.md  # Detailed observability guide
```

## Development Workflow

### Adding a New Feature

1. Create the feature in the appropriate package
2. Add tests before implementing (TDD)
3. Update the handler to record metrics
4. Test with `go test ./...`
5. Verify metrics with Prometheus

### Database Migrations

Migrations should be created in the `migrations/` directory.

### API Documentation

Update Swagger comments in handlers:
```go
// CreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task with title and description
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Router       /tasks [post]
```

Generate documentation:
```bash
swag init -g cmd/api/main.go
```

## Docker Services

The docker-compose file includes:

1. **app** - Task Manager API
2. **postgres** - PostgreSQL database
3. **redis** - Redis cache
4. **prometheus** - Prometheus for metrics collection

## Monitoring & Alerting

### Prometheus Queries

```promql
# Request rate per minute
rate(task_manager_http_requests_total[1m])

# Error rate
rate(task_manager_http_requests_total{status="5xx"}[1m])

# Average latency
rate(task_manager_http_request_latency_seconds_sum[5m]) / rate(task_manager_http_request_latency_seconds_count[5m])

# Task count by status
task_manager_tasks_count
```

### Setting Up Alerts

Add alert rules to `prometheus.yml`:

```yaml
rule_files:
  - 'alerts.yml'
```

## Troubleshooting

### Services not starting
```bash
# Check service status
docker-compose ps

# View logs
docker-compose logs -f app
```

### Database connection errors
```bash
# Check PostgreSQL health
docker-compose exec postgres pg_isready -U postgres
```

### Metrics not appearing
```bash
# Verify metrics endpoint
curl http://localhost:8080/metrics

# Check Prometheus targets
open http://localhost:9090/targets
```

## Performance Optimization

- **Caching**: Tasks are cached in Redis with 30-second TTL
- **Database**: Connection pooling is configured in PostgreSQL
- **Indexing**: Task ID is indexed in the database

## License

MIT

## Support

For issues or questions, please refer to the [OBSERVABILITY.md](./OBSERVABILITY.md) guide or check the test files for usage examples.
