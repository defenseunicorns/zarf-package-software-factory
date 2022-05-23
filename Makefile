# The version of Big Bang to use. If you change this you need to also do a couple of other things:
#    1. Run `make vendor-big-bang-base` and commit any changes to the repo.
#    2. Additionally update the following files to use the new version of Big Bang:
#        - zarf.yaml
BIGBANG_VERSION := 1.28.0

# The version of Zarf to use. To keep this repo as portable as possible the Zarf binary will be downloaded and added to
# the build folder.
ZARF_VERSION := v0.17.0

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

.PHONY: build-harness-shell
build-harness-shell: ## Open a shell in the build harness container with the project mounted
	docker run -it --rm -v "${PWD}:/app" --workdir "/app" -e "PRE_COMMIT_HOME=/app/.cache/pre-commit" ghcr.io/defenseunicorns/zarf-package-software-factory/build-harness:0.0.2 bash

.PHONY: run-pre-commit-hooks
run-pre-commit-hooks: ## Run all pre-commit hooks. Returns nonzero exit code if any hooks fail. Recommend running with `make build-harness-shell` followed by `make run-pre-commit-hooks`
	echo "hello world"

.PHONY: vm-init
vm-init: vm-destroy ## Stripped-down vagrant box to reduce friction for basic user testing. Note the need to perform disk resizing for some examples
	@VAGRANT_EXPERIMENTAL="disks" vagrant up --no-color
	@vagrant ssh

.PHONY: vm-ssh
vm-ssh: ## SSH into the Vagrant VM
	vagrant ssh

.PHONY: vm-destroy
vm-destroy: ## Destroy the Vagrant VM
	@vagrant destroy -f

.PHONY: clean
clean: ## Clean up build files
	@rm -rf ./build

.PHONY: all
all: | build/zarf build/zarf-mac-intel build/zarf-init-amd64.tar.zst build/zarf-package-flux-amd64.tar.zst build/zarf-package-software-factory-amd64.tar.zst ## Make everything. Will skip downloading/generating dependencies if they already exist.

.PHONY: vendor-big-bang-base
vendor-big-bang-base: ## Vendor the BigBang base kustomization, since Flux doesn't support private bases. This only needs to be run if you change the version of Big Bang used. Don't forget to commit the changes to the repo.
	@rm -rf kustomizations/bigbang/vendor/bigbang && \
	mkdir -p kustomizations/bigbang/vendor && \
	cd kustomizations/bigbang/vendor && \
	git init bigbang && \
	cd bigbang && \
	git remote add -f origin https://repo1.dso.mil/platform-one/big-bang/bigbang.git && \
	git config core.sparseCheckout true && \
	echo "base/" > .git/info/sparse-checkout && \
	git checkout tags/$(BIGBANG_VERSION) -b tagbranch && \
	rm -rf .git && \
	rm -rf base/flux

build:
	@mkdir -p build

build/zarf: | build
	@echo "Downloading zarf"
	@wget https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf -O build/zarf
	@chmod +x build/zarf

build/zarf-mac-intel: | build
	@echo "Downloading zarf-mac-intel"
	@wget https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf-mac-intel -O build/zarf-mac-intel
	@chmod +x build/zarf-mac-intel

build/zarf-init-amd64.tar.zst: | build
	@echo "Downloading zarf-init-amd64.tar.zst"
	@wget https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf-init-amd64.tar.zst -O build/zarf-init-amd64.tar.zst

build/zarf-package-flux-amd64.tar.zst: | build/$(ZARF_BIN)
	@rm -rf ./tmp
	@mkdir -p ./tmp
	@git clone -b $(ZARF_VERSION) --depth 1 https://github.com/defenseunicorns/zarf.git tmp/zarf
	@cd tmp/zarf/packages/flux-iron-bank && ../../../../build/$(ZARF_BIN) package create --confirm
	@mv tmp/zarf/packages/flux-iron-bank/zarf-package-flux-amd64.tar.zst build/zarf-package-flux-amd64.tar.zst
	@rm -rf ./tmp

build/zarf-package-software-factory-amd64.tar.zst: FORCE | build/$(ZARF_BIN)
	@echo "Creating the deploy package"
	@build/$(ZARF_BIN) package create --confirm
	@mv zarf-package-software-factory-amd64.tar.zst build/zarf-package-software-factory-amd64.tar.zst
