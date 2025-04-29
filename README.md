# Order Packs Calculator

A Go application that calculates optimal pack sizes for shipping orders, with a basic web UI.
For the production I am recommending:
- Add Configuration, Environment variables: PORT, GIN_MODE, CORS origins ...
- Users' authentications and authorization
- Use TLS, CSRF middleware
- Add max pack size in validation

## Features

- Calculate optimal pack combinations
- Customizable pack sizes
- Simple web interface
- REST API

## Prerequisites

- Go 1.23+
- Docker 20.10+
- Make (optional)

## Quick Start

### Using Docker (recommended)

```bash
# Build and start containers
make docker-build
make docker-run

# Stop containers
make docker-stop

# Access the application
open http://localhost:8086/

# Run tests
make test
```

### Run from the project root:

go run cmd/server/main.go

## UI

http://127.0.0.1:8086/

## Example
GET
url: http://127.0.0.1:8086/api/v1/packs?items_no=12001

PUT
url: http://127.0.0.1:8086/api/v1/packs/sizes
{
"sizes": [23, 31, 53]
}

The default port is 8086, or it can be added in .env as PORT

### Test

## Run all tests
go test ./...

## Run service tests
go test ./internal/services/...

## Run controller tests
go test ./controllers/v1/pack/...
