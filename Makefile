
.PHONY: neuro
neuro:
	go build -o neuro

.PHONY: install
install:
	install -D neuro ~/.local/bin

.PHONY: clean
clean:
	rm -f neuro
