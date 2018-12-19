#PACKAGE  = github.com/mexisme/get-secrets
PACKAGE = $(shell test -f go.mod && head -1 go.mod |sed -e 's/^module *//')
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell \
            git describe --tags --always --dirty 2>/dev/null || \
			cat $(CURDIR)/.version 2>/dev/null)
BINDIR   = $(PWD)/bin
BIN      = $(BINDIR)/$(notdir $(PACKAGE))
PKGS     = $(or $(PKG),$(shell $(GO) list ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS = $(shell $(GO) list -f '{{ if .TestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))

DOCKER_TAG = ...CHANGEME...

COVERAGE_PROFILE = all.coverprofile
# COVERAGE_XML = coverage.xml
COVERAGE_HTML = coverage.index.html

ifneq ($(DATE),)
LDFLAGS_DATE = -X $(PACKAGE)/version.buildDate=$(DATE)
endif
ifneq ($(VERSION),)
LDFLAGS_VERSION = -X $(PACKAGE)/version.release=$(VERSION)
endif

GO      = go
GO_TEST = $(GINKGO) -r -p
GODOC   = godoc
GOFMT   = gofmt
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m>>\033[0m")

.PHONY: all all-debug strip
strip: LDFLAGS_STRIP = -s -w
all all-debug: dep-update fmt lint test-default build

build strip: ; $(info $(M) building executable...) @ ## Build program binary
	$Q env GOOS="$(GOOS)" GOARCH="$(GOARCH)" $(GO) build \
		-tags release \
		-ldflags '$(LDFLAGS_STRIP) $(LDFLAGS_VERSION) $(LDFLAGS_DATE)' \
		-o $(BIN) .

# Cross-compile for Linux:
build-linux strip-linux: GOOS=linux
build-linux strip-linux: GOARCH=amd64
build-linux strip-linux: BIN:=$(BIN)-linux
build-linux: build
strip-linux: strip

# Tools

GOLINT = golangci-lint
$(GOLINT): ; $(info $(M) building $(GOLINT)...)
	$Q cd && GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

GINKGO = ginkgo
$(GINKGO): ; $(info $(M) building $(GINKGO)...)
	$Q cd && GO111MODULE=off go get -u github.com/onsi/ginkgo/ginkgo

GOCOVMERGE = gocovmerge
$(GOCOVMERGE): ; $(info $(M) building $(GOCOVMERGE)...)
	$Q cd && GO111MODULE=off go get -u github.com/wadey/gocovmerge

GOCOV = gocov
$(GOCOV): ; $(info $(M) building $(GOCOV)...)
	$Q cd && GO111MODULE=off go get -u github.com/axw/gocov/...

.PHONY: $(GOLINT) $(GINKGO) $(GOCOVMERGE) $(GOCOV)

# Tests

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
INTEGRATION_TEST_TARGETS := test-integration test-integration-verbose

.PHONY: $(TEST_TARGETS) $(INTEGRATION_TEST_TARGETS) \
	check test tests run-test
## Run benchmarks
test-bench: ARGS=-run=__absolutelynothing__ -bench=.
## Run only short tests
test-short: ARGS=-short
## Run tests in verbose mode with coverage reporting
test-verbose test-integration-verbose: ARGS=-v -trace -cover
## Run tests with race detector
test-race: ARGS=-race
$(TEST_TARGETS): SKIP_ARGS=-skip=Integration
$(INTEGRATION_TEST_TARGETS): SKIP_ARGS=-focus=Integration
$(TEST_TARGETS) $(INTEGRATION_TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS) $(INTEGRATION_TEST_TARGETS): test
check test tests: run-test
run-test: | $(GINKGO) ; $(info $(M) running $(NAME:%=% )tests...) @ ## Run tests
	$Q $(GO_TEST) $(ARGS) $(SKIP_ARGS) $(TESTPKGS)

docker-build:
	echo "docker build --tag $(DOCKER_TAG) ."

.PHONY: cover coverage
cover coverage clean: COVERAGE_DIR:=$(PWD)
cover coverage clean: COVERAGE_FILES=$(shell find . -name '*.coverprofile')
cover coverage: | $(GOCOVMERGE) $(GOCOV)
	$Q $(GOCOVMERGE) $(COVERAGE_FILES) > $(COVERAGE_DIR)/$(COVERAGE_PROFILE)
	$Q $(GO) tool cover -func=$(COVERAGE_DIR)/$(COVERAGE_PROFILE)
	$Q $(GO) tool cover -html=$(COVERAGE_DIR)/$(COVERAGE_PROFILE) -o $(COVERAGE_DIR)/$(COVERAGE_HTML)

# Code format

.PHONY: lint
lint: $(GOLINT) ; $(info $(M) running $(GOLINT)...) @ ## Run golint
	$Q $(GOLINT) run ./...

.PHONY: fmt
fmt: ; $(info $(M) running gofmt...) @ ## Run gofmt on all source files
	$Q go fmt ./...

# Dependency management

.PHONY: dep-update
go.sum: go.mod
vendor: go.sum ; $(info $(M) retrieving dependencies...)
	$Q go mod vendor
	@touch $@
dep-update: go.mod ; $(info $(M) updating dependencies...)
	$Q go mod tidy

# Misc

.PHONY: clean clean-go clean-rm-files
clean: clean-go clean-rm-files dep-update ; $(info $(M) cleaning...)	@ ## Cleanup everything
clean-go:
	$Q go clean -v -modcache
clean-rm-files:
	$Q rm -rf $(BINDIR) $(COVERAGE_FILES) $(COVERAGE_DIR)/$(COVERAGE_PROFILE) $(COVERAGE_DIR)/$(COVERAGE_HTML)

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
