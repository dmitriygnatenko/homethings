include .env

GOOSE_DB_STRING = "user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} port=${DB_PORT} host=${DB_HOST} sslmode=disable"

usage:
	@echo "make run"
	@echo "make test"
	@echo "make test-cover"
	@echo "make swag"
	@echo "make lint"
	@echo "make migration-status"
	@echo "make migration-up"
	@echo "make migration-down"
	@echo "make docker-build"
	@echo "make docker-up"
	@echo "make docker-down"
	@echo "make install-deps"
	@echo "make build"

run:
	cd cmd/app && go run main.go

test:
	go test ./...

test-cover:
	go test ./... -coverprofile=./coverage.out
	go tool cover -html=./coverage.out

lint:
	golangci-lint run --timeout=3m

swag:
	swag init -o "docs" -d "cmd/app,internal/api/v1,internal/dto"

migration-status:
	goose -dir migrations postgres ${GOOSE_DB_STRING} status

migration-up:
	goose -dir migrations postgres ${GOOSE_DB_STRING} up

migration-down:
	goose -dir migrations postgres ${GOOSE_DB_STRING} down

docker-build:
	docker compose up --build --detach

docker-up:
	docker compose up --detach

docker-down:
	docker compose down

install-deps: install-lint install-goose install-swagger

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: build
build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/app/app cmd/app/main.go