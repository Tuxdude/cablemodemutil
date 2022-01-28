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
GO_CMD         := go
GO_IMPORTS_CMD := goimports
GO_FMT_CMD     := gofmt
GO_LINT_CMD    := golint
GO_CI_LINT_CMD := golangci-lint

# Commands invoked from rules.
GOCLEAN        := $(GO_CMD) clean
GOMOD          := $(GO_CMD) mod
GOIMPORTS      := $(GO_IMPORTS_CMD) -w
GOFMT          := $(GO_FMT_CMD) -s -w
GOLINT         := $(GO_LINT_CMD) -set_exit_status -min_confidence 0.200001
GOLINTAGG      := $(GO_LINT_CMD) -set_exit_status -min_confidence 0
GOLANGCILINT   := $(GO_CI_LINT_CMD) run
GOVET          := $(GO_CMD) vet
GOBUILD        := $(GO_CMD) build
GOTEST         := $(GO_CMD) test -v
GOCOVERAGE     := $(GO_CMD) test -v -race -coverprofile coverage.out -covermode atomic

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

all: fiximports fmt lint vet build test

clean:
	$(call ExecWithMsg,Cleaning,$(GOCLEAN))

fiximports:
	$(call ExecWithMsg,Fixing imports,$(GOIMPORTS) .)

fmt:
	$(call ExecWithMsg,Fixing formatting,$(GOFMT) .)

lint:
	$(call ExecWithMsg,Linting,$(GOLINT) . && $(GOLANGCILINT))

lintagg:
	$(call ExecWithMsg,Aggressive Linting,$(GOLINTAGG) . && $(GOLANGCILINT))

vet:
	$(call ExecWithMsg,Vetting,$(GOVET) ./...)

tidy:
	$(call ExecWithMsg,Tidying module,$(GOMOD) tidy)

build: tidy
	$(call ExecWithMsg,Building,$(GOBUILD) ./...)

test: tidy
	$(call ExecWithMsg,Testing,$(GOTEST) ./...)

coverage: tidy
	$(call ExecWithMsg,Testing with Coverage generation,$(GOCOVERAGE) ./...)

.PHONY: all clean tidy fiximports fmt lint lintagg vet build test
