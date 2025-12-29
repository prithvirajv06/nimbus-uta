FROM golang:1.24-alpine AS builder

WORKDIR /app

# 1. ADD: Install 'build-base' for GCC and CGO dependencies.
# This fixes the "C compiler not found" error.
RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# The CGO_ENABLED=1 command will now succeed.
RUN go build -a -o main .

# --- START OF RUNNER STAGE ---
# This stage remains largely unchanged, ensuring a small final image.
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

# Ensure the binary has execute permissions
RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]