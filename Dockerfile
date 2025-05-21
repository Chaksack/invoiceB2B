# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for private repositories if needed
# RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application for the cmd/api service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api_server ./cmd/api/main.go

# Stage 2: Create the final small image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/api_server .

# Expose the port the app runs on
EXPOSE 3000

# Command to run the application
CMD ["./api_server"]