.PHONY: build test run-scan

build:
	go build -o zt ./cmd/zt
	go build -o plugins/plugin-os/plugin-os ./plugins/plugin-os

test:
	go test ./...

run-scan: build
	./zt scan
