build:
	@go build -o bin/DisBlocker
run:
	@./bin/DisBlocker
test:
	@go test -v ./...