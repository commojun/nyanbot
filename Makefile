PROJECT_NAME=nyanbot
VERSION=v0.0.3

testall:
	go test ./...

export:
	go run cmd/init/init.go

release:
	git tag -a $(VERSION) -m "new tag by make"
	git push origin tag $(VERSION)
