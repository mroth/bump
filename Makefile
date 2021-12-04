# This makefile is currently just for convenience in my own local development,
# you probably don't want to be using it otherwise.
NAME=bump

bin/$(NAME): *.go bin
	go build -o $@

bin:
	mkdir -p bin

snapshot:
	goreleaser release --rm-dist --snapshot

package:
	# never want to actually publish from local, thats what CI is for
	goreleaser release --rm-dist --skip-publish

clean:
	rm -rf bin
	rm -rf dist

install:
	go install

uninstall:
	rm ${GOPATH}/bin/$(NAME)

.PHONY: snapshot package clean install uninstall
