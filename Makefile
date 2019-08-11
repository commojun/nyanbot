PROJECT_NAME=nyanbot

testall:
	go test ./...

init:
	go run cmd/init/init.go
