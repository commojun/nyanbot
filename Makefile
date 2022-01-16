PROJECT_NAME=nyanbot
VERSION?=0.0.5

testall:
	go test ./...

release:
	git tag -a v$(VERSION) -m "new tag by make"
	git push origin tag v$(VERSION)

dockerbuild:
	cd docker/ && \
	docker buildx build \
		-t commojun/nyanb:$(VERSION) \
		--platform linux/amd64,linux/arm/v7,linux/arm/v6 \
		--build-arg VERSION=$(VERSION) \
		--push \
		-f ./Dockerfile .
