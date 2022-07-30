test:
	@echo "Running tests database-app..."
	go test ./... --cover

cover:
	@echo "Running tests database-app..."
	go test ./... --coverprofile coverage.out
	go tool cover -func coverage.out

lint:
	@echo "Running golangci-lint database-app..."
	golangci-lint run

.PHONY: all clean test cover lint
