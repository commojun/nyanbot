PROJECT_NAME=nyanbot
VERSION=v0.0.3

testall:
	go test ./...

release:
	git tag -a $(VERSION) -m "new tag by make"
	git push origin tag $(VERSION)

dockerbuild:
	cd docker/ && \
	docker buildx build -t commojun/nyanbot:$(VERSION) --platform linux/amd64,linux/arm -f ./Dockerfile .
