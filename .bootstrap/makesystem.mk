ifndef BOOTSTRAP_MAKESYSTEM_MK
BOOTSTRAP_MAKESYSTEM_MK := 1

MAKESYSTEM_BASE_DIR ?= ./.makesystem

SHELL := /bin/bash -e -o pipefail

ifneq ("$(wildcard $(MAKESYSTEM_BASE_DIR)/.id)","")
    MAKESYSTEM_FOUND := 1
else
    MAKESYSTEM_FOUND := 0
endif

ifneq ($(MAKECMDGOALS),makesystem_install)
    ifneq ($(MAKESYSTEM_FOUND),1)
        $(error makesystem not installed, please install the makesystem by running "make makesystem_install")
    endif
endif

all:
.PHONY: all

makesystem_install:
	@./.bootstrap/setup-makesystem.sh "$(MAKESYSTEM_BASE_DIR)"
.PHONY: makesystem_install

endif
