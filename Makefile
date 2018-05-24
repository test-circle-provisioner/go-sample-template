VERSION := $(shell git rev-parse --short HEAD)_$(CIRCLE_BUILD_NUM)
LDFLAGS := -ldflags='-X "main.version=$(VERSION)"'

Q=@

.PHONY: deps
deps:
	$Qdep ensure

.PHONY: vet
vet:
	$Qgo vet ./...

.PHONY: generate
generate:
	$Qgo generate ./...

.PHONY: test
test: vet
	$Qgo test -cover -race ./...

.PHONY: build
build:
	$Qdocker build --build-arg VERSION=$(VERSION) \
		-t 528451384384.dkr.ecr.us-west-2.amazonaws.com/{{ .Name }}:$(VERSION) \
		.

.PHONY: release
release: build
	$Qdocker push 528451384384.dkr.ecr.us-west-2.amazonaws.com/{{ .Name }}:$(VERSION)
