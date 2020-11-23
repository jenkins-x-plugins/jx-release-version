NAME := jx-release-version
ORG := jenkins-x
ROOT_PACKAGE := main.go

GO := GO15VENDOREXPERIMENT=1 go

PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
FORMATTED := $(shell $(GO) fmt $(PACKAGE_DIRS))

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_DIR ?= ./bin
BUILDFLAGS := '-w -s'

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
	goreleaser release

.PHONY: goreleaser
release:
	goreleaser release

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -rf dist

.PHONY: docker
docker: $(BUILD_DIR)/$(NAME)-linux
	docker build -t "${ORG}/$(NAME):dev" .
