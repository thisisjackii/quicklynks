# File: /quicklynks/backend/Dockerfile
# --- Build Stage ---
FROM golang:1.21-alpine3.16 AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for a static binary
# -ldflags="-w -s" strips debug information to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /quicklynks-api ./cmd/api

# --- Final Stage ---
FROM alpine:3.19

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /quicklynks-api /quicklynks-api

# Update Alpine packages to latest security patches
RUN apk update && apk upgrade

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/quicklynks-api"]