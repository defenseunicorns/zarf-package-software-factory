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

# Silent mode by default. Run `make VERBOSE=1` to turn off silent mode.
ifndef VERBOSE
.SILENT:
endif

# Idiomatic way to force a target to always run, by having it depend on this dummy target
FORCE:

.PHONY: help
help: ## Show a list of all targets
	grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1:\3/p' \
	| column -t -s ":"

.PHONY: build-and-deploy
build-and-deploy: ## Build the day2 package and immediately deploy it
	../build/$(ZARF_BIN) package create --confirm
	../build/$(ZARF_BIN) package deploy zarf-package-day-two-update-amd64.tar.zst --confirm

.PHONY: k9s
k9s: ## Pull up K9s
	../build/$(ZARF_BIN) tools k9s

.PHONY: nuke-and-rebuild
nuke-and-rebuild: ## Nuke the cluster and rebuild it
	cd ~/app/build && ./zarf destroy --confirm || true
	cd ~/app && make build/zarf-package-software-factory-amd64.tar.zst
	cd ~/app/build && ./zarf init --components k3s,git-server --confirm
	cd ~/app/build && ./zarf package deploy zarf-package-flux-amd64.tar.zst --confirm
	gpg --list-secret-keys user@example.com || gpg --batch --passphrase '' --quick-gen-key user@example.com default default
	gpg --export-secret-keys --armor user@example.com | kubectl create secret generic sops-gpg -n flux-system --from-file=sops.asc=/dev/stdin
	cd ~/app/build && ./zarf package deploy zarf-package-software-factory-amd64.tar.zst --components optional-tools-linux-amd64 --confirm
	kubectl patch gitrepositories.source.toolkit.fluxcd.io -n flux-system zarf-package-software-factory --type=json -p '[{"op": "replace", "path": "/spec/ref/branch", "value": "$(shell git rev-parse --abbrev-ref HEAD)"}]'
