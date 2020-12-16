PROGRAM=rmap

$(PROGRAM):
	go build ./cmd/rmap

.PHONY: tests
tests:
	go test ./...

.PHONY: clean
clean:
	rm -f $(PROGRAM)
