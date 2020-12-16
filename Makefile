PROGRAM=rmap

$(PROGRAM):
	go build ./cmd/rmap

.PHONY: clean
clean:
	rm -f $(PROGRAM)
