.PHONY: build test clean install

build:
	go build -o ts2go ./cmd/ts2go

test:
	cd tests && go test -v ./...

clean:
	go clean ./...
	rm -f ts2go

install:
	go install ./cmd/ts2go