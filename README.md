#  Build a Real-Time Crypto Tracker with TimescaleDB

[Build a Real-Time Crypto Tracker with TimescaleDB](https://www.youtube.com/watch?v=GlEaGCMbCpw)

## Docker

Database Up
```sh
docker compose up -d
```
Database Down
```sh
docker compose down
```

## Environment Variables

```sh
# This file contains environmental variables for use with your Tiger Cloud instance
TIMESCALE_SERVICE_URL=
PGPASSWORD=
PGUSER=
PGDATABASE=
PGHOST=
PGPORT=
PGSSLMODE=

COIN_GECKO_API_KEY=
```

## Run

```go
go run main.go
```
