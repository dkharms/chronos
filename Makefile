SHELL := /bin/bash

OS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH ?= $(shell uname -m | tr '[:upper:]' '[:lower:]')

ifeq ($(ARCH), x86_64)
ARCH := amd64
endif

LOCAL_BIN:=$(CURDIR)/bin
.PHONY: build-binaries
build-binaries:
	CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build \
      -o ./bin/${OS}-${ARCH}/ \
      ./cmd/...

.PHONY: imports
imports:
	goimports -w -local 'github.com/dkharms/chronos' .
