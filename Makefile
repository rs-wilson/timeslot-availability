.PHONY: all build vet test clean run

all: clean build vet test

build: clean
		go build -o timeslot-server main.go

vet:
		go vet ./...

test: build
		go test ./...

clean:
		rm -f timeslot-server

run: build
	./timeslot-server
