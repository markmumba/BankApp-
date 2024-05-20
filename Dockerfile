# Use the official Golang image for building the Go application
FROM golang:1.20 AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Set the working directory to where the main.go file is located
WORKDIR /app/cmd/api

# Build the Go app
RUN go build -o /app/out

# Use a minimal image as the base for the final stage
FROM gcr.io/distroless/base-debian10

# Copy the prebuilt binary from the build stage
COPY --from=build /app/out /app/out

# Expose the port that your application will run on (optional, adjust as needed)
EXPOSE 8080

# Command to run the executable
CMD ["/app/out"]
