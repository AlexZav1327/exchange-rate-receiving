build:
	go build -o ./bin/currencies-service ./cmd/currencies-service/main.go

fmt:
	gofumpt -w .

tidy:	
	go mod tidy

lint: build fmt tidy
	golangci-lint run ./...

run:
	go run ./cmd/currencies-service/main.go

up:
	docker compose up -d

down:
	docker compose down