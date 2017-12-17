SOURCES=main.go slack.go

all: gofmt
	go build -o miria-chan github.com/utgwkk/miria-chan

gofmt:
	gofmt -w $(SOURCES)
