# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO
RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -a -o main .

# --- Runner stage ---
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 5050

CMD ["./main"]

