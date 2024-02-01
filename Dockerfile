# Use the official Golang image as the base image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# Expose the port that the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./main"]