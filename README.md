# Build a Real-Time Crypto Tracker with TimescaleDB and Coin Gecko

[Build a Real-Time Crypto Tracker with TimescaleDB](https://www.youtube.com/watch?v=GlEaGCMbCpw)

## Prerequisites

- Go 1.18+
- Node.js 18+ and npm/yarn
- Docker and Docker Compose

## Getting Started

### 1. Start the Database

```sh
docker compose up -d
```

This starts TimescaleDB which is required for both the backend and frontend.

To stop the database:

```sh
docker compose down -v
```

### 2. Environment Variables

Both the backend and frontend use `.env` files. Configure them before running:

**Backend** (`.env` in root):
```sh
TIMESCALE_SERVICE_URL=postgres://postgres:password@127.0.0.1:5432/timescale_db
PGPASSWORD=password
PGUSER=postgres
PGDATABASE=timescale_db
PGHOST=localhost
PGPORT=5432
PGSSLMODE=disable
COIN_GECKO_API_KEY=<your-coingecko-api-key>
```

**Frontend** (`.env` in `tiger-tracker-frontend/`):
```sh
TIMESCALE_SERVICE_URL=postgres://postgres:password@127.0.0.1:5432/timescale_db?sslmode=disable
PGPASSWORD=password
PGUSER=postgres
PGDATABASE=timescale_db
PGHOST=localhost
PGPORT=5432
PGSSLMODE=disable
COIN_GECKO_API_KEY=<your-coingecko-api-key>
```

### 3. Run the Go Backend

The Go backend ingests crypto prices from the CoinGecko API and stores them in TimescaleDB every 30 seconds.

```sh
go run main.go
```

You should see output like:
```
Calling Coin Gecko API...
API Response as struct {...}
```

The backend will continue running in the foreground, periodically fetching and storing price data.

### 4. Run the Next.js Frontend

In a new terminal, navigate to the frontend directory and start the development server:

```sh
cd tiger-tracker-frontend
npm install  # if you haven't already
npm run dev
```

The frontend will be available at `http://localhost:3000`

## Architecture

- **Backend (Go)**: Runs a continuous price ingestion process that fetches BTC, ETH, and SOL prices from CoinGecko and stores them in TimescaleDB
- **Frontend (Next.js)**: Web UI that displays crypto price data stored in the database
- **Database (TimescaleDB)**: Time-series database for storing historical price data
