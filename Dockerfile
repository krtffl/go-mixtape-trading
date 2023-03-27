# Use the official Golang image as the base image
FROM golang:1.17.1-alpine3.14

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN go build -o app ./cmd/server

# Expose port 8080 for incoming traffic
EXPOSE 8080

# Start the server when the container starts
CMD ["./app"]
