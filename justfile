
# display all available commands
default: help

# build the Go applications for all supported platforms
build service:
  # osx (apple silicon)
  @GOARCH=arm64 GOOS=darwin go build -o bin/{{service}}-darwin-arm64 ./cmd/{{service}}

  # osx (intel)
  @GOARCH=amd64 GOOS=darwin go build -o bin/{{service}}-darwin-amd64 ./cmd/{{service}}

  # linux
  @GOARCH=amd64 GOOS=linux go build -o bin/{{service}}-linux-amd64 ./cmd/{{service}}

  # windows
  @GOARCH=amd64 GOOS=windows go build -o bin/{{service}}-windows-amd64.exe ./cmd/{{service}}

build-web:
  @GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/web

# run all tests
test:
  @go test ./... -v -cover

# lint codebase
lint:
  @golangci-lint run --enable-all

# format the codebase
format:
  # format Go code
  @gofmt -s -w .
  # format swag docs
  @swag fmt --dir services/assets/http

# start the docker-compose stack for local development
up profile:  
  @docker-compose \
    -f deployments/docker/docker-compose.yml \
    -f deployments/docker/{{profile}}/docker-compose.yml \
    --env-file deployments/docker/{{profile}}/.{{profile}}.env \
    up -d --build --remove-orphans --force-recreate

# stop the docker-compose stack for local development
down profile:
  @docker-compose \
    -f deployments/docker/docker-compose.yml \
    -f deployments/docker/{{profile}}/docker-compose.yml \
    --env-file deployments/docker/{{profile}}/.{{profile}}.env \
    down --remove-orphans --volumes

# clean up the project's build and test artifacts
clean:
  @rm -rf bin/
  @rm -rf vendor/
  @rm -rf .env
  @rm -rf .cache
  @rm -rf .coverage

clean-docker:
  @docker system prune -f

# generate and serve documentation
docs:
  # generate assets swagger docs
  @swag init \
    --parseDependency \
    --parseDepth 2 \
    --outputTypes go,yaml \
    -d services/assets/http \
    -g http_server.go \
    -o docs/swagger-assets-http-api \

# run the vegeta load tests
vegeta service:
    @cd services/{{service}}/test && go run vegeta_targeter.go
    @echo "Plot it at 'https://hdrhistogram.github.io/HdrHistogram/plotFiles.html'"

# display all available commands
help:
  @just --list


  
