# Start from the official Golang image
FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o notificationapp main.go
# Use a minimal image for running
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/notificationapp .
EXPOSE 8082
CMD ["./notificationapp"]