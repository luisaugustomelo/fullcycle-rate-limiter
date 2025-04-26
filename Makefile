run:
	go run cmd/server/main.go

curl:
	curl -H "API_KEY: abc123" http://localhost:8080

test:
	go test ./internal/middleware -v 

