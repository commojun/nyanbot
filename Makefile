PROJECT_NAME=nyanbot

testall:
	go test ./...

init:
	go run cmd/init/init.go

build:
	gox \
	-os=linux \
	-arch="arm amd64" \
	-output=bin/{{.Dir}}_{{.Arch}} \
	./cmd/hello \
	./cmd/alarm \
	./cmd/echo \
	./cmd/init
