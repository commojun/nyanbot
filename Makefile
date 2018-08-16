PROJECT_NAME=nyanbot
CMD=cmd
PUSH=nyanpush

test:
	go test
deps:
	go get github.com/golang/dep/cmd/dep
	dep ensure
install:
	go install ./$(CMD)/$(PUSH)
	@echo 'please add following crontab'
	@echo '*/5 * * * * path/to/nyanpush --config="/path/to/your/config.yml"'
