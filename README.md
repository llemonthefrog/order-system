## Order Processing System

A Go-based service for managing orders with a guaranteed deferred payment mechanism and automatic cancellation on timeout.

## Tech Stack

- **Language:** Go 1.25
- **Database:** PostgreSQL (pgx driver)
- **Architecture:** Hexagonal (Ports and Adapters)
- **Infrastructure:** Docker, Docker Compose

### Key Mechanisms

- **Scheduled Tasks:** A pattern for asynchronous order cancellation. Tasks are created within the same transaction as the order.
- **Background Worker:** A background process that checks for expired orders and transitions them to the `CANCELED` state.
- **Graceful Shutdown:** Proper shutdown of the server with clean database connection handling.

## Quick Start

1. Clone the repository.
2. Setup .env in the root directory
```
DB_USER=user
DB_PASSWORD=secret_pass
DB_NAME=db_name
DB_HOST=db
DB_PORT=some_port

DATABASE_URL=postgres://user:secret_pass@db:some_port/db_name?sslmode=disable
```
3. Make sure Docker is installed.
4. Run the project:

```bash
docker-compose up --build
```

## API Endpoints

### 1. Create Order

`POST /orders/create`

**Payload:**
```json
{ "price": 1500.50 }
```

**Response:** 201(Created)
```json
{
  "payment_url": "some_url",
  "order_id": "some_id"
}
```

### 2. Confirm Order

`POST /orders/{id}/confirm`

**Response:** 200(OK)
