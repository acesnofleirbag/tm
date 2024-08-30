DESTDIR ?= /usr/local/bin/

.PHONY: build
build:
	go build -o build/tm

.PHONY: install
install: build
	sudo cp build/tm "$(DESTDIR)"

