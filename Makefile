# Copyright 2021 VMware, Inc.
# SPDX-License-Identifier: BSD-2

HASH := $(shell git rev-parse --short HEAD)
COMMIT_DATE := $(shell git show -s --format=%ci ${HASH})
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
VERSION := ${HASH} (${COMMIT_DATE})

BUILDDIR ?= .
SRCDIR ?= .

.PHONY: help
help:
	@echo "make [TARGETS...]"
	@echo
	@echo "    help:               Print this usage information."
	@echo "    build:              Builds project"
	@echo "    install:            Installs binary, configuration and unit files"
	@echo "    clean:              Cleans the build"

$(BUILDDIR)/:
	mkdir -p "$@"

$(BUILDDIR)/%/:
	mkdir -p "$@"

.PHONY: build
build:
	- mkdir -p bin
	go build -ldflags="-X 'main.buildVersion=${VERSION}' -X 'main.buildDate=${BUILD_DATE}'" -o bin/photon-os-containerd ./cmd/photon-os-containerd
	go build -ldflags="-X 'main.buildVersion=${VERSION}' -X 'main.buildDate=${BUILD_DATE}'" -o bin/containerctl ./cmd/containerctl

.PHONY: install
install:
	install bin/containerctl /usr/bin/
	install bin/photon-os-containerd /usr/bin/

	install -vdm 755 /etc/photon-os-container
	install -m 755 conf/photon-os-container.toml /etc/photon-os-container

	install -m 0644 units/photon-os-container.service /lib/systemd/system/
	systemctl daemon-reload

PHONY: clean
clean:
	go clean
	rm -rf bin
