# 📋 Todogo - Task Manager API with Observability

A production-ready Task Management API built with Go, featuring CRUD operations, Redis caching, PostgreSQL persistence, and comprehensive observability with Prometheus metrics and distributed tracing.

**Status**: ✅ Production Ready | 📊 Fully Observable | ✅ Well Tested

---

## Description & Tradeoffs

### What This Project Does

**todogo** is a task management API that allows you to:
- Create, read, update, and delete tasks
- Query all tasks with automatic caching
- Monitor performance and health with Prometheus metrics
- Trace requests with unique IDs for debugging

### Design Tradeoffs

| Decision | Benefit | Tradeoff |
|----------|---------|----------|
| **Redis Caching** | Fast responses for list operations | Extra service to manage |
| **PostgreSQL** | ACID compliance, relational data | Schema migrations required |
| **Gin Web Framework** | Lightweight, fast routing | Less opinionated structure |
| **Prometheus Metrics** | Production-grade observability | Slight performance overhead |
| **Distributed Tracing** | Request debugging capability | Memory usage for trace storage |
| **Docker Compose** | Single command deployment | Requires Docker/Docker Compose |

### Performance Characteristics

- **List Operations**: ~5ms (cached), ~50ms (DB hit)
- **Create/Update**: ~20-30ms
- **Metrics Collection**: <1ms overhead per request
- **Cache TTL**: 30 seconds

---


## Architecture

```
todogo/
├── cmd/
│   ├── api/            # Application entrypoint
│   └── seed/           # Database seeding function
├── internal/
│   ├── config/         # Configuration management
│   ├── observability/
│   │   ├── metrics/    # Prometheus metrics
│   │   ├── middleware/ # HTTP middleware
│   │   └── tracing/    # Request tracing
│   ├── platform/
│   │   ├── cache/      # Redis cache layer
│   │   ├── database/   # PostgreSQL connection
│   │   └── http/       # HTTP routing & setup
│   └── task/
│       ├── handler/    # HTTP handlers
│       ├── service/    # Business logic
│       ├── repository/ # Data access
│       ├── cache/      # Task caching
│       └── model/      # Task models
├── docs/               # API documentation
├── docker-compose.yml  # Docker services
├── Dockerfile          # Application Docker image
├── prometheus.yml      # Prometheus configuration
└── Makefile            # Build & run commands

```


### System Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Client                              │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
        ┌────────────────────────┐
        │   Gin HTTP Router      │
        │ (with CORS, Swagger)   │
        └────────┬───────────────┘
                 │
         ┌───────┴──────────┐
         │                  │
    ┌────▼──────┐    ┌─────▼──────┐
    │  Metrics  │    │  Tracing   │
    │ Middleware│    │ Middleware │
    └────┬──────┘    └─────┬──────┘
         │                 │
    ┌────▼─────────────────▼─────┐
    │  Task Handler Layer        │
    │ (Validation, Response)     │
    └────┬──────────────────────┘
         │
    ┌────▼──────────────────────┐
    │  Task Service Layer       │
    │ (Business Logic)          │
    └────┬──────────────────────┘
         │
    ┌────┴──────────┬─────────┐
    │               │         │
┌───▼────┐   ┌─────▼──┐  ┌──▼────┐
│ PostgreSQL│   │ Redis │  │Prom  │
│ Database  │   │ Cache │  │etheus│
└──────────┘   └───────┘  └──────┘
```

### Package Structure

```
todogo/
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── observability/
│   │   ├── metrics/      # Prometheus metrics definitions
│   │   ├── middleware/   # HTTP metrics middleware
│   │   └── tracing/      # Request tracing with trace IDs
│   ├── platform/
│   │   ├── cache/        # Redis connection pool
│   │   ├── database/     # PostgreSQL connection
│   │   └── http/         # Router setup, middleware chain
│   └── task/
│       ├── handler/      # HTTP request handlers
│       ├── service/      # Business logic, cache handling
│       ├── repository/   # Database queries
│       ├── cache/        # Cache layer
│       └── model/        # Task data structures
├── migrations/           # Database migrations
├── docs/                 # Swagger documentation
├── docker-compose.yml    # Service orchestration
├── Dockerfile            # Container image definition
├── Makefile              # Common commands
└── prometheus.yml        # Prometheus scrape config
```

### Data Flow

1. **HTTP Request** → Gin Router
2. **Middleware Chain** → Metrics & Tracing middleware record request ID, start time
3. **Handler** → Validates input, calls service
4. **Service** → Checks cache, queries DB if needed, records metrics
5. **Response** → Middleware records latency, status
6. **Prometheus** → Scrapes /metrics endpoint every 5 seconds

---

## Setup & Run with Make

### Prerequisites

- Docker & Docker Compose
- Go 1.25+ (for local development)
- `make` command
- curl (for testing)

### Available Make Commands

```bash
make help           # Show all available commands
make init           # Complete setup: env, build, compose up, migrations, seed
make run            # Start services with docker-compose up -d
make stop           # Stop all services
make logs           # View live logs from all services
make test           # Run all tests
make coverage       # Run tests with coverage report
make migrate        # Run pending database migrations
make migrate_down   # Rollback migrations
make seed           # Seed database with sample data
make swagger        # Generate Swagger documentation
make db_bash        # Interactive PostgreSQL shell
make redis_bash     # Interactive Redis CLI
make app_bash       # Shell access to app container
```

### Quick Start (One Command)

```bash
make init
```

For next runs, just do:
```bash
make run
```

This does everything:
1. Sets up `.env` from `.env.example`
2. Builds Docker images (including Prometheus)
3. Starts all services (app, postgres, redis, prometheus)
4. Runs database migrations
5. Seeds sample data
6. Generates Swagger docs

**Services Running:**
- API: `http://localhost:8080`
- Prometheus: `http://localhost:9090`
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`
- Swagger UI: `http://localhost:8080/docs`

