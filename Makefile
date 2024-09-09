build:
	@go build -o bin/ws main.go

run: build
	./bin/ws

