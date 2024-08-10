.PHONY: build 

all: build
	
build:
	go build -o build/server ./cmd/server/*.go 

test:
	go test -count 1 ./...

help:
	@echo "make build: \t\t compile packages"
	@echo "make test: \t\t run units for the all packages"