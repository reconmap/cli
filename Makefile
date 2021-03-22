SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
LATEST_TAG = $(shell git describe --tags)

PROGRAM=rmap

$(PROGRAM):
	go build -v -ldflags="-X 'github.com/reconmap/cli/internal/build.BuildVersion=$(LATEST_TAG)'" ./cmd/rmap

.PHONY: tests
tests:
	go test ./...

.PHONY: clean
clean:
	rm -f $(PROGRAM)

