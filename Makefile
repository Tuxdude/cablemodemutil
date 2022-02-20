# Set the shell to bash.
SHELL := /bin/bash -e -o pipefail

# Enable a verbose output from the makesystem.
VERBOSE ?= no

# Disable the colorized output from make if either
# explicitly overridden or if no tty is attached.
DISABLE_COLORS ?= $(shell [ -t 0 ] && echo no)

# Silence echoing the commands being invoked unless
# overridden to be verbose.
ifneq ($(VERBOSE),yes)
    silent := @
else
    silent :=
endif

# Configure colors.
ifeq ($(DISABLE_COLORS),no)
    COLOR_BLUE    := \x1b[1;34m
    COLOR_RESET   := \x1b[0m
else
    COLOR_BLUE    :=
    COLOR_RESET   :=
endif

GEN_FILES         := dist/ coverage.out

# Common utilities.
ECHO                              := echo -e

# go and related binaries.
GO_CMD                            := go
GO_IMPORTS_CMD                    := goimports
GO_FMT_CMD                        := gofmt
GO_LINT_CMD                       := golint
GO_CI_LINT_CMD                    := golangci-lint
GO_RELEASER                       := goreleaser

# Commands invoked from rules.
GOBUILD                           := $(GO_CMD) build
GOSTRIPPEDBUILD                   := CGO_ENABLED=0 GOOS=linux $(GO_CMD) build -a -ldflags "-s -w" -installsuffix cgo
GOCLEAN                           := $(GO_CMD) clean
GOGENERATE                        := $(GO_CMD) generate
GOGET                             := $(GO_CMD) get
GOLIST                            := $(GO_CMD) list
GOMOD                             := $(GO_CMD) mod
GOTEST                            := $(GO_CMD) test -v
GOCOVERAGE                        := $(GO_CMD) test -v -race -coverprofile coverage.out -covermode atomic
GOVET                             := $(GO_CMD) vet
GOIMPORTS                         := $(GO_IMPORTS_CMD) -w
GOFMT                             := $(GO_FMT_CMD) -s -w
GOLINT                            := $(GO_LINT_CMD) -set_exit_status -min_confidence 0.200001
GOLINTAGG                         := $(GO_LINT_CMD) -set_exit_status -min_confidence 0
GOLANGCILINT                      := $(GO_CI_LINT_CMD) run
GOLANGCILINTAGG                   := $(GO_CI_LINT_CMD) run --enable-all
GORELEASERRELEASE                 := $(GO_RELEASER) release
GORELEASERCHECK                   := $(GO_RELEASER) check
INSTALL_GORELEASER_HOOK_PREREQS   := $(GO_CMD) install \
    github.com/golangci/golangci-lint/cmd/golangci-lint@latest
CLEAN_ALL                         := $(GOCLEAN) ./... && rm -rf $(GEN_FILES)

# Alternative for running golangci-lint, using docker instead:
# docker run \
#   --rm \
#   --tty \
#   --volume $$(pwd):/go-src:ro \
#   --workdir /go-src \
#   golangci/golangci-lint:v1.44.0 \
#   golangci-lint run

# Helpful functions
# ExecWithMsg
# $(1) - Message
# $(2) - Command to be executed
define ExecWithMsg
    $(silent)$(ECHO) "\n===  $(COLOR_BLUE)$(1)$(COLOR_RESET)  ==="
    $(silent)$(2)
endef

# List of packages in the current directory.
PKGS ?= $(shell $(GOLIST) ./... | grep -v /vendor/)
# Define tags.
TAGS ?=

DEP_PKGS := $(shell $(GOLIST) -f '{{ join .Imports "\n" }}' | grep tuxdude || true)
ifeq ($(DEP_PKGS),)
    DEP_PKGS_TEXT := None
else
    DEP_PKGS_TEXT := $(DEP_PKGS)
    DEP_PKGS := $(addsuffix @master,$(DEP_PKGS))
endif

all: fix_imports generate fmt lint vet build test
.PHONY: all

build: tidy
	$(call ExecWithMsg,Building,$(GOBUILD) ./...)
.PHONY: build

build_stripped: tidy
	$(call ExecWithMsg,Building Stripped,$(GOSTRIPPEDBUILD) ./...)
.PHONY: build_stripped

clean:
	$(call ExecWithMsg,Cleaning,$(CLEAN_ALL))
.PHONY: clean

coverage: tidy
	$(call ExecWithMsg,Testing with Coverage generation,$(GOCOVERAGE) ./...)
.PHONY: coverage

deps_list:
	$(call ExecWithMsg,Listing dependencies,$(GOLIST) -m all)
.PHONY: deps_list

deps_list_latest_version:
	$(call ExecWithMsg,Listing latest dependency versions,$(GOLIST) -u -m all)
.PHONY: deps_list_latest_version

deps_update_tuxdude_latest_only:
	$(call ExecWithMsg,Updating to the latest version of dependencies for \"$(DEP_PKGS_TEXT)\",$(GOGET) -t -u $(DEP_PKGS))
.PHONY: deps_update_tuxdude_latest_only

deps_update_tuxdude_latest: deps_update_tuxdude_latest_only tidy
.PHONY: deps_update_tuxdude_latest

deps_update_only:
	$(call ExecWithMsg,Updating to the latest version of all direct dependencies,$(GOGET) -t -u ./...)
.PHONY: deps_update_only

deps_update: deps_update_only tidy
.PHONY: deps_update

fix_imports:
	$(call ExecWithMsg,Fixing imports,$(GOIMPORTS) .)
.PHONY: fix_imports

fmt:
	$(call ExecWithMsg,Fixing formatting,$(GOFMT) .)
.PHONY: fmt

generate:
	$(call ExecWithMsg,Generating,$(GOCLEAN) ./...)
.PHONY: generate

goreleaser_check_config:
	$(call ExecWithMsg,GoReleaser Checking config,$(GORELEASERCHECK))
.PHONY: goreleaser_check_config

goreleaser_local_release:
	$(call ExecWithMsg,GoReleaser Building Local Release,$(GORELEASERRELEASE) --snapshot --rm-dist)
.PHONY: goreleaser_local_release

goreleaser_verify_install_prereqs:
	$(call ExecWithMsg,GoReleaser Pre-Release Installing Prereqs,$(INSTALL_GORELEASER_HOOK_PREREQS))
.PHONY: goreleaser_verify_install_prereqs

goreleaser_verify: goreleaser_verify_install_prereqs generate fmt lint vet build test
.PHONY: goreleaser_verify

lint: tidy
	$(call ExecWithMsg,Linting,$(GOLANGCILINT))
.PHONY: lint

lint_agg: tidy
	$(call ExecWithMsg,Aggressive Linting,$(GOLANGCILINTAGG))
.PHONY: lint_agg

lint_deprecated: tidy
	$(call ExecWithMsg,Linting (Deprecated),$(GOLINT) .)
.PHONY:lint_deprecated

lint_deprecated_agg: tidy
	$(call ExecWithMsg,Aggressive Linting (Deprecated),$(GOLINTAGG) .)
.PHONY:lint_deprecated_agg

test: tidy
	$(call ExecWithMsg,Testing,$(GOTEST) ./...)
.PHONY: test

tidy:
	$(call ExecWithMsg,Tidying module,$(GOMOD) tidy)
.PHONY: tidy

vet: tidy
	$(call ExecWithMsg,Vetting,$(GOVET) ./...)
.PHONY: vet
