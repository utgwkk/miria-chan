all: gofmt
	go build -o miria-chan main

gofmt:
	gofmt -w src/main/main.go
