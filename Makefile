# DI2ME package version
PACKAGE_VERSION := 0.1.0
# The version of Big Bang to use. If you change this you need to also do a couple of other things:
#    1. Run `make vendor-big-bang-base` and commit any changes to the repo.
#    2. Additionally update the following files to use the new version of Big Bang:
#        - zarf.yaml
#        - flux/zarf.yaml
BIGBANG_VERSION := 2.2.0

# The version of Zarf to use. To keep this repo as portable as possible the Zarf binary will be downloaded and added to
# the build folder.
ZARF_VERSION := v0.27.0

# The version of the build harness container to use
BUILD_HARNESS_REPO := ghcr.io/defenseunicorns/not-a-build-harness/not-a-build-harness
BUILD_HARNESS_VERSION := 0.0.25

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

# Silent mode by default. Run `make VERBOSE=1` to turn off silent mode.
ifndef VERBOSE
.SILENT:
endif

# Optionally add the "-it" flag for docker run commands if the env var "CI" is not set (meaning we are on a local machine and not in github actions)
TTY_ARG :=
ifndef CI
	TTY_ARG := -it
endif

.DEFAULT_GOAL := help

# Idiomatic way to force a target to always run, by having it depend on this dummy target
FORCE:

.PHONY: help
help: ## Show a list of all targets
	grep -E '^\S*:.*##.*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1:\3/p' \
	| column -t -s ":"

.PHONY: docker-save-build-harness
docker-save-build-harness: ## Pulls the build harness docker image and saves it to a tarball
	mkdir -p .cache/docker
	docker pull $(BUILD_HARNESS_REPO):$(BUILD_HARNESS_VERSION)
	docker save -o .cache/docker/build-harness.tar $(BUILD_HARNESS_REPO):$(BUILD_HARNESS_VERSION)

.PHONY: docker-load-build-harness
docker-load-build-harness: ## Loads the saved build harness docker image
	docker load -i .cache/docker/build-harness.tar

.PHONY: run-pre-commit-hooks
run-pre-commit-hooks: ## Run all pre-commit hooks. Returns nonzero exit code if any hooks fail. Uses Docker for maximum compatibility
	mkdir -p .cache/pre-commit
	docker run --rm -v "${PWD}:/app" --workdir "/app" -e "PRE_COMMIT_HOME=/app/.cache/pre-commit" $(BUILD_HARNESS_REPO):$(BUILD_HARNESS_VERSION) bash -c 'git config --global --add safe.directory /app && asdf install && pre-commit run -a'

.PHONY: fix-cache-permissions
fix-cache-permissions: ## Fixes the permissions on the pre-commit cache
	docker run --rm -v "${PWD}:/app" --workdir "/app" -e "PRE_COMMIT_HOME=/app/.cache/pre-commit" $(BUILD_HARNESS_REPO):$(BUILD_HARNESS_VERSION) chmod -R a+rx .cache

# TODO: Figure out how to make it log to the console in real time so the user isn't sitting there wondering if it is working or not.
.PHONY: test
test: ## Run all automated tests. Requires access to an AWS account. Costs money. Requires env vars "REPO_URL", "GIT_BRANCH", "REGISTRY1_USERNAME", "REGISTRY1_PASSWORD", and standard AWS env vars.
	mkdir -p .cache/go
	mkdir -p .cache/go-build
	echo "Running automated tests. This will take several minutes. At times it does not log anything to the console. If you interrupt the test run you will need to log into AWS console and manually delete any orphaned infrastructure."
	docker run $(TTY_ARG) --rm -v "${PWD}:/app" -v "${PWD}/.cache/go:/root/go" -v "${PWD}/.cache/go-build:/root/.cache/go-build" --workdir "/app/test/e2e" -e GOPATH=/root/go -e GOCACHE=/root/.cache/go-build -e REPO_URL -e GIT_BRANCH -e REGISTRY1_USERNAME -e REGISTRY1_PASSWORD -e AWS_REGION -e AWS_DEFAULT_REGION -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_SESSION_TOKEN -e AWS_SECURITY_TOKEN -e AWS_SESSION_EXPIRATION -e SKIP_SETUP -e SKIP_TEST -e SKIP_TEARDOWN $(BUILD_HARNESS_REPO):$(BUILD_HARNESS_VERSION) bash -c 'asdf install && go test -v -timeout 2h -p 1 ./...'

.PHONY: test-ssh
test-ssh: ## Run this if you set SKIP_TEARDOWN=1 and want to SSH into the still-running test server. Don't forget to unset SKIP_TEARDOWN when you're done
	cd test/tf/public-ec2-instance && terraform init
	cd test/tf/public-ec2-instance/.test-data && cat Ec2KeyPair.json | jq -r .PrivateKey > privatekey.pem && chmod 600 privatekey.pem
	cd test/tf/public-ec2-instance && ssh -i .test-data/privatekey.pem ubuntu@$$(terraform output public_instance_ip | tr -d '"')

create-cluster:
	kind create cluster --name di2me --config test/kind-config/noCNI.yaml
	echo
	echo "Waiting for cluster to be ready..."
	kubectl wait --for=condition=Ready pods --all --all-namespaces 2>&1 >/dev/null
	echo
	echo "Installing Calico..."
	kubectl apply --wait=true -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.1/manifests/calico.yaml 2>&1 >/dev/null
	echo "Waiting for Calico to be ready..."
	kubectl rollout status deployment/calico-kube-controllers -n kube-system --watch --timeout=90s 2>&1 >/dev/null
	kubectl rollout status daemonset/calico-node -n kube-system --watch --timeout=90s 2>&1 >/dev/null
	kubectl wait --for=condition=Ready pods --all --all-namespaces 2>&1 >/dev/null
	echo
	test/metallb/install.sh
	kubectl wait --for=condition=Ready pods --all --all-namespaces 2>&1 >/dev/null
	echo
	echo "Cluster is ready!"

