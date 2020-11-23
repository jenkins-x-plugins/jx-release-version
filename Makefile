NAME := jx-release-version
ORG := jenkins-x-plugins
ROOT_PACKAGE := main.go

GO := GO15VENDOREXPERIMENT=1 go

PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
FORMATTED := $(shell $(GO) fmt $(PACKAGE_DIRS))

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_DIR ?= ./bin
BUILDFLAGS := '-w -s'

REV := $(shell git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
BRANCH     := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown')
BUILD_DATE := $(shell date +%Y%m%d-%H:%M:%S)
ORG_REPO := $(ORG)/$(NAME)
ROOT_PACKAGE := github.com/$(ORG_REPO)
GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')

all: test build

check: fmt test

.PHONY: build
build:
	CGO_ENABLED=0 GOARCH=amd64 go build -ldflags $(BUILDFLAGS) -o $(BUILD_DIR)/$(NAME) $(ROOT_PACKAGE)

fmt:
	@FORMATTED=`$(GO) fmt $(PACKAGE_DIRS)`
	@([[ ! -z "$(FORMATTED)" ]] && printf "Fixed unformatted files:\n$(FORMATTED)") || true

.PHONY: test
test:
	go test -v $(GOPACKAGES)

.PHONY: release
release: clean test

.PHONY: goreleaser
release:
	step-go-releaser --organisation=$(ORG) --revision=$(REV) --branch=$(BRANCH) --build-date=$(BUILD_DATE) --go-version=$(GO_VERSION) --root-package=$(ROOT_PACKAGE) --version=$(VERSION)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -rf dist

.PHONY: docker
docker: $(BUILD_DIR)/$(NAME)-linux
	docker build -t "${ORG}/$(NAME):dev" .
