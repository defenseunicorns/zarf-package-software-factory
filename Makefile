# The version of Big Bang to use. If you change this you need to also do a couple of other things:
#    1. Run `make vendor-big-bang-base` and commit any changes to the repo.
#    2. Additionally update the following files to use the new version of Big Bang:
#        - zarf.yaml
BIGBANG_VERSION := 1.28.0

# The version of Zarf to use. To keep this repo as portable as possible the Zarf binary will be downloaded and added to
# the build folder.
ZARF_VERSION := v0.17.0

# The version of the build harness container to use
BUILD_HARNESS_VERSION := 0.0.7

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
	@grep -E '^\S*:.*##.*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1:\3/p' \
	| column -t -s ":"

.PHONY: docker-save-build-harness
docker-save-build-harness: ## Pulls the build harness docker image and saves it to a tarball
	mkdir -p .cache/docker
	docker pull ghcr.io/defenseunicorns/zarf-package-software-factory/build-harness:$(BUILD_HARNESS_VERSION)
	docker save -o .cache/docker/build-harness.tar ghcr.io/defenseunicorns/zarf-package-software-factory/build-harness:$(BUILD_HARNESS_VERSION)

.PHONY: docker-load-build-harness
docker-load-build-harness: ## Loads the saved build harness docker image
	docker load -i .cache/docker/build-harness.tar

.PHONY: run-pre-commit-hooks
run-pre-commit-hooks: ## Run all pre-commit hooks. Returns nonzero exit code if any hooks fail. Uses Docker for maximum compatibility
	mkdir -p .cache/pre-commit
	docker run --rm -v "${PWD}:/app" --workdir "/app" -e "PRE_COMMIT_HOME=/app/.cache/pre-commit" ghcr.io/defenseunicorns/zarf-package-software-factory/build-harness:$(BUILD_HARNESS_VERSION) pre-commit run -a

.PHONY: fix-cache-permissions
fix-cache-permissions: ## Fixes the permissions on the pre-commit cache
	docker run --rm -v "${PWD}:/app" --workdir "/app" -e "PRE_COMMIT_HOME=/app/.cache/pre-commit" ghcr.io/defenseunicorns/zarf-package-software-factory/build-harness:$(BUILD_HARNESS_VERSION) chmod -R a+rx .cache

# TODO: Figure out how to make it log to the console in real time so the user isn't sitting there wondering if it is working or not.
.PHONY: test
test: ## Run all automated tests. Requires access to an AWS account. Costs money.
	@mkdir -p .cache/go
	@mkdir -p .cache/go-build
	@echo "Running automated tests. This will take several minutes. At times it does not log anything to the console. If you interrupt the test run you will need to log into AWS console and manually delete any orphaned infrastructure."
	@docker run --rm -v "${PWD}:/app" -v "${PWD}/.cache/go:/root/go" -v "${PWD}/.cache/go-build:/root/.cache/go-build" --workdir "/app/test/e2e" -e GOPATH=/root/go -e GOCACHE=/root/.cache/go-build -e REPO_URL -e GIT_BRANCH -e REGISTRY1_USERNAME -e REGISTRY1_PASSWORD -e AWS_REGION -e AWS_DEFAULT_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY ghcr.io/defenseunicorns/zarf-package-software-factory/build-harness:$(BUILD_HARNESS_VERSION) go test -v -timeout 1h -p 1 ./...

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

build/zarf: | build ## Download the Linux flavor of Zarf to the build dir
	@echo "Downloading zarf"
	@wget https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf -O build/zarf
	@chmod +x build/zarf

build/zarf-mac-intel: | build ## Download the Mac (Intel) flavor of Zarf to the build dir
	@echo "Downloading zarf-mac-intel"
	@wget https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf-mac-intel -O build/zarf-mac-intel
	@chmod +x build/zarf-mac-intel

build/zarf-init-amd64.tar.zst: | build ## Download the init package
	@echo "Downloading zarf-init-amd64.tar.zst"
	@wget https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf-init-amd64.tar.zst -O build/zarf-init-amd64.tar.zst

build/zarf-package-flux-amd64.tar.zst: | build/$(ZARF_BIN) ## Build the Flux package
	@rm -rf ./tmp
	@mkdir -p ./tmp
	@git clone -b $(ZARF_VERSION) --depth 1 https://github.com/defenseunicorns/zarf.git tmp/zarf
	@cd tmp/zarf/packages/flux-iron-bank && ../../../../build/$(ZARF_BIN) package create --confirm
	@mv tmp/zarf/packages/flux-iron-bank/zarf-package-flux-amd64.tar.zst build/zarf-package-flux-amd64.tar.zst
	@rm -rf ./tmp

build/zarf-package-software-factory-amd64.tar.zst: FORCE | build/$(ZARF_BIN) ## Build the Software Factory package
	@echo "Creating the deploy package"
	@build/$(ZARF_BIN) package create --confirm
	@mv zarf-package-software-factory-amd64.tar.zst build/zarf-package-software-factory-amd64.tar.zst
