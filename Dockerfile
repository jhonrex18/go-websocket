# Use the official Golang image
FROM golang:1.20 AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go_websocket/go.mod go_websocket/go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY go_websocket/. ./

# Build the Go application
RUN go build -o main .

# Use a smaller base image to run the application
FROM gcr.io/distroless/base

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file to the runtime image
COPY go_websocket/.env ./

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["/main"]
