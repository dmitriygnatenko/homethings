
usage:
	@echo "make run"
	@echo "make lint"
	@echo "make install-lint"
run:
	cd cmd/app && go run main.go

lint:
	golangci-lint run --timeout=3m

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@master
