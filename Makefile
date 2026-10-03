.DEFAULT_GOAL := test

export GO111MODULE := on
export SHELL := /bin/bash

GOPKG := github.com/veraison/da
GOPKG += github.com/veraison/da/dat
GOPKG += github.com/veraison/da/dat/cmd

# dat only holds the main() wrapper and has no tests
COVER_PKG := $(filter-out github.com/veraison/da/dat,$(GOPKG))

GOLINT ?= golangci-lint

GOLINT_ARGS ?= run

.PHONY: lint
lint: ; $(GOLINT) $(GOLINT_ARGS)

ifeq ($(MAKECMDGOALS),test)
GOTEST_ARGS ?= -v -race $(GOPKG)
else
  ifeq ($(MAKECMDGOALS),test-cover)
  GOTEST_ARGS ?= -short -cover $(COVER_PKG)
  endif
endif

COVER_THRESHOLD := $(shell sed -n "s/^ *min-coverage: '\(.*\)'/≥\1%/p" .github/workflows/ci-go-cover.yml)

.PHONY: test test-cover
test test-cover: ; go test $(GOTEST_ARGS)

.PHONY: presubmit
presubmit:
	@echo
	@echo ">>> Check that the reported coverage figures are $(COVER_THRESHOLD)"
	@echo
	$(MAKE) test-cover
	@echo
	@echo ">>> Fix any lint error"
	@echo
	$(MAKE) lint

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  * test:       run unit tests for $(GOPKG)"
	@echo "  * test-cover: run unit tests and measure coverage for $(COVER_PKG)"
	@echo "  * lint:       lint sources using .golangci.yml"
	@echo "  * presubmit:  check you are ready to push your local branch to remote"
	@echo "  * help:       print this menu"
