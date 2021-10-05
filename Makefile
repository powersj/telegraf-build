.PHONY: all
all: build

.PHONY: build
build:
	go build

.PHONY: clean
clean:
	rm -f telegraf-mod

.PHONY: lint
lint:
	golangci-lint run
