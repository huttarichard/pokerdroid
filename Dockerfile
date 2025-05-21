FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc g++ musl-dev 

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with explicit output binary
RUN go build -v -o /app/server ./cmd/main.go

# Use a smaller base image for the final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies for C++ libraries
RUN apk add --no-cache libstdc++ libgcc

# Copy only the binary from builder stage
COPY --from=builder /app/server .

# Expose the port the app runs on
EXPOSE 8080

ENTRYPOINT ["/app/server"]
