VERSION := $(shell git describe --tags --always --dirty="-dev")
LDFLAGS := -ldflags='-X "main.version=$(VERSION)"'

Q=@

deps:
	$Qdep ensure

vet:
	$Qgo vet ./...

test: vet
	$Qgo test -v -cover -race ./...

build:
	$Qdocker build --build-arg VERSION=$(VERSION) \
		-t segment/scrubber:$(VERSION) \
		-t 528451384384.dkr.ecr.us-west-2.amazonaws.com/go-hello-world:$(VERSION) \
		.

release: build
	$Qdocker push segment/scrubber:$(VERSION)
	$Qdocker push 528451384384.dkr.ecr.us-west-2.amazonaws.com/go-hello-world:$(VERSION)

.PHONY: *
