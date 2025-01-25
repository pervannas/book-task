# Use the official Golang image to build the Go application
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o myapp .

# Use a smaller image to run the application
# FROM alpine:latest
# # Install build dependencies
# RUN apk add --no-cache build-base
# # Set the Current Working Directory inside the container
# WORKDIR /app
#
# # Copy the binary from the builder stage
# COPY --from=builder /app/myapp .
#
# Expose the port the app runs on
EXPOSE 1323

# Command to run the executable
CMD ["./myapp"]
# ENTRYPOINT [ "./myapp" ]
