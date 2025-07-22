include .env

GOOSE_DB_STRING = ${DB_USER}:${DB_PASSWORD}@tcp\(${DB_HOST}:${DB_PORT}\)/${DB_NAME}

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

run-frontend:
	npm run dev

test:
	go test ./...

test-cover:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/dmitriygnatenko/homethings-v1/internal/api/...
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;

lint:
	golangci-lint run --timeout=3m

swag:
	swag init -o "doc" -d "cmd/app,internal/api/v1,internal/dto"

migration-status:
	goose -dir migration ${DB_DRIVER} ${GOOSE_DB_STRING} status

migration-up:
	goose -dir migration ${DB_DRIVER} ${GOOSE_DB_STRING} up

migration-down:
	goose -dir migration ${DB_DRIVER} ${GOOSE_DB_STRING} down

docker-build:
	docker compose up --build --detach

docker-up:
	docker compose up --detach

docker-down:
	docker compose down

install-deps: install-lint install-goose install-minimock install-swagger

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

install-minimock:
	go install github.com/gojuno/minimock/v3/cmd/minimock@latest

install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

install-frontend:
	npm install

.PHONY: build
build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/app/app cmd/app/main.go

build-frontend:
	npm run build --emptyOutDir
