# Stock Consolidation Service

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

Run all tests with coverage:
```bash
go test ./... -cover
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License.
