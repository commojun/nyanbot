PROJECT_NAME=nyanbot

testall:
	go test ./...

export:
	go run cmd/init/init.go
