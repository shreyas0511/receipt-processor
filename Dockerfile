# ---- Build Stage ----
    FROM golang:1.24-alpine AS builder

    # Install git (needed for go get if you're using external packages)
    RUN apk add --no-cache git
    
    # Set environment variables
    ENV GO111MODULE=on \
        CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64
    
    # Set working directory
    WORKDIR /app
    
    # Copy go mod and sum files
    COPY go.mod go.sum ./
    
    # Download dependencies
    RUN go mod download
    
    # Copy the source code
    COPY . .
    
    # Build the Go app
    RUN go build -o receipt-processor main.go
    
    # ---- Runtime Stage ----
    FROM alpine:latest
    
    # Set working directory
    WORKDIR /root/
    
    # Copy the binary from the builder stage
    COPY --from=builder /app/receipt-processor .
    
    # Expose port 8080
    EXPOSE 8080
    
    # Command to run the binary
    CMD ["./receipt-processor"]
    