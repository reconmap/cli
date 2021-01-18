SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

PROGRAM=rmap

$(PROGRAM):
	go build ./cmd/rmap

.PHONY: tests
tests:
	go test ./...

.PHONY: clean
clean:
	rm -f $(PROGRAM)
