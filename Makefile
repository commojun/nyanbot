PROJECT_NAME=nyanbot
BINARY_DIR=cmd/$(PROJECT_NAME)
BINARY_NAME=$(PROJECT_NAME)

all: test build
build:
	cd $(BUNARY_DIR)
	go build -o $(BINARY_NAME) -v
test:
	go test -v ./...
clean:
	go clean
	rm -f $(BINARY_DIR)/$(BINARY_NAME)
run:
	build
	./$(BINARY_NAME)
deps:
	go get github.com/golang/dep/cmd/dep
	dep ensure