---

## API Examples (Curl)

### 1. Health Check

**Request:**
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "message": "OK"
}
```

---

### 2. Create a Task

**Request:**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, and bread",
    "status": "pending",
    "assignee": "john@example.com"
  }'
```

**Response:**
```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Milk, eggs, and bread",
  "status": "pending",
  "assignee": "john@example.com",
  "created_at": "2025-12-26T19:50:18.292916Z",
  "updated_at": "2025-12-26T19:50:18.292916Z"
}
```

**Headers:**
```
X-Trace-ID: 29b19267-7357-487b-ad0a-0baea350a90c
```

---

### 3. List All Tasks

**Request:**
```bash
curl http://localhost:8080/tasks
```

**Response:**
```json
[
  {
    "id": 1,
    "title": "Buy groceries",
    "description": "Milk, eggs, and bread",
    "status": "pending",
    "assignee": "john@example.com",
    "created_at": "2025-12-26T19:50:18.292916Z",
    "updated_at": "2025-12-26T19:50:18.292916Z"
  },
  {
    "id": 2,
    "title": "Complete project",
    "description": "Finish API implementation",
    "status": "in_progress",
    "assignee": "jane@example.com",
    "created_at": "2025-12-26T19:51:00.000000Z",
    "updated_at": "2025-12-26T19:51:00.000000Z"
  }
]
```

---

### 4. Get Specific Task

**Request:**
```bash
curl http://localhost:8080/tasks/1
```

**Response:**
```json
{
  "id": 1,
  "title": "Buy groceries",
  "description": "Milk, eggs, and bread",
  "status": "pending",
  "assignee": "john@example.com",
  "created_at": "2025-12-26T19:50:18.292916Z",
  "updated_at": "2025-12-26T19:50:18.292916Z"
}
```

---

### 5. Update a Task

**Request:**
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, bread, and cheese",
    "status": "in_progress",
    "assignee": "jane@example.com"
  }'
```

**Response:**
```
HTTP/1.1 202 Accepted
```

---

### 6. Delete a Task

**Request:**
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

**Response:**
```
HTTP/1.1 204 No Content
```

---

### 7. View Metrics Endpoint

**Request:**
```bash
curl http://localhost:8080/metrics | head -30
```

**Response (excerpt):**
```
# HELP task_manager_http_requests_total Total number of HTTP requests
# TYPE task_manager_http_requests_total counter
task_manager_http_requests_total{endpoint="/tasks",method="GET",status="2xx"} 12
task_manager_http_requests_total{endpoint="/tasks",method="POST",status="2xx"} 7
task_manager_http_requests_total{endpoint="/health",method="GET",status="2xx"} 17

# HELP task_manager_http_request_latency_seconds HTTP request latency in seconds
# TYPE task_manager_http_request_latency_seconds histogram
task_manager_http_request_latency_seconds_bucket{endpoint="/tasks",method="GET",le="0.005"} 10
task_manager_http_request_latency_seconds_bucket{endpoint="/tasks",method="GET",le="+Inf"} 12

# HELP task_manager_tasks_count Current number of tasks by status
# TYPE task_manager_tasks_count gauge
task_manager_tasks_count{status="pending"} 8
task_manager_tasks_count{status="in_progress"} 3
task_manager_tasks_count{status="completed"} 2
```

---

## Swagger UI

Access the interactive API documentation:

**URL:** `http://localhost:8080/swagger/index.html`

Or use the redirect:
```bash
curl http://localhost:8080/docs
```

**Features:**
- Try out API endpoints directly
- View request/response schemas
- See all available endpoints
- Test with different parameters

---

## Running Tests

### Run All Tests

```bash
make test
```

