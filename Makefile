.PHONY: build
build:
	mkdir -p bin
	go build -o ./bin/annotations ./...

.PHONY: run
run:
	go run ./...
