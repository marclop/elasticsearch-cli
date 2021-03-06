SHELL := /bin/bash
BINARY := elasticsearch-cli
AUTHOR = marclop
VERSION ?= 0.3.0
INSTALL_PATH ?= $(GOPATH)/bin
# Test all major version of ES
ES_VERSION ?= 1.7 2.4 5.4 5.5 5.6
ES_PORT ?= 9200
ES_CONTAINER_NAME ?= elasticsearch-cli_es
ES_IMAGE ?= elasticsearch
ES_TAG_SUFFIX ?= -alpine
BUILD_PLATFORMS ?= "darwin linux"
BUILD_ARCHITECTURES ?= "386 amd64"
BUILD_OUTPUT ?= "pkg/{{.Dir}}_{{.OS}}_{{.Arch}}"
REPORT_PATH ?= reports
REPORT_FORMAT ?= html
COMPLETIONS_FILE ?= $(HOME)/.$(BINARY).auto
export GO111MODULE=on
define HELP

$(BINARY) v$(VERSION) Makefile
=================================

## Build target

- build:                  It will cross build $(BINARY) for $(BUILD_PLATFORMS).
- install:                It will install $(BINARY) in the current system (by default in $(INSTALL_PATH)/$(BINARY)).
- deps:                   It will install golint and gox.

## Development targets

- vendor:                 Installs vendor dependencies.
- start-es:               Starts Elasticsearch containers ($(shell echo $(ES_VERSION) | tr " " ","m)).
- stop-es:                Stops Elasticsearch containers ($(shell echo $(ES_VERSION) | tr " " ","m)).
- unit:                   Runs unit tests.
- acceptance:             Runs acceptance tests.
- test:                   Runs all tests.
- code-quality:           Runs a code quality against $(elasticsearch-cli).
- get-quality-report:     Returns the path for the latest code-quality report.

## Release targets

- release:                Releases the package to GitHub.

endef
export HELP

.DEFAULT: help
.PHONY: help
help:
	@ echo "$$HELP"

.PHONY: deps
deps:
	@ go get -u github.com/golang/lint/golint
	@ go get -u github.com/mitchellh/gox

.PHONY: vendor
vendor: deps
	@ echo "-> Installing $(BINARY) dependencies..."
	@ go get

.PHONY: docker-build
docker-build:
	@ echo "-> Building $(BINARY) inside Docker..."
	@ docker run --rm \
	-v $(shell pwd):/go/src/github.com/$(AUTHOR)/$(BINARY) \
	golang:1.9-alpine \
	sh -c 'apk --update add curl bash git make && cd /go/src/github.com/$(AUTHOR)/$(BINARY) && make build'

.PHONY: build
build: deps vendor
	@ echo "-> Building $(BINARY)..."
	@ gox -os=$(BUILD_PLATFORMS) -arch=$(BUILD_ARCHITECTURES) \
		-ldflags="-X main.Version=$(VERSION)" -output=$(BUILD_OUTPUT)

.PHONY: _set_build_current_arch
_set_build_current_arch:
	$(eval BUILD_PLATFORMS := $(shell echo $(shell uname -s) | tr '[:upper:]' '[:lower:]'))
	$(eval BUILD_OUTPUT := bin/$(BINARY))
	$(eval BUILD_ARCHITECTURES := "386")
ifeq ($(shell uname -m), x86_64)
	$(eval BUILD_ARCHITECTURES := "amd64")
endif

.PHONY: install
install: _set_build_current_arch build
	@ echo "-> Moving binary to $(INSTALL_PATH)/$(BINARY)"
	@ mv bin/$(BINARY) $(INSTALL_PATH)/$(BINARY)

.PHONY: release
release: build
	@ go get -u github.com/tcnksm/ghr
	@ echo "-> Publishing $(BINARY) to GitHub..."
	@ ghr -u $(AUTHOR) $(VERSION) pkg

.PHONY: start-es
start-es:
	$(eval ES_CONTAINER_NAME := $(ES_CONTAINER_NAME)_dev)
	@ $(foreach es_version,$(ES_VERSION),$(MAKE) start_elasticsearch_docker \ES_CONTAINER_NAME=$(ES_CONTAINER_NAME) \
	ES_VERSION=$(es_version) ES_PORT=920$(shell echo $(es_version) | cut -d '.' -f1);)
	@ echo "-> Port bindings are:"
	@ $(foreach es_version,$(ES_VERSION), echo "-> Elasticsearch $(es_version) => 920$(shell echo $(es_version) | cut -d '.' -f1)";)

.PHONY: stop-es
stop-es:
	$(eval ES_CONTAINER_NAME := $(ES_CONTAINER_NAME)_dev)
	@ $(foreach es_version,$(ES_VERSION),echo "Stopped $$(docker kill $(ES_CONTAINER_NAME)_$(es_version))";)

.PHONY: test
test: lint unit acceptance

.PHONY: unit
unit:
	@ echo "-> Running unit tests for $(BINARY)..."
	@ echo "" > coverage.txt
	@ for d in $(shell go list ./...); do \
			go test -race -coverprofile=coverage.txt -covermode=atomic $$d; \
			if [ -f profile.out ]; then \
				cat profile.out >> coverage.txt ; \
				rm profile.out ; \
			fi ; \
		done;

.PHONY: acceptance
acceptance: _set_build_current_arch build
	@ $(foreach es_version,$(ES_VERSION),$(MAKE) acc ES_VERSION=$(es_version);)

.PHONY: lint
lint: deps
	@ golint -set_exit_status $(shell go list ./...)

.PHONY: acc
acc: start_elasticsearch_docker
	@ echo "-> Running acceptance tests for $(BINARY) in Elasticsearch $(ES_VERSION)..."
	@ go test -tags acceptance -count 1 . || (echo "-> Killing Docker container $$( docker kill $(ES_CONTAINER_NAME)_$(ES_VERSION) )" && exit 1) \
	&& (echo "-> Killing Docker container $$(docker kill $(ES_CONTAINER_NAME)_$(ES_VERSION))")

.PHONY: start_elasticsearch_docker
start_elasticsearch_docker:
	@ printf "=> Starting Elasticsearch $(ES_VERSION)... "
	@ docker run -d --rm -p '$(ES_PORT):9200' --name $(ES_CONTAINER_NAME)_$(ES_VERSION) $(ES_IMAGE):$(ES_VERSION)$(ES_TAG_SUFFIX) > /dev/null
	@ while ! docker logs $(ES_CONTAINER_NAME)_$(ES_VERSION) | grep recovered > /dev/null; do sleep 1; done
	@ echo "Done."

.PHONY: code-quality
code-quality:
	@ go get -u github.com/wgliang/goreporter
	@ mkdir -p $(REPORT_PATH)
	@ goreporter -p ../$(shell basename $(PWD)) -r $(REPORT_PATH) -f $(REPORT_FORMAT)
	@ $$(open $$(make get-quality-report))

.PHONY: get-quality-report
get-quality-report:
	@ ls -dt1 $(REPORT_PATH)/*.html | head -1
