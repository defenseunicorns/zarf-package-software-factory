# Figure out which Zarf binary we should use based on the operating system we are on
ZARF_BIN := zarf
UNAME_S := $(shell uname -s)
UNAME_P := $(shell uname -p)
ifneq ($(UNAME_S),Linux)
	ifeq ($(UNAME_S),Darwin)
		ZARF_BIN := $(addsuffix -mac,$(ZARF_BIN))
	endif
	ifeq ($(UNAME_P),i386)
		ZARF_BIN := $(addsuffix -intel,$(ZARF_BIN))
	endif
	ifeq ($(UNAME_P),arm64)
		ZARF_BIN := $(addsuffix -apple,$(ZARF_BIN))
	endif
endif

.DEFAULT_GOAL := help

# Idiomatic way to force a target to always run, by having it depend on this dummy target
FORCE:

.PHONY: help
help: ## Show a list of all targets
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1:\3/p' \
	| column -t -s ":"

.PHONY: clean
clean: ## Clean up the build directory
	rm -rf ../build

.PHONY: build-and-deploy
build-and-deploy: ## Build the disaster-recovery package and immediately deploy it
	../build/$(ZARF_BIN) package create --confirm
	../build/$(ZARF_BIN) package deploy zarf-package-alt-backupstoragelocation-amd64.tar.zst --confirm

#.PHONY: destroy
#destroy: | ../build/helm ## Destroy the disaster-recovery package using helm
#

.PHONY: k9s
k9s: ## Pull up K9s
	../build/$(ZARF_BIN) tools k9s

build:
	@mkdir -p ../build

../build/helm: | build ## Download the Helm binary
	curl -sL https://get.helm.sh/helm-v3.10.1-linux-amd64.tar.gz -o ../build/helm.tar.gz
	tar -xzvf ../build/helm.tar.gz -C ../build --strip-components 1 linux-amd64/helm
	rm ../build/helm.tar.gz
	chmod +x ../build/helm
