PROJECT_NAME=nyanbot
BINARY_DIR=cmd/$(PROJECT_NAME)
BINARY_NAME=$(PROJECT_NAME)

all: test build
build:
	cd $(BINARY_DIR) && go build $(BINARY_NAME).go
test:
	go test
clean:
	go clean
	rm -f $(BINARY_DIR)/$(BINARY_NAME)
run: build
	cd $(BINARY_DIR) && ./$(BINARY_NAME)
deps:
	go get github.com/golang/dep/cmd/dep
	dep ensure