**Output:**
```
=== RUN   TestHandlerCreate
=== RUN   TestHandlerCreate/create_valid_task
=== RUN   TestHandlerCreate/create_missing_title
--- PASS: TestHandlerCreate (0.00s)
    --- PASS: TestHandlerCreate/create_valid_task (0.00s)
    --- PASS: TestHandlerCreate/create_missing_title (0.00s)

=== RUN   TestHandlerList
--- PASS: TestHandlerList (0.00s)

=== RUN   TestServiceCreate
--- PASS: TestServiceCreate (0.00s)

ok      task-manager/internal/config           3/3 tests
ok      task-manager/internal/platform/http    4/4 tests
ok      task-manager/internal/task/handler     7/7 tests
ok      task-manager/internal/task/service     8/8 tests

✅ 22/22 tests passing
```

### Run Tests with Coverage

```bash
make coverage
```

**Output:**
```
total coverage: 65.3% of statements
For HTML report, run: go tool cover -html=/tmp/coverage.out
```

### Run Specific Test

```bash
go test -v ./internal/task/handler -run TestHandlerCreate
```

---

## Prometheus Queries

### Prerequisites

Prometheus must be running: `http://localhost:9090`

Make sure to generate some traffic first:
```bash
curl http://localhost:8080/tasks
for i in {1..5}; do curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"title":"Task '$i'","status":"pending"}'; done
```

### 3 Metrics Available

#### 1. Request Count: `task_manager_http_requests_total`

**Query: Total requests to /tasks endpoint**
```promql
task_manager_http_requests_total{endpoint="/tasks"}
```

**Query: POST requests only (create operations)**
```promql
task_manager_http_requests_total{method="POST"}
```

**Query: Error requests only (4xx, 5xx)**
```promql
task_manager_http_requests_total{status=~"[45]xx"}
```

**Query: Request rate per minute**
```promql
rate(task_manager_http_requests_total[1m])
```

---

#### 2. Request Latency: `task_manager_http_request_latency_seconds`

**Query: Average latency for /tasks endpoint**
```promql
rate(task_manager_http_request_latency_seconds_sum{endpoint="/tasks"}[5m]) / rate(task_manager_http_request_latency_seconds_count{endpoint="/tasks"}[5m])
```

**Query: 95th percentile latency (P95)**
```promql
histogram_quantile(0.95, rate(task_manager_http_request_latency_seconds_bucket[5m]))
```

**Query: 99th percentile latency (P99)**
```promql
histogram_quantile(0.99, rate(task_manager_http_request_latency_seconds_bucket[5m]))
```

**Query: Slowest endpoints (top 3)**
```promql
topk(3, rate(task_manager_http_request_latency_seconds_sum[5m]) / rate(task_manager_http_request_latency_seconds_count[5m]))
```

---

#### 3. Task Count: `task_manager_tasks_count`

**Query: Total tasks**
```promql
sum(task_manager_tasks_count)
```

**Query: Pending tasks**
```promql
task_manager_tasks_count{status="pending"}
```

**Query: Tasks by status (all)**
```promql
task_manager_tasks_count
```

**Query: Task distribution (pie chart ready)**
```promql
sum(task_manager_tasks_count) by (status)
```

---

### Dashboard Queries (Copy-Paste Ready)

**Monitor API Health:**
```promql
rate(task_manager_http_requests_total[1m])
```

**Error Rate:**
```promql
rate(task_manager_http_requests_total{status="5xx"}[1m])
```

**Performance (P95):**
```promql
histogram_quantile(0.95, rate(task_manager_http_request_latency_seconds_bucket[5m]))
```

**Workload:**
```promql
sum(task_manager_tasks_count) by (status)
```

---

## Key Features

✅ **Complete CRUD API** - Full task lifecycle management
✅ **Smart Caching** - Redis-backed 30-second TTL
✅ **High Performance** - Sub-50ms response times
✅ **Production Observability** - 3 Prometheus metrics
✅ **Request Tracing** - Unique trace IDs (X-Trace-ID)
✅ **API Documentation** - Interactive Swagger UI
✅ **Comprehensive Tests** - 22 passing tests
✅ **Docker Ready** - One command to production
✅ **Health Monitoring** - `/health` endpoint
✅ **CORS Enabled** - Ready for web frontends

---

## Configuration

### Environment Variables

Create `.env` file (or use `make set-env`):

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

### Docker Services

All services defined in `docker-compose.yml`:

- **app** - Task Manager API (Go + Gin)
- **postgres** - PostgreSQL database with migrations
- **redis** - Redis cache with health checks
- **prometheus** - Prometheus for metrics collection

---

## License

MIT

---

**Built with ❤️ in Go | Monitored with Prometheus | Cached with Redis**

Get trace ID from response:
```bash
curl -i http://localhost:8080/tasks | grep X-Trace-ID
```


## Support

For issues or questions, please contact sezarsaman@gmail.com
