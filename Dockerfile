# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum secret.env ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main ./cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /app

# Copy binary and env file into final image
COPY --from=builder /app/main .
COPY secret.env .

# Option A: Source the env file in the ENTRYPOINT/CMD
# This ensures environment variables are set before running your binary
CMD ["/bin/sh", "-c", "source /app/secret.env && exec ./main"]

EXPOSE 8080