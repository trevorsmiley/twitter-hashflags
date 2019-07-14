.PHONY: compile

build-all: mac windows

windows:
	GOOS=windows GOARCH=amd64 go build

mac:
	GOOS=darwin GOARCH=amd64 go build

test-all:
	go test -v ./...
