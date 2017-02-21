SHELL := /bin/bash
DOCKER_PRESENT := $(shell command -v docker 2> /dev/null)
GLIDE_PRESENT := $(shell command -v glide 2> /dev/null)
BINARY := elasticsearch-cli
VERSION ?= 0.1.0
export CGO_ENABLED := 0

.PHONY: deps
deps:
ifndef GLIDE_PRESENT
	@ curl -sL https://glide.sh/get | bash
endif
	@ go get github.com/mitchellh/gox

.PHONY: build
build: deps
	@ echo "-> Installing $(BINARY) dependencies..."
	@ glide install
	@ echo "-> Building $(BINARY)..."
	@ gox -os "darwin linux" -arch="386 amd64" \
		-output="pkg/{{.Dir}}_{{.OS}}_{{.Arch}}"

.PHONY: docker-build
docker-build:
	@ echo "-> Running $(BINARY) build in Docker..."
	@ docker run --rm \
	-v $(shell pwd):/go/src/github.com/marclop/$(BINARY) \
	golang:1.7-alpine \
	sh -c 'apk --update add curl bash git make && cd /go/src/github.com/marclop/$(BINARY) && make build'

.PHONY: _get_sys_arch
_get_sys_arch:
	$(eval OS_NAME := $(shell echo $(shell uname -s) | tr '[:upper:]' '[:lower:]'))
	$(eval ARCH := "386")
ifeq ($(shell uname -m), x86_64)
	$(eval ARCH := "amd64")
endif

.PHONY: install
install: docker-build _get_sys_arch
	@ echo "-> Moving binary to /usr/local/bin/$(BINARY)"
	@ mv pkg/$(BINARY)_$(OS_NAME)_$(ARCH) /usr/local/bin/$(BINARY)

.PHONY: release
release: build
	@ go get github.com/tcnksm/ghr
	@ echo "-> Publishing $(BINARY) to GitHub..."
	@ ghr -u elastic $(VERSION) pkg
