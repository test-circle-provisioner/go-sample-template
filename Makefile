VERSION := $(shell git describe --tags --always --dirty="-dev")
LDFLAGS := -ldflags='-X "main.version=$(VERSION)"'

Q=@

.PHONY: deps
deps:
	$Qdep ensure

.PHONY: vet
vet:
	$Qgo vet ./...

.PHONY: test
test: vet
	$Qgo test -cover -race ./...

.PHONY: build
build:
	$Qdocker build --build-arg VERSION=$(VERSION) \
		-t segment/scrubber:$(VERSION) \
		-t 528451384384.dkr.ecr.us-west-2.amazonaws.com/go-hello-world:$(VERSION) \
		.

.PHONY: release
release: build
	$Qdocker push segment/scrubber:$(VERSION)
	$Qdocker push 528451384384.dkr.ecr.us-west-2.amazonaws.com/go-hello-world:$(VERSION)
