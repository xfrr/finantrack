# Development Stage
FROM golang:1.23-alpine AS dev

ARG APP_NAME "finances-manager"

# Set app_name environment variable
ENV APP_NAME=${APP_NAME}

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Install Just cli for task automation
RUN apk add --no-cache git just bash gcc musl-dev

# Install Go tools
RUN go install github.com/go-delve/delve/cmd/dlv@latest && \
  go install github.com/air-verse/air@latest && \
  go install github.com/swaggo/swag/cmd/swag@latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first, to cache dependencies
COPY ./internal/contexts/assets/go.mod ./internal/contexts/assets/go.sum ./
RUN go mod download

# Expose debugging port
EXPOSE 40000

# Start Air for live-reloading, exposing Delve for debugging on port 40000
ENTRYPOINT ["/bin/bash", "-c", "/app/deployments/docker/entrypoint.sh"]
