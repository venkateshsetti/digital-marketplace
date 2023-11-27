# Use an official Go runtime as a parent image
FROM golang:1.21.4-alpine3.18

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download and install Go module dependencies
RUN go mod download

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
WORKDIR  /app/cmd/
RUN go build -o main .

COPY config/default_config.yml /app/cmd/config/

RUN apk add --no-cache docker-compose


# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
