.PHONY: run deps build clean

run: build
	./bin/app

build:
	go build -o bin/app -ldflags '-s -w -extldflags "-static"'

clean:
	rm -rf bin

deps:
	go get -d -u -v ./...

