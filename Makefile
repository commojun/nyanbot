PROJECT_NAME=nyanbot
VERSION=0.0.3

testall:
	go test ./...

release:
	git tag -a v$(VERSION) -m "new tag by make"
	git push origin tag v$(VERSION)

dockerbuild:
	cd docker/ && \
	docker buildx build \
		-t commojun/nyanbot:$(VERSION) \
		--platform linux/amd64,linux/arm \
		--build-arg VERSION=$(VERSION) \
		-f ./Dockerfile .
