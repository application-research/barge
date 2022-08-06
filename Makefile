SHELL=/usr/bin/env bash

GO_BUILD_IMAGE?=golang:1.18

# VERSION is the nearest tag plus number of commits and short form of most recent commit since the tag, if any
VERSION=$(shell git describe --always --tag --dirty)

unexport GOFLAGS

CLEAN:=
BINS:=
GOFLAGS:=

## FFI
FFI_PATH:=extern/filecoin-ffi/
FFI_DEPS:=.install-filcrypto
FFI_DEPS:=$(addprefix $(FFI_PATH),$(FFI_DEPS))

$(FFI_DEPS): build/.filecoin-install ;

build/.filecoin-install: $(FFI_PATH)
	$(MAKE) -C $(FFI_PATH) $(FFI_DEPS:$(FFI_PATH)%=%)
	@touch $@

MODULES+=$(FFI_PATH)
BUILD_DEPS+=build/.filecoin-install
CLEAN+=build/.filecoin-install

ffi-version-check:
	@[[ "$$(awk '/const Version/{print $$5}' extern/filecoin-ffi/version.go)" -eq 3 ]] || (echo "FFI version mismatch, update submodules"; exit 1)
BUILD_DEPS+=ffi-version-check

.PHONY: ffi-version-check

$(MODULES): build/.update-modules ;
# dummy file that marks the last time modules were updated
build/.update-modules:
	git submodule update --init --recursive
ifneq ($(FFI_COMMIT_HASH),"")
	git submodule update --init --recursive && cd extern/filecoin-ffi/ && git checkout -q $(FFI_COMMIT_HASH)
endif
	touch $@

CLEAN+=build/.update-modules

# Add version information to the package
ldflags=-X=main.appVersion=$(VERSION)
ifneq ($(strip $(LDFLAGS)),)
	ldflags+=-extldflags=$(LDFLAGS)
endif

GOFLAGS+=-ldflags="$(ldflags)"

.PHONY: all
all: build

.PHONY: debug
debug: debug-build

.PHONY: deps
deps: $(BUILD_DEPS)

.PHONY: build
build: deps barge

.PHONY: debug-build
debug-build: deps debug-barg

.PHONY: barge
barge:
	go build $(GOFLAGS) -o barge .
BINS+=barge

.PHONY: debug-barg
debug-barg:
	go build -ldflags="all=-w" -o barge .
BINS+=barge


.PHONY: tests
tests:
	cd test
	go test -v $(GOFLAGS) ./...