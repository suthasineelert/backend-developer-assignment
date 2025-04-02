.PHONY: clean critic security lint test build run

APP_NAME = backend-assignment
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/platform/migrations
DATABASE_URL = mysql://mysql:mysql@tcp(127.0.0.1:3306)/assignment

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go test -v -coverprofile=cover.out -coverpkg=./app/controllers,./app/services ./pkg/tests/...
	go tool cover -html=cover.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

build: test
	go build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME)

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "${DATABASE_URL}" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "${DATABASE_URL}" down

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "${DATABASE_URL}" force $(version)

docker-compose.up:
	docker-compose up --build -d

docker-compose.up.backend:
	docker-compose up --build -d backend
	
docker-compose.down:
	docker-compose down

swag:
	swag init -g cmd/main.go

seed.pins:
	go run cmd/seed/seed.go

mock.generate:
	mockery --dir=app/services --output=pkg/mocks/services --outpkg=mocks --case=snake --all && mockery --dir=app/repositories --output=pkg/mocks/repositories --outpkg=mocks --case=snake --all

# Run stress test with output to HTML report
stress.report:
	K6_WEB_DASHBOARD=true K6_WEB_DASHBOARD_EXPORT=html-report.html k6 run platform/stress-tests/transaction-test.js