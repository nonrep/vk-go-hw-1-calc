# Makefile

BINARY=calc

all: build

build:
	go build -o $(BINARY) .

.PHONY: all build