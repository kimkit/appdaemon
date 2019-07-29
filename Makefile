VERSION_TAG:=$(shell git describe --always --tag)

all: fmt
	rm -rf bin/appdaemon.*
	go build -o bin/appdaemon cmd/appdaemon/main.go

simple:
	rm -rf bin/appdaemon.*
	go build -o bin/appdaemon cmd/appdaemon-simple/main.go

release: fmt
	rm -rf bin/appdaemon.*
	GOOS=linux GOARCH=amd64 go build -o bin/appdaemon.linux.$(VERSION_TAG) cmd/appdaemon/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/appdaemon.darwin.$(VERSION_TAG) cmd/appdaemon/main.go
	cp bin/appdaemon.linux.$(VERSION_TAG) bin/appdaemon.linux

image: release
	docker build -t kimkit/appdaemon:$(VERSION_TAG) .

fmt:
	if [ -d static/dist ]; then statik -src static/dist -dest static; fi
	find . -name '*.go' | grep -v '^\./vendor/' | xargs -i go fmt {}

static:
	rm -rf static/dist
	cd static && npm install && npm run lint && npm run build

update:
	go list -m -u all

run: fmt static
	rm -rf bin/appdaemon.*
	GOOS=linux GOARCH=amd64 go build -o bin/appdaemon.linux cmd/appdaemon/main.go
	docker build -t appdaemon .
	docker-compose rm -f
	docker-compose up

runcluster: fmt static
	rm -rf bin/appdaemon.*
	GOOS=linux GOARCH=amd64 go build -o bin/appdaemon.linux cmd/appdaemon/main.go
	docker build -t appdaemon .
	docker-compose -f docker-compose-cluster.yaml rm -f
	docker-compose -f docker-compose-cluster.yaml up

.PHONY: static
