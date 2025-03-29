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
	go test -v -timeout 30s -coverprofile=cover.out -cover ./pkg/tests/...
	go tool cover -func=cover.out

# build: test
build:	
	go build -o $(BUILD_DIR)/$(APP_NAME) main.go

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
	
docker-compose.down:
	docker-compose down

swag:
	swag init

seed.pins:
	go run cmd/seed/seed.go