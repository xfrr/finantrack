# display all available commands
default: help

# build the Go applications
go-build:
  go build -o bin/ ./...

# run the Go tests
go-test:
  go test ./... -v -cover

# lint the Go codebase
go-lint:
  golangci-lint run --enable-all

# format the Go codebase using gofmt
go-format:
  gofmt -s -w .

# clean up the project's build and test artifacts
clean:
  rm -rf bin/
  rm -rf vendor/
  rm -rf .env
  rm -rf .cache
  rm -rf .coverage

up:
  docker-compose -f deployments/docker/docker-compose.yml --env-file deployments/docker/.env up -d --build

# display all available commands
help:
  @echo "Usage: make <target>"
  @echo ""
  @echo "Targets:"
  @echo "  go-build    - Build the Go project"
  @echo "  go-test     - Run the Go tests"
  @echo "  go-lint     - Lint the Go codebase"
  @echo "  go-format   - Format the Go codebase"
  @echo "  clean       - Clean up the project's build and test artifacts"
  @echo "  docs        - Generate and serve documentation"
  @echo "  help        - Display this help message"


  
