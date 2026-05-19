# To-Do REST API (Go)

A simple to-do application written in Go that exposes a REST API for creating,
reading, listing, updating, and deleting to-do items. Items are stored
in-memory; the list endpoint sorts by due date and excludes completed items by
default.

## Run it locally

### Prerequisites

- Docker Desktop running
- Go 1.25+ installed

### 1. Start PostgreSQL

In PowerShell:

```powershell
cd D:\Study\NorthStar\AssignmentTest
docker compose up -d
docker compose ps          # wait until status shows (healthy)
```

### 2. Start the API server

In a second PowerShell window:

```powershell
cd D:\Study\NorthStar\AssignmentTest
$env:DATABASE_URL = "postgres://todo:todo@localhost:5432/todo?sslmode=disable"
go run ./cmd/todoserver
```

You should see `listening on :8080`. The server pings the database and applies
`internal/todo/schema.sql` on startup.

