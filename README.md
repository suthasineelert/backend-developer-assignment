# Backend Developer Assignment

<img src="https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go" alt="go version" />&nbsp;<a href="https://goreportcard.com/report/github.com/create-go-app/fiber-go-template" target="_blank">

[Fiber](https://gofiber.io/) is an Express.js inspired web framework build on top of Fasthttp, the fastest HTTP engine for Go. Designed to ease things up for **fast** development with **zero memory allocation** and **performance** in mind.

## ⚡️ Quick start
1. Clone the repo
2. Rename `.env.example` to `.env` and fill it with your environment values.
3. Install [Docker](https://www.docker.com/get-started) and the following useful Go tools to your system:

   - [golang-migrate/migrate](https://github.com/golang-migrate/migrate#cli-usage) for apply migrations
   - [github.com/swaggo/swag](https://github.com/swaggo/swag) for auto-generating Swagger API docs
   - [github.com/securego/gosec](https://github.com/securego/gosec) for checking Go security issues
   - [github.com/go-critic/go-critic](https://github.com/go-critic/go-critic) for checking Go the best practice issues
   - [github.com/golangci/golangci-lint](https://github.com/golangci/golangci-lint) for checking Go linter issues
   - [github.com/air-verse/air](https://github.com/air-verse/air) for live reload

4. Run project by this command:

```bash
make docker-compose.up
```

5. Go to API Docs page (Swagger): [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## API Functionality
- [x] Verify User Pin
- [x] Get User Profile
- [x] Update User Profile
- [x] Get User Greeting
- [x] Update User Greeting
- [x] Get User Transaction
- [x] Add Money to Account
- [x] Withdraw Money from Account
- [x] List Account
- [x] Get Account
- [x] Update Main Account
- [x] Update Account (Name, Color, Border Color)
- [x] List Debit Card
- [x] Update Debit Card (Name, Color, Border Color)
- [x] Create Debit Card
- [x] Delete Debit Card
- [] Get Banner

## Key Notes
- Database transaction will be handled in repository layer, but the logic will be pass through function


## TODO
- [] Unit test for all apis
- [] Optimized services latency (optional)
- [] Optimized database schema (optional)
- [] Including services stress test report such as locust, k6 (optional)

## Updated Schema
- Add column `pin` to `users` table
- Add column `created_at`, `updated_at`, `deleted_at` to all table
- Add column `amount`, `transaction_type` to `transactions` table to log transaction history, type is `debit`, `withdrawal`, or `transfer`



## Database Migration

### Apply migrations

Database will be automatically migrated on startup.

To apply migration manually, specify DATABASE_URL before running migrate command to specify database connection
```bash
make migrate.up
# or with DATABASE_URL
DATABASE_URL="mysql://user:password@tcp(localhost:3306)/assignment" make migrate.up
```

### Rollback migrations
```bash
make migrate.down
```

### Create new migration
```bash
migrate create -ext sql -dir platform/migrations -seq create_users_table
```


## Seeding Data

### Seeding Mock Data
```bash
mysql -h 127.0.0.1 -p assignment < /path/to/mock/*.sql
```
and enter mysql root password

### Mock user pin data
```
make seed.pins
```


## Testing
### Run all tests
```bash
make test
```

### Create mocks for test

Install mockery
```bash
go install github.com/vektra/mockery/v2@latest
```

Create mocks
```bash
 make mock.generate
```

## Assumptions
- User already registered and has pin



## 🗄 Project Structure

### ./app

**Folder with business logic only**. This directory contains the core business logic of the application, independent of external implementations.

- `./app/controllers` folder for functional controllers (used in routes)
- `./app/models` folder for business models representing database tables
- `./app/queries` folder for database queries related to models

### ./docs

**Folder with API Documentation**. Contains Swagger configuration files for auto-generated API documentation.

### ./pkg

**Folder with project-specific functionality**. Contains code tailored for this specific application.

- `./pkg/configs` folder for configuration functions (Fiber settings, etc.)
- `./pkg/middleware` folder for HTTP middleware components
- `./pkg/repository` folder for constants and repository interfaces
- `./pkg/routes` folder for API route definitions
- `./pkg/utils` folder with utility functions (server starter, connection URL builder, etc.)

### ./platform

**Folder with platform-level logic**. Contains infrastructure code that connects the application to external services.

- `./platform/database` folder with MySQL database setup and connection functions
- `./platform/migrations` folder with SQL migration files for database schema
- `./platform/seeds` folder for database seed files to populate test data


## ⚙️ Configuration

```ini
# .env

# Stage status to start server:
#   - "dev", for start server without graceful shutdown
#   - "prod", for start server with graceful shutdown
APP_ENV="dev"

# Server settings:
SERVER_HOST="0.0.0.0"
SERVER_PORT=8080
SERVER_READ_TIMEOUT=60

# JWT settings:
JWT_SECRET_KEY="secret"
JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT=15
JWT_REFRESH_KEY="refresh"
JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT=720

# Database settings:
DB_HOST="localhost"
DB_PORT=5432
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_NAME="assignment"
DB_SSL_MODE="disable"
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=2
```

## ⚠️ License

This project is based on the Fiber Go Template created by [Vic Shóstak](https://shostak.dev/) & [True web artisans](https://1wa.co/), licensed under Apache 2.0.

Original template: [create-go-app/fiber-go-template](https://github.com/create-go-app/fiber-go-template)
