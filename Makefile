BINARY := opencode-configurator

.PHONY: all build run test vet fmt clean

all: fmt vet test build

build:
	go build -o $(BINARY) .

run:
	go run main.go

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

clean:
	rm -f $(BINARY)
