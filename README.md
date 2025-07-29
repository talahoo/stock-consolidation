# Stock Consolidation Service

![Go CI/CD](https://github.com/talahoo/stock-consolidation/actions/workflows/go.yml/badge.svg)

A Change Data Capture (CDC) service that monitors stock changes in branch databases and consolidates them to a central HQ system.

## Features

- Real-time stock change monitoring using PostgreSQL notifications
- Automatic synchronization with HQ system
- Support for multiple stock operations (insert, update)
- REST API health check endpoint
- Docker containerization for easy deployment
- Comprehensive test coverage

## Technology Stack

- Go 1.20+
- PostgreSQL 15
- Docker & Docker Compose
- REST API
- Change Data Capture (CDC) Pattern

## Prerequisites

- Docker and Docker Compose installed
- Go 1.20 or higher (for local development)
- PostgreSQL 15 or higher (for local development)

## Installation & Setup

### Using Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/talahoo/stock-consolidation.git
   cd stock-consolidation
   ```

2. Configure environment variables (optional):
   - Copy `.env.example` to `.env` (if needed)
   - Default configurations are already set in the docker-compose.yml

3. Start the services:
   ```bash
   docker-compose up -d
   ```

This will start both the PostgreSQL database and the stock consolidation service.

### Local Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/talahoo/stock-consolidation.git
   cd stock-consolidation
   ```

2. Set up environment variables:
   ```bash
   # Windows PowerShell
   $env:DB_HOST="localhost"
   $env:DB_PORT="5432"
   $env:DB_USER="admin"
   $env:DB_PASSWORD="admin123"
   $env:DB_NAME="stockdb"
   $env:SERVICE_PORT="3000"
   $env:HQ_END_POINT="http://localhost:8085/stock"
   $env:HQ_BASIC_AUTHORIZATION="Basic dXNlcjpwYXNz"
   ```

3. Run the tests:
   ```bash
   go test ./... -cover
   ```

4. Build and run:
   ```bash
   go build -o stockconsolidation ./cmd/stockconsolidation
   ./stockconsolidation
   ```

## API Endpoints

### Health Check
- `GET /health`
  - Returns the health status of the service
  - Response: `200 OK` with body `{"status": "up"}`

## Database Structure

### Stock Table
```sql
CREATE TABLE stock (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id INTEGER NOT NULL,
    branch_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    reserved INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, branch_id)
);
```

## CDC Notification System

The service uses PostgreSQL's NOTIFY/LISTEN feature for Change Data Capture:

1. A trigger on the stock table captures changes
2. Changes are sent as notifications on the 'stock_changes' channel
3. The service listens for these notifications and forwards them to HQ

## Testing

### End-to-End Testing Flow

#### 1. Testing Stock Insert & Update Flow
```bash
# 1. Start the service and PostgreSQL
docker-compose up -d

# 2. Open PostgreSQL session to monitor notifications
docker-compose exec postgres psql -U admin -d stockdb
# In psql, execute:
LISTEN stock_changes;

# 3. In another terminal, insert test data
docker-compose exec postgres psql -U admin -d stockdb
# In psql, execute:
INSERT INTO stock (product_id, branch_id, quantity) 
VALUES (1001, 1, 100);

# You should see notification in the first terminal:
Asynchronous notification "stock_changes" received from server process with payload:
'{"product_id":1001,"branch_id":1,"quantity":100,...}'

# 4. Check service logs for HQ API call
docker-compose logs -f stock-consolidation
# You should see logs like:
# "Sending stock update to HQ for product 1001 in branch 1"
# "Successfully sent stock update to HQ"

# 5. Test update scenario
UPDATE stock 
SET quantity = 150 
WHERE product_id = 1001 AND branch_id = 1;

# Check notifications and logs again
```

#### 2. Testing Multiple Branch Updates
```sql
-- Insert stocks for different branches
INSERT INTO stock (product_id, branch_id, quantity) VALUES 
(2001, 1, 100),
(2001, 2, 150),
(2001, 3, 200);

-- Update stock in branch 1
UPDATE stock SET quantity = 120 
WHERE product_id = 2001 AND branch_id = 1;

-- Update stock in branch 2
UPDATE stock SET quantity = 180 
WHERE product_id = 2001 AND branch_id = 2;
```

#### 3. Testing Edge Cases
```sql
-- Test zero quantity
UPDATE stock SET quantity = 0 
WHERE product_id = 2001 AND branch_id = 1;

-- Test large numbers
UPDATE stock SET quantity = 999999 
WHERE product_id = 2001 AND branch_id = 2;

-- Test reserved stock
UPDATE stock SET reserved = 50 
WHERE product_id = 2001 AND branch_id = 3;
```

### Running Unit Tests
Run all tests with coverage:
```bash
go test ./... -cover
```

To run specific package tests:
```bash
# Run database adapter tests
go test ./internal/adapter/db/... -v

# Run service tests
go test ./internal/service/... -v

# Run integration tests
go test ./test/... -v
```

### Test Coverage Requirements
- Maintain minimum 80% test coverage for all packages
- Critical components (database, service layer) should have >90% coverage
- Integration tests must cover all main workflows

### Code Quality Checks

The project uses golangci-lint for code quality enforcement. Configuration is in `.golangci.yml`.

1. Install golangci-lint:
   ```bash
   # Windows
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

   # Make sure it's in your PATH
   golangci-lint --version
   ```

2. Run linter:
   ```bash
   # Run all configured linters
   golangci-lint run

   # Run specific linter
   golangci-lint run --disable-all -E errcheck

   # Run with verbose output
   golangci-lint run -v
   ```

3. Configured Linters:
   - `gofmt`: Check code formatting
   - `golint`: Check coding style
   - `govet`: Check common mistakes
   - `errcheck`: Check error handling
   - `staticcheck`: Check for various issues
   - `gosimple`: Check for code simplification
   - `ineffassign`: Check unused assignments

4. Fix Common Issues:
   ```bash
   # Auto-fix formatting issues
   gofmt -w .

   # Run linter with fix flag
   golangci-lint run --fix
   ```

## Troubleshooting

### Common Issues

1. PostgreSQL Connection Issues
   - Verify PostgreSQL is running: `docker ps | grep postgres`
   - Check connection settings in environment variables
   - Ensure database user has proper permissions
   - Check PostgreSQL logs: `docker logs <postgres-container-id>`

2. Empty Notifications
   - Verify trigger function is properly installed
   - Check PostgreSQL notification channel: `stock_changes`
   - Monitor PostgreSQL logs for notification events
   - Use `psql` to test NOTIFY/LISTEN manually

3. HQ Sync Failures
   - Verify HQ endpoint is accessible
   - Check Basic Authorization credentials
   - Inspect network connectivity between services
   - Review service logs for detailed error messages

4. Docker Issues
   - Ensure all required ports are available
   - Check Docker logs: `docker-compose logs -f`
   - Verify Docker network connectivity
   - Try rebuilding containers: `docker-compose up -d --build`

## Debug Tips

### Local Development

1. Enable Debug Logging
   ```bash
   # Set environment variable for verbose logging
   $env:LOG_LEVEL="debug"
   ```

2. Monitor PostgreSQL Notifications
   ```sql
   -- In psql
   LISTEN stock_changes;
   -- Make changes to stock table and watch notifications
   ```

3. Test Database Triggers
   ```sql
   -- Insert test data
   INSERT INTO stock (product_id, branch_id, quantity) VALUES (1, 1, 10);
   -- Update to trigger notification
   UPDATE stock SET quantity = 20 WHERE product_id = 1 AND branch_id = 1;
   ```

4. Use VS Code Debugging
   - Launch configuration is provided in `.vscode/launch.json`
   - Set breakpoints in key functions
   - Use debug console to inspect variables
   - Monitor goroutines and stack traces

### Docker Environment

1. Access Container Logs
   ```bash
   # View service logs
   docker-compose logs -f stock-consolidation
   # View database logs
   docker-compose logs -f postgres
   ```

2. Connect to Running Containers
   ```bash
   # Access PostgreSQL container
   docker-compose exec postgres psql -U admin -d stockdb
   # Access service container
   docker-compose exec stock-consolidation sh
   ```

3. Network Debugging
   ```bash
   # Test HQ endpoint
   curl -v -H "Authorization: $env:HQ_BASIC_AUTHORIZATION" $env:HQ_END_POINT
   # Check service health
   curl http://localhost:3000/health
   ```

## License

This project is licensed under the MIT License.