destroy-cluster:
	kind delete cluster --name di2me

day2-create:
	cd day2 && ../build/$(ZARF_BIN) package create --skip-sbom --confirm --set DI2ME_REPO="https://github.com/defenseunicorns/zarf-package-software-factory.git@$$(git show-ref --heads --tags | grep /$$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match)$$ | cut -d ' ' -f2)"

day2-deploy:
	build/$(ZARF_BIN) package deploy --components=neuvector-cve-update --confirm day2/zarf-package-software-factory-amd64.tar.zst

default-build: ## All in one make target for the default di2me repo (only x86) - uses the current branch/tag of the repo
	make build
	make build/zarf
	make build/zarf-init.sha256
	make build/zarf-package-flux-amd64.tar.zst
	make build/zarf-package-software-factory-amd64.tar.zst DI2ME_REPO="https://github.com/defenseunicorns/zarf-package-software-factory.git@$$(git show-ref --heads --tags | grep /$$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match)$$ | cut -d ' ' -f2)"

.PHONY: deploy-local
deploy-local: ## Deploy created zarf package to local cluster
	cat test/e2e/zarf-config.toml | grep -v progress > build/zarf-config.toml
	cd build && ./zarf init --components git-server --confirm
	cd build && ./zarf package deploy zarf-package-flux-amd64.tar.zst --confirm
	gpg --list-secret-keys user@example.com || gpg --batch --passphrase '' --quick-gen-key user@example.com default default
	gpg --export-secret-keys --armor user@example.com | kubectl create secret generic sops-gpg -n flux-system --from-file=sops.asc=/dev/stdin
	cd build && ./zarf package deploy zarf-package-software-factory-amd64-$(PACKAGE_VERSION).tar.zst --confirm
	kubectl patch gitrepositories.source.toolkit.fluxcd.io -n flux-system zarf-package-software-factory --type=json -p '[{"op": "replace", "path": "/spec/ref/branch", "value": "$(shell git rev-parse --abbrev-ref HEAD)"}]'
	timeout 2400 bash -c "while ! kubectl get cronjob gitlab-toolbox-backup -n gitlab; do sleep 5; done"
	kubectl create job -n gitlab --from=cronjob/gitlab-toolbox-backup gitlab-toolbox-backup-manual

.PHONY: clean
clean: ## Clean up build files
	rm -rf ./build

.PHONY: all
all: | build/zarf build/zarf-mac-intel build/zarf-init.sha256 build/zarf-package-flux-amd64.tar.zst build/zarf-package-software-factory-amd64.tar.zst ## Make everything. Will skip downloading/generating dependencies if they already exist.

.PHONY: vendor-big-bang-base
vendor-big-bang-base: ## Vendor the BigBang base kustomization, since Flux doesn't support private bases. This only needs to be run if you change the version of Big Bang used. Don't forget to commit the changes to the repo.
	rm -rf kustomizations/bigbang/vendor/bigbang && \
	mkdir -p kustomizations/bigbang/vendor && \
	cd kustomizations/bigbang/vendor && \
	git init bigbang && \
	cd bigbang && \
	git remote add -f origin https://repo1.dso.mil/big-bang/bigbang.git && \
	git config core.sparseCheckout true && \
	echo "base/" > .git/info/sparse-checkout && \
	git checkout tags/$(BIGBANG_VERSION) -b tagbranch && \
	rm -rf .git && \
	rm -rf base/flux

build:
	mkdir -p build

build/zarf: | build ## Download the Linux flavor of Zarf to the build dir
	echo "Downloading zarf"
	curl -sL https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf_$(ZARF_VERSION)_Linux_amd64 -o build/zarf
	chmod +x build/zarf

build/zarf-mac-intel: | build ## Download the Mac (Intel) flavor of Zarf to the build dir
	echo "Downloading zarf-mac-intel"
	curl -sL https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf_$(ZARF_VERSION)_Darwin_amd64 -o build/zarf-mac-intel
	chmod +x build/zarf-mac-intel

build/zarf-init.sha256: | build ## Download the init package and create a small file with the sha256sum of the package so the Makefile can check whether it needs to be updated
	echo "Downloading zarf-init-amd64-$(ZARF_VERSION).tar.zst"
	curl -sL https://github.com/defenseunicorns/zarf/releases/download/$(ZARF_VERSION)/zarf-init-amd64-$(ZARF_VERSION).tar.zst -o build/zarf-init-amd64-$(ZARF_VERSION).tar.zst
	echo "Creating shasum of the init package"
	shasum -a 256 build/zarf-init-amd64-$(ZARF_VERSION).tar.zst | awk '{print $$1}' > build/zarf-init.sha256

build/zarf-package-flux-amd64.tar.zst: | build/$(ZARF_BIN) ## Build the Flux package
	cd flux && ../build/$(ZARF_BIN) package create --skip-sbom --confirm
	mv flux/zarf-package-flux-amd64.tar.zst build/zarf-package-flux-amd64.tar.zst

build/zarf-package-software-factory-amd64.tar.zst: FORCE | build/$(ZARF_BIN) ## Build the Software Factory package
	echo "Creating the deploy package"
	build/$(ZARF_BIN) package create --skip-sbom --confirm --set PACKAGE_VERSION=$(PACKAGE_VERSION) --set DI2ME_REPO=$(DI2ME_REPO)
	mv zarf-package-software-factory-amd64-$(PACKAGE_VERSION).tar.zst build/zarf-package-software-factory-amd64-$(PACKAGE_VERSION).tar.zst
