.PHONY: build clean

BINARY_NAME=hello-world

build:
	go build -o $(BINARY_NAME) main.go

clean:
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME) --rpc
