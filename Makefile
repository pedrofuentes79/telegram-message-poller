.PHONY: build run

build:
	go build -o bin/telegram-check ./src

run: build
	./bin/telegram-check
