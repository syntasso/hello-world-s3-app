FROM golang:1.23 as builder

# Set the working directory
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . ./

# Build the Go application
RUN go build -o main .

# Create a minimal runtime image
FROM ubuntu:22.04

# Set the working directory
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Command to run the application
CMD ["./main"]
