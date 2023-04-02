# First stage: Build Golang application
FROM golang:1.19 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source files
COPY . .

# Build the Golang application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Second stage: Setup the PostgreSQL database
FROM postgres:13

# Initialize the database with schema and data
COPY ./init.sql /docker-entrypoint-initdb.d/

# Third stage: Final image with Golang application and PostgreSQL
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Create a directory for the app
WORKDIR /app

# Copy the Golang binary from the builder stage
COPY --from=builder /app/main /app/

# Expose the application's port
EXPOSE 8080

# Set the timezone (optional)
ENV TZ=Etc/UTC

# Start the Golang application
CMD ["/app/main"]
