.SILENT:
.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.PHONY: run deps build clean exec

name=$(shell basename $(CURDIR))
# commit_hash=$(shell git log -1 --pretty=format:"%H")
# commit_date=$(shell git log -1 --pretty=format:"%cD")
# commit_timestamp=$(shell git log -1 --pretty=format:"%at")
# current_dir = $(shell pwd)

run: buildPublic build exec clean

exec:
	./bin/${name}

buildPublic:
	go-bindata -pkg statik -o ./pkg/statik/statik.go ./public

build:
	CGO_ENABLED=0 go build -o bin/${name} -ldflags '-s -w -extldflags "-static"'

clean:
	rm -rf bin
	rm -rf upload

deps:
	govendor init
	govendor add +e
	govendor update +v

dev:
	go get -u -v github.com/kardianos/govendor
	go get -u -v github.com/go-bindata/go-bindata/...

