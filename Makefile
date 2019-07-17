all:
	find . -name '*.go' | grep -v '^\./vendor/' | xargs -i go fmt {}
	go build -o bin/appdaemon cmd/cmdsvr/main.go
	git describe --always --tag
	GOOS=linux GOARCH=amd64 go build -o bin/appdaemon.linux cmd/cmdsvr/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/appdaemon.darwin cmd/cmdsvr/main.go

update:
	go list -m -u all

run: all
	docker build -t appdaemon .
	docker-compose rm -f
	docker-compose up
