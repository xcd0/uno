BUILDDIR     := .
VERSION      := 0.0.1
REVISION     := `git rev-parse --short HEAD`
FLAG         := -ldflags='-X main.version=$(VERSION) -X main.revision='$(REVISION)' -s -w -extldflags="-static" -buildid=' -a -tags netgo -installsuffix -trimpath
MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# Detect OS
ifeq ($(OS),Windows_NT)
    EXE = .exe
else
    EXE =
endif

all: build
build:
	@echo "Building..."
	go build -C cmd/uno -o $(MAKEFILE_DIR)uno$(EXE)
	@echo "Built successfully."
release:
	@echo "Building..."
	go build -C cmd/uno -o $(MAKEFILE_DIR)uno$(EXE) $(FLAG)
	@echo "Built successfully."
	upx --lzma $(BUILDDIR)/uno$(EXE)

# PHONY targets
.PHONY: all server client clean release

