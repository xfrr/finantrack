# Development Stage
FROM golang:1.23-alpine AS dev

ARG APP_NAME

# Set app_name environment variable
ENV APP_NAME=${APP_NAME}

# Install dependencies for Go and Air
RUN apk add --no-cache git bash curl \
  && curl -fLo /bin/air https://github.com/air-verse/air/releases/download/v1.61.0/air_1.61.0_linux_amd64 \
  && chmod +x /bin/air

# Install Delve for debugging Go applications
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first, to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Expose the application port (default Go port 8080)
EXPOSE 3000

# Start Air for live-reloading, exposing Delve for debugging on port 40000
ENTRYPOINT ["/bin/bash", "-c", "/app/deployments/docker/entrypoint.sh"]
