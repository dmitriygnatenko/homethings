include config/.env

usage:
	@echo "make run"
	@echo "make test"
	@echo "make test-cover"
	@echo "---------- Linter ----------"
	@echo "make lint"
	@echo "make install-lint"
	@echo "---------- Docker ----------"
	@echo "make docker-build"
	@echo "make docker-up"
	@echo "make docker-down"

run:
	cd cmd/app && go run main.go

test:
	go test ./...

test-cover:
	go test ./... -coverprofile=build/coverage.out
	go tool cover -html=build/coverage.out

lint:
	golangci-lint run --timeout=3m

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@master

docker-build:
	docker compose --file build/docker/docker-compose.yml --env-file config/.env --project-name homethings up --build --detach

docker-up:
	docker compose --file build/docker/docker-compose.yml --env-file config/.env --project-name homethings up --detach

docker-down:
	docker compose --file build/docker/docker-compose.yml --env-file config/.env --project-name homethings down