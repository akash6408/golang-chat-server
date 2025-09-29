# Use official Go image
FROM golang:1.25

WORKDIR /app

# Copy go.mod first and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go server binary
RUN go build -o server ./cmd/server

# Expose app port
EXPOSE 8080

# Run the server
CMD ["./cmd/server/main.go"]
