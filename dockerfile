# Use the official Golang image
FROM golang:1.23.5-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy Go mod/sum files first (for caching dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go binary (note: entrypoint moved!)
RUN go build -o /app/bin/server ./cmd/server

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/bin/server .

# Copy .env if you want it baked in (optional)
# COPY .env . 

# Expose API port
EXPOSE 8080

# Run the binary
CMD ["./server"]
