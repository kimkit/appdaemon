all:
	find . -name '*.go' | grep -v '^\./vendor/' | gxargs -i go fmt {}
	go build -o bin/appdaemon cmd/cmdsvr/main.go

update:
	go list -m -u all
