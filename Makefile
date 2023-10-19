build:
	go build -o ./bin/currencies ./cmd/currencies/main.go

fmt:
	gofumpt -w .

tidy:	
	go mod tidy

lint: build fmt tidy
	golangci-lint run ./...

run:
	go run ./cmd/currencies/main.go
