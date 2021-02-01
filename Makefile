NAME := jx-release-version
ORG := jenkins-x-plugins
ORG_REPO := $(ORG)/$(NAME)
ROOT_PACKAGE := github.com/$(ORG_REPO)/v2
MAIN_SRC_FILE := main.go

GO := GO15VENDOREXPERIMENT=1 go
CGO_ENABLED = 0

PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
FORMATTED := $(shell $(GO) fmt $(PACKAGE_DIRS))

# set dev version unless VERSION is explicitly set via environment
VERSION ?= $(shell echo "$$(git for-each-ref refs/tags/ --count=1 --sort=-version:refname --format='%(refname:short)' 2>/dev/null)-dev+$(REV)" | sed 's/^v//')

REV        := $(shell git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
BRANCH     := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown')
BUILD_DATE := $(shell date +%Y%m%d-%H:%M:%S)
GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_DIR ?= ./bin
BUILDFLAGS := -ldflags \
  " -X main.Version=$(VERSION)\
    -X main.Revision='$(REV)'\
    -X main.Date='$(BUILD_DATE)'"

all: test build

check: fmt test

.PHONY: build
build:
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILDFLAGS) -o $(BUILD_DIR)/$(NAME) $(MAIN_SRC_FILE)

linux: ## Build for Linux
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 go build $(BUILDFLAGS) -o build/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x build/linux/$(NAME)

fmt:
	@FORMATTED=`$(GO) fmt $(PACKAGE_DIRS)`
	@([[ ! -z "$(FORMATTED)" ]] && printf "Fixed unformatted files:\n$(FORMATTED)") || true

.PHONY: test
test:
	go test -v ./...

.PHONY: release
release: clean test linux

.PHONY: goreleaser
goreleaser:
	step-go-releaser --organisation=$(ORG) --revision=$(REV) --branch=$(BRANCH) --build-date=$(BUILD_DATE) --go-version=$(GO_VERSION) --root-package=$(ROOT_PACKAGE) --version=$(VERSION)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -rf dist

.PHONY: docker
docker: $(BUILD_DIR)/$(NAME)-linux
	docker build -t "${ORG}/$(NAME):dev" .
