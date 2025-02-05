# Start with a Golang base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go application from cmd/main.go
RUN go build -o main ./cmd/main.go

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"]
