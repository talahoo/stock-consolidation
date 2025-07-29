.PHONY: lint test coverage run

lint:
	golangci-lint run

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

run:
	go run ./cmd/stockconsolidation/main.go
