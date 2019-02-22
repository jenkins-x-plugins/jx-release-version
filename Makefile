NAME := jx-release-version
ORG := jenkins-x
ROOT_PACKAGE := main.go
VERSION := $(shell jx-release-version)

GO := GO15VENDOREXPERIMENT=1 go
REVISION        := $(shell git rev-parse --short HEAD 2> /dev/null  || echo 'unknown')
BRANCH     := $(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo 'unknown')
BUILD_DATE := $(shell date +%Y%m%d-%H:%M:%S)

GO_VERSION=$(shell go version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
FORMATTED := $(shell $(GO) fmt $(PACKAGE_DIRS))

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_DIR ?= ./bin

BUILDFLAGS := -ldflags \
  " -X $(ROOT_PACKAGE)/version.Version='$(VERSION)'\
    -X $(ROOT_PACKAGE)/version.Revision='$(REVISION)'\
    -X $(ROOT_PACKAGE)/version.Branch='$(BRANCH)'\
    -X $(ROOT_PACKAGE)/version.BuildDate='$(BUILD_DATE)'\
    -X $(ROOT_PACKAGE)/version.GoVersion='$(GO_VERSION)'\
    -s -w -extldflags '-static'"

all: test $(GOOS)-build

check: fmt test

.PHONY: build
build:
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILDFLAGS) -o $(BUILD_DIR)/$(NAME) $(ROOT_PACKAGE)

fmt:
	@FORMATTED=`$(GO) fmt $(PACKAGE_DIRS)`
	@([[ ! -z "$(FORMATTED)" ]] && printf "Fixed unformatted files:\n$(FORMATTED)") || true

darwin-build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build $(BUILDFLAGS) -o $(BUILD_DIR)/$(NAME)-darwin $(ROOT_PACKAGE)

linux-build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build $(BUILDFLAGS) -o $(BUILD_DIR)/$(NAME)-linux $(ROOT_PACKAGE)

windows-build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build $(BUILDFLAGS) -o $(BUILD_DIR)/$(NAME)-windows.exe $(ROOT_PACKAGE)

.PHONY: test
test:
	go test -v $(GOPACKAGES)

.PHONY: release
release: clean test cross
	mkdir -p release
	cp $(BUILD_DIR)/$(NAME)-* release
	gh-release checksums sha256
	gh-release create $(ORG)/$(NAME) $(VERSION) master v$(VERSION)
	updatebot push-version --kind docker JX_RELEASE_VERSION $(VERSION)

.PHONY: cross
cross: darwin-build linux-build windows-build

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -rf release

.PHONY: docker
docker: $(BUILD_DIR)/$(NAME)-linux
	docker build -t "${ORG}/$(NAME):dev" .
