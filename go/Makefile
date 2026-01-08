.PHONY: run build test clean docker-up docker-down docker-build

docker-build:
    docker-compose build

docker-up:
    docker-compose up -d --build

docker-down:
    docker-compose down

install:
    go mod download
    go mod tidy

lint:
    golangci-lint run