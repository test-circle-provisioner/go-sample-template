VERSION := $(shell git describe --tags --always --dirty="-dev")

deps:
	govendor sync

vet:
	govendor vet ./...

test: vet
	govendor test -v --cover --race -coverprofile=profile.out -covermode=atomic

build:
	docker build --build-arg VERSION=$(VERSION) \
		-t segment/go-hello-world:$(VERSION) \
		-t 528451384384.dkr.ecr.us-west-2.amazonaws.com/go-hello-world:$(VERSION) \
		.

release: build
	docker push segment/go-hello-world:$(VERSION)
	docker push 528451384384.dkr.ecr.us-west-2.amazonaws.com/go-hello-world:$(VERSION)

.PHONY: *
