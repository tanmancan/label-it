.PHONY: all clean

BINDIR := ./bin
CMDMAIN := ./cmd/label-it.go

GOBUILD := go build
LINKERFLAG := -ldflags

TARGET_OS := darwin linux windows
TARGET_ARCH := amd64

GITBRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GITCOMMIT := $(shell git rev-parse HEAD)
GITSHORT := $(shell git rev-parse --short HEAD)

# Versioning is done via git tags using semantic versioning
VERSION := $(shell git describe --tags --abbrev=0)
ifndef VERSION
	VERSION := dev
endif

# YAML Api version are based on branch name and major versions v1, v2, etc.
# API Breaking changes will only happen on major releases
APIVERSION 	:= $(GITBRANCH)

LDFLAGS = -s -w
LDFLAGS += -X github.com/tanmancan/label-it/v1/internal/config.BuildVersion=$(VERSION)
LDFLAGS += -X github.com/tanmancan/label-it/v1/internal/config.APIVersion=$(APIVERSION)
LDFLAGS += -X github.com/tanmancan/label-it/v1/internal/config.GitSHA=$(GITCOMMIT)

SRC := $(wildcard ./cmd/*.go)
BINS := $(TARGET_OS:%=${BINDIR}/%-${TARGET_ARCH}-${VERSION}-${GITSHORT})
BINSDIR := ${BINDIR}/%-${TARGET_ARCH}-${VERSION}-${GITSHORT}

RELEASE := $(BINS:%=%.tar.gz)
BINZIP := %.tar.gz

all: test ${BINS} ${RELEASE}

TESTS := $(wildcard ./internal/*/*_test.go)
TESTDIR := $(dir $(TESTS))

test: $(TESTDIR)
	@go test $(TESTDIR)

# Build for various platforms. See TARGET_OS
${BINSDIR}: ${SRC}
	@echo Building for $*-${TARGET_ARCH}
	@mkdir -p $@
	@env GOOS=$* GOARCH=${TARGET_ARCH} go build -ldflags '${LDFLAGS}' -o $@ $<

# Package and compress binaries and docs for release
${BINZIP}: ./LICENSE ./readme.md
	@echo Creating archive $@
	@cp ./LICENSE $*/LICENSE
	@cp ./readme.md $*/readme.md
	@tar -zcf $@ -C $* .

clean:
	@echo "Cleaning Up Binaries..."
	@rm -rf ./bin/*

