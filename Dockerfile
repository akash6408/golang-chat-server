# Use official Go image
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod first and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go server binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Expose app port
EXPOSE 8080

# Run the server
CMD ["./server"]
