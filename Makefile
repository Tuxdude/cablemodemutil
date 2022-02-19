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

# Common utilities.
ECHO := echo -e

# go and related binaries.
GO_CMD            := go
GO_IMPORTS_CMD    := goimports
GO_FMT_CMD        := gofmt
GO_LINT_CMD       := golint
GO_CI_LINT_CMD    := golangci-lint
GO_RELEASER       := goreleaser

# Commands invoked from rules.
GOBUILD           := $(GO_CMD) build
GOSTRIPPEDBUILD   := CGO_ENABLED=0 GOOS=linux $(GO_CMD) build -a -ldflags "-s -w" -installsuffix cgo
GOCLEAN           := $(GO_CMD) clean
GOGENERATE        := $(GO_CMD) generate
GOGET             := $(GO_CMD) get -u
GOLIST            := $(GO_CMD) list
GOMOD             := $(GO_CMD) mod
GOTEST            := $(GO_CMD) test -v
GOCOVERAGE        := $(GO_CMD) test -v -race -coverprofile coverage.out -covermode atomic
GOVET             := $(GO_CMD) vet
GOIMPORTS         := $(GO_IMPORTS_CMD) -w
GOFMT             := $(GO_FMT_CMD) -s -w
GOLINT            := $(GO_LINT_CMD) -set_exit_status -min_confidence 0.200001
GOLINTAGG         := $(GO_LINT_CMD) -set_exit_status -min_confidence 0
GOLANGCILINT      := $(GO_CI_LINT_CMD) run
GORELEASERRELEASE := $(GO_RELEASER) release
GORELEASERCHECK   := $(GO_RELEASER) check
INSTALL_GORELEASER_HOOK_PREREQS   := $(GO_CMD) install \
    github.com/golangci/golangci-lint/cmd/golangci-lint@latest

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
PKGS ?= $(shell $(GO_CMD) list ./... | grep -v /vendor/)
# Define tags.
TAGS ?=

DEP_PKGS := $(shell $(GO_CMD) list -f '{{ join .Imports "\n" }}' | grep tuxdude || true)
ifeq ($(DEP_PKGS),)
    DEP_PKGS_TEXT := None
else
    DEP_PKGS_TEXT := $(DEP_PKGS)
    DEP_PKGS := $(addsuffix @master,$(DEP_PKGS))
endif

all: fiximports generate fmt lint vet build test

build: tidy
	$(call ExecWithMsg,Building,$(GOBUILD) ./...)

buildstripped: tidy
	$(call ExecWithMsg,Building Stripped,$(GOSTRIPPEDBUILD) ./...)

clean:
	$(call ExecWithMsg,Cleaning,$(GOCLEAN) ./...)

coverage: tidy
	$(call ExecWithMsg,Testing with Coverage generation,$(GOCOVERAGE) ./...)

deps_update:
	$(call ExecWithMsg,Updating to the latest version of dependencies for \"$(DEP_PKGS_TEXT)\",$(GOGET) $(DEP_PKGS))

fiximports:
	$(call ExecWithMsg,Fixing imports,$(GOIMPORTS) .)

fmt:
	$(call ExecWithMsg,Fixing formatting,$(GOFMT) .)

generate:
	$(call ExecWithMsg,Generating,$(GOCLEAN) ./...)

goreleaser_check_config:
	$(call ExecWithMsg,GoReleaser Checking config,$(GORELEASERCHECK))

goreleaser_local_release:
	$(call ExecWithMsg,GoReleaser Building Local Release,$(GORELEASERRELEASE) --snapshot --rm-dist)

goreleaser_verify_install_prereqs:
	$(call ExecWithMsg,GoReleaser Pre-Release Installing Prereqs,$(INSTALL_GORELEASER_HOOK_PREREQS))

goreleaser_verify: goreleaser_verify_install_prereqs generate fmt lint_golangci_lint_only vet build test

test: tidy
	$(call ExecWithMsg,Testing,$(GOTEST) ./...)

tidy:
	$(call ExecWithMsg,Tidying module,$(GOMOD) tidy)

lint: tidy
	$(call ExecWithMsg,Linting,$(GOLINT) . && $(GOLANGCILINT))

lint_agg: tidy
	$(call ExecWithMsg,Aggressive Linting,$(GOLINTAGG) . && $(GOLANGCILINT))

lint_golangci_lint_only: tidy
	$(call ExecWithMsg,Linting (golangci-lint only),$(GOLANGCILINT))

vet: tidy
	$(call ExecWithMsg,Vetting,$(GOVET) ./...)

.PHONY: all build buildstripped clean coverage deps_update fiximports
.PHONY: fmt generate goreleaser_check_config goreleaser_local_release
.PHONY: goreleaser_verify_install_prereqs goreleaser_verify test tidy lint
.PHONY: lint_agg lint_golangci_lint_only vet
