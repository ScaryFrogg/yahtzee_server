build:
	@echo "Building..."
	@go build -o main cmd/yahtzee/main.go
run:
	@go run cmd/yahtzee/main.go

.PHONY: build run
