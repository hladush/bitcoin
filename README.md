# Bitcoin Address Tracker

A REST API service for tracking Bitcoin addresses, synchronizing their transactions, and monitoring balances. Built with Go and SQLite, using the Blockchair API for blockchain data.

## Features

- **Address Management**: Add/remove Bitcoin addresses for tracking
- **Transaction Synchronization**: Fetch and store transaction history from blockchain APIs
- **Balance Tracking**: Monitor confirmed and unconfirmed balances
- **Background Sync**: Automatic periodic synchronization of all tracked addresses
- **REST API**: Clean HTTP API for all operations
- **Data Persistence**: SQLite database for reliable data storage

## Architecture

### High-Level Components

1. **REST API Server** (`cmd/server/main.go`)
   - HTTP handlers for all endpoints
   - Middleware for CORS and logging
   - Graceful shutdown handling

2. **Service Layer** (`internal/services/`)
   - Business logic for address management
   - Transaction synchronization logic
   - Validation and error handling

3. **Repository Layer** (`internal/repository/`)
   - SQLite database operations
   - Data persistence with proper indexing
   - Transaction management

4. **External Client** (`internal/clients/`)
   - Blockchair API integration
   - Bitcoin address validation
   - Transaction data fetching

5. **Data Models** (`internal/models/`)
   - Core data structures
   - API request/response models
   - Standardized error handling

### Architecture Decisions

- **SQLite Database**: Chosen for simplicity and portability. Easy to deploy without external dependencies.
- **Blockchair API**: Selected for reliable blockchain data and good documentation.
- **Repository Pattern**: Separates data access logic for better testability and maintainability.
- **Service Layer**: Encapsulates business logic and coordinates between repository and external APIs.
- **Background Sync**: Automatic synchronization every 5 minutes ensures data freshness.

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Address Management
- `GET /addresses` - List all tracked addresses with balances
- `POST /addresses` - Add a new address to track
- `GET /addresses/{address}` - Get specific address details
- `DELETE /addresses/{address}` - Remove address from tracking

### Balance and Transactions
- `GET /addresses/{address}/balance` - Get current balance
- `GET /addresses/{address}/transactions` - Get transaction history (with pagination)

### Synchronization
- `POST /addresses/{address}/sync` - Manually sync specific address
- `POST /sync` - Sync all tracked addresses

## Setup and Installation

### Prerequisites
- Go 1.19 or higher
- CGO enabled (for SQLite)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd bitcoin
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build -o bitcoin-tracker cmd/server/main.go
```

4. Run the server:
```bash
./bitcoin-tracker
```

The server will start on port 8080 and create a SQLite database file `bitcoin_tracker.db` in the current directory.

### Development Mode

To run in development mode with automatic recompilation:
```bash
go run cmd/server/main.go
```

## Usage Examples

### Add an Address
```bash
curl -X POST http://localhost:8080/addresses \
  -H "Content-Type: application/json" \
  -d '{"address": "bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5", "label": "My Wallet"}'
```

### Get All Addresses
```bash
curl http://localhost:8080/addresses
```

### Get Address Balance
```bash
curl http://localhost:8080/addresses/bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5/balance
```

### Get Address Transactions
```bash
curl "http://localhost:8080/addresses/bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5/transactions?limit=10&offset=0"
```

### Sync Address Manually
```bash
curl -X POST http://localhost:8080/addresses/bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5/sync
```

## Sample Addresses for Testing

Use these Bitcoin addresses for testing (from the assignment):

1. `3E8ociqZa9mZUSwGdSmAEMAoAxBK3FNDcd`
2. `bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5`
3. `bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h` (high activity - 156k+ transactions)
4. `12xQ9k5ousS8MqNsMBqHKtjAtCuKezm2Ju` (high activity - 900+ transactions)

## Data Models

### Address
```json
{
  "id": 1,
  "address": "bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5",
  "label": "My Wallet",
  "created_at": "2024-01-01T00:00:00Z",
  "last_synced": "2024-01-01T00:05:00Z"
}
```

### Transaction
```json
{
  "id": 1,
  "hash": "abcd1234...",
  "address": "bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5",
  "amount": 100000000,
  "confirmations": 6,
  "block_height": 800000,
  "timestamp": "2024-01-01T00:00:00Z",
  "type": "received"
}
```

### Balance
```json
{
  "address": "bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5",
  "confirmed_balance": 100000000,
  "unconfirmed_balance": 0,
  "total_balance": 100000000,
  "balance_btc": 1.0
}
```

## Configuration

### Environment Variables
- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database file path (default: bitcoin_tracker.db)
- `SYNC_INTERVAL`: Background sync interval (default: 5m)

### Database Schema

The application creates two main tables:

**addresses**
- `id`: Primary key
- `address`: Unique Bitcoin address
- `label`: Optional user-defined label
- `created_at`: Creation timestamp
- `last_synced`: Last synchronization timestamp

**transactions**
- `id`: Primary key
- `hash`: Transaction hash
- `address`: Associated Bitcoin address
- `amount`: Amount in satoshis
- `confirmations`: Number of confirmations
- `block_height`: Block height
- `timestamp`: Transaction timestamp
- `type`: Transaction type (sent/received)

## Assumptions Made

1. **Transaction Types**: Simplified to "sent" and "received" based on balance change direction
2. **Confirmations**: Uses a simplified confirmation model (6 confirmations for confirmed transactions)
3. **Rate Limiting**: Relies on Blockchair API rate limits (no additional client-side limiting)
4. **Error Handling**: Graceful degradation - sync failures don't block other operations
5. **Pagination**: Default limit of 50 transactions, maximum of 100 per request
6. **Address Validation**: Basic format validation (length and prefix checking)
7. **Concurrent Access**: SQLite handles concurrent reads; writes are synchronized
8. **Background Sync**: Runs every 5 minutes; configurable via environment variable

## Testing

### Manual Testing
1. Start the server: `go run cmd/server/main.go`
2. Add a test address using curl or a REST client
3. Wait for automatic sync or trigger manual sync
4. Check balance and transactions via API

### Unit Tests
```bash
go test ./...
```

## Deployment

### Docker (Future Enhancement)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bitcoin-tracker cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bitcoin-tracker .
EXPOSE 8080
CMD ["./bitcoin-tracker"]
```

### Production Considerations
1. Use environment variables for configuration
2. Implement proper logging (structured logging)
3. Add rate limiting and authentication
4. Consider PostgreSQL for production database
5. Add comprehensive monitoring and metrics
6. Implement graceful shutdown handling
7. Add database migrations system

## Future Enhancements

1. **WebSocket Support**: Real-time balance and transaction updates
2. **Multiple Cryptocurrencies**: Extend beyond Bitcoin
3. **Advanced Analytics**: Portfolio tracking, profit/loss calculations
4. **Notification System**: Email/SMS alerts for transactions
5. **Web Dashboard**: Simple UI for address management
6. **API Authentication**: JWT or API key-based auth
7. **Caching Layer**: Redis for improved performance
8. **Metrics and Monitoring**: Prometheus integration

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is created as part of a technical assignment for CoinTracker.

---

**Note**: This is a prototype system designed for demonstration purposes. For production use, additional security, monitoring, and scalability considerations would be necessary.
