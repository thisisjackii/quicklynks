# File: backend/Dockerfile
# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for a static binary, especially when using sqlite
# -ldflags="-w -s" strips debug information to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o quicklynks .

# Stage 2: Create a minimal final image
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/quicklynks .

# Copy the database schema
COPY db/schema.sql .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
# The app will create the sqlite db on first run if it doesn't exist
CMD ["./quicklynks"]