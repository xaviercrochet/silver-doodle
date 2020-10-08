server:
	go run main.go -port 8080
test:
	go test -cover -race ./...
