
# Start with the official Golang image to build the application
FROM golang:1.18-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the entire project directory to the working directory inside the container
COPY . .

COPY .env .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux  go build -o main ./cmd/api/main.go

# Start a new stage from a minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates
# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port that the application will run on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
