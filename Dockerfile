# Start from the official Golang image to build the binary.
FROM golang:1.18 as builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum to download dependencies.
# Doing this before copying the entire source code helps to cache the dependencies layer.
COPY go.mod go.sum ./

# Download dependencies.
RUN go mod download

# Copy the source code from the current directory to the /app directory in the container.
COPY . .

# Build the application. Adjust the build command if needed.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch for a smaller, final image.
FROM alpine:latest  

# Set the working directory.
WORKDIR /root/

# Copy the binary from the builder stage.
COPY --from=builder /app/main .

# Expose the port your application listens on.
EXPOSE 3000

# Command to run the executable.
CMD ["./main"]
