.DEFAULT_GOAL := help


.PHONY: help
help: ## Show a list of all targets
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1:\3/p' \
	| column -t -s ":"

.PHONY: vm-init
vm-init: vm-destroy ## Stripped-down vagrant box to reduce friction for basic user testing. Note the need to perform disk resizing for some examples
	@VAGRANT_EXPERIMENTAL="disks" vagrant up --no-color
	@vagrant ssh

.PHONY: vm-destroy
vm-destroy: ## Cleanup plz
	@vagrant destroy -f

clean: ## Clean up build files
	@rm -rf ./build

.PHONY: build
build: ## Create the software factory deploy package
	@mkdir -p ./build
	@zarf package create --confirm && mv zarf-package-* ./build/

.PHONY: ssh
ssh: ## SSH into the Vagrant VM
	vagrant ssh
