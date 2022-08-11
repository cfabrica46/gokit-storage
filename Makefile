BUILDPATH=$(CURDIR)
API_NAME=storage

build: 
	@echo "Creating Binary ..."
	@go build -ldflags '-s -w' -o $(BUILDPATH)/build/bin/${API_NAME} cmd/main.go
	@echo "Binary generated in build/bin/${API_NAME}"
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
