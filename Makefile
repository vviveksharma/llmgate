.PHONY: build run test loadtest tidy docker-up docker-down

build:
	go build -o bin/gateway ./cmd/gateway
	go build -o bin/loadtest ./cmd/loadtest

run:
	go run ./cmd/gateway

test:
	go test ./...

loadtest:
	go run ./cmd/loadtest -addr http://localhost:8080 -n 200 -c 20

tidy:
	go mod tidy

docker-up:
	docker compose -f deploy/docker-compose.yml up --build

docker-down:
	docker compose -f deploy/docker-compose.yml down
