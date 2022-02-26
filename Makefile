#
# some housekeeping tasks
#

# variable definitions
NAME := certcal
DESC := provide an iCal web feed for certificate expiration
PREFIX ?= usr/local
VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)

BUILD_GOOS ?= $(shell go env GOOS)
BUILD_GOARCH ?= $(shell go env GOARCH)

RELEASE_ARTIFACTS_DIR := .release_artifacts
CHECKSUM_FILE := checksums.txt

$(RELEASE_ARTIFACTS_DIR):
	install -d $@

BUILDER := $(shell echo "${BUILDER_NAME} <${EMAIL}>")

PKG_RELEASE ?= 1
PROJECT_URL := "https://github.com/mrtazz/$(NAME)"
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.goversion=$(GOVERSION)'

TARGETS := certcal
INSTALLED_TARGETS = $(addprefix $(PREFIX)/bin/, $(TARGETS))

certcal: certcal.go
	GOOS=$(BUILD_GOOS) GOARCH=$(BUILD_GOARCH) go build -ldflags "$(LDFLAGS)" -o $@ $<

.PHONY: all
all: $(TARGETS) $(MAN_TARGETS)
.DEFAULT_GOAL:=all

# development tasks
.PHONY: test
test:
	go test -v ./...

.PHONY: coverage
coverage:
	go test -v -race -coverprofile=cover.out ./...
	@-go tool cover -html=cover.out -o cover.html

.PHONY: benchmark
benchmark:
	@echo "Running tests..."
	@go test -bench=. ${NAME}

# install tasks
$(PREFIX)/bin/%: %
	install -d $$(dirname $@)
	install -m 755 $< $@

.PHONY: install
install: $(INSTALLED_TARGETS) $(INSTALLED_MAN_TARGETS)

.PHONY: local-install
local-install:
	$(MAKE) install PREFIX=usr/local

.PHONY: build-artifact
build-artifact: certcal $(RELEASE_ARTIFACTS_DIR)
	mv certcal $(RELEASE_ARTIFACTS_DIR)/certcal-$(VERSION).$(BUILD_GOOS).$(BUILD_GOARCH)
	cd $(RELEASE_ARTIFACTS_DIR) && shasum -a 256 certcal-$(VERSION).$(BUILD_GOOS).$(BUILD_GOARCH) >> $(CHECKSUM_FILE)

.PHONY: github-release
github-release:
	gh release create $(VERSION) --title 'Release $(VERSION)' --notes-file docs/releases/$(VERSION).md $(RELEASE_ARTIFACTS_DIR)/*

# clean up tasks
.PHONY: clean
clean:
	git clean -fdx
