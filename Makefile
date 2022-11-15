include config/.env

GOOSE_DB_STRING = ${DB_USER}:${DB_PASSWORD}@/${DB_NAME}?parseTime=true

usage:
	@echo "make run"
	@echo "make test"
	@echo "make test-cover"
	@echo "make lint"
	@echo "make migration-status"
	@echo "make migration-up"
	@echo "make migration-down"
	@echo "make docker-build"
	@echo "make docker-up"
	@echo "make docker-down"
	@echo "make install-deps"

run: run-swag
	cd cmd/app && go run main.go

test:
	go test ./...

test-cover:
	go test ./... -coverprofile=build/coverage.out
	go tool cover -html=build/coverage.out

lint:
	golangci-lint run --timeout=3m

run-swag:
	swag init --pd -o "./internal/docs"  -d "cmd/app,internal/api/v1"

migration-status:
	goose -dir migrations mysql ${GOOSE_DB_STRING} status

migration-up:
	goose -dir migrations mysql ${GOOSE_DB_STRING} up

migration-down:
	goose -dir migrations mysql ${GOOSE_DB_STRING} down

docker-build:
	docker compose --file build/docker/docker-compose.yml --env-file config/.env --project-name homethings up --build --detach

docker-up:
	docker compose --file build/docker/docker-compose.yml --env-file config/.env --project-name homethings up --detach

docker-down:
	docker compose --file build/docker/docker-compose.yml --env-file config/.env --project-name homethings down

install-deps: install-lint install-goose install-swagger

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest