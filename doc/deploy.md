# Deploy the packages

This guide assumes you have already [Forked, customized, and built the packages](fork-and-build.md). If you haven't, please do that first.

Depending on where you want to run the package you just created, there are a few different paths

## Airgap

1. Burn to removable media
    - `build/zarf`
    - `zarf-package-k3s-amd64.tar.zst`
    - `zarf-package-k3s-images-amd64.tar.zst`
    - `build/zarf-init-amd64.tar.zst`
    - `build/zarf-package-flux-amd64.tar.zst`
    - `build/zarf-package-software-factory-amd64.tar.zst`
    - `secret-sops-gpg.yaml` (See [SOPS configuration](sops.md))

2. Use [Sneakernet](https://en.wikipedia.org/wiki/Sneakernet) or whatever other method you want to get it where it needs to go

3. Deploy

   ```shell
   # Assuming you want to use the built-in single-node K3s cluster. If you don't, skip the "k3s" and "k3s-images" packages
   ./zarf package deploy zarf-package-k3s-amd64.tar.zst --confirm
   ./zarf init --components git-server --confirm
   ./zarf package deploy zarf-package-k3s-images-amd64.tar.zst --confirm
   ./zarf package deploy zarf-package-flux-amd64.tar.zst --confirm
   kubectl apply -f secret-sops-gpg.yaml
   ./zarf package deploy zarf-package-software-factory-amd64.tar.zst --confirm
   ```

4. Wait for everything to come up. Use `./zarf tools k9s` to monitor using the [K9s](https://github.com/derailed/k9s) tool

## Vagrant

1. Update the package to disable enough services so that it will fit on your computer. Big Bang deployments can be disabled by changing `enabled: true` to `enabled: false` in `kustomizations/bigbang/values.yaml`. Other addons like Jira, Confluence, and Jenkins can be disabled by modifying the Kustomization in the `kustomizations/softwarefactoryaddons` directory.

1. Rebuild the package -- `make all`

1. Spin up the VM with `make vm-init`. You will automatically be dropped into a root shell inside the VM with the `build` folder mounted as the current working directory. If you need to leave the VM type `exit`. To get back in type `make vm-ssh`.

1. Create `secret-sops-gpg.yaml` (See [SOPS configuration](sops.md))

1. Deploy

    ```shell
    ./zarf package deploy zarf-package-k3s-amd64.tar.zst --confirm
    ./zarf init --components git-server --confirm
    ./zarf package deploy zarf-package-k3s-images-amd64.tar.zst --confirm
    ./zarf package deploy zarf-package-flux-amd64.tar.zst --confirm
    kubectl apply -f secret-sops-gpg.yaml
    ./zarf package deploy zarf-package-software-factory-amd64.tar.zst --confirm
    ```

1. Wait for everything to come up. Use `./zarf tools k9s` to monitor using the [K9s](https://github.com/derailed/k9s) tool

1. When you're done, `exit`, then bring everything down with `make vm-destroy`

## Cloud

1. Copy files to the machine that will be running `zarf`. Could be your local computer, or could be a Bastion Host, or something else entirely. Zarf doesn't care, as long as it has a connection to the cluster.
    - `build/zarf`
    - `build/zarf-init-amd64.tar.zst`
    - `build/zarf-package-flux-amd64.tar.zst`
    - `build/zarf-package-software-factory-amd64.tar.zst`
    - `secret-sops-gpg.yaml` (See [SOPS configuration](sops.md))

2. Have a Kubernetes cluster ready that you'll be deploying to. Have your KUBECONFIG be configured such that running `kubectl get nodes` will connect to the right cluster. OR use the built-in K3s that Zarf comes bundled with.

3. Deploy

    ```shell
    ./zarf init --components gitops-service --confirm
    ./zarf package deploy zarf-package-flux-amd64.tar.zst --confirm
    kubectl apply -f secret-sops-gpg.yaml
    ./zarf package deploy zarf-package-software-factory-amd64.tar.zst --confirm
    ```

4. Wait for everything to come up. Use `./zarf tools k9s` to monitor using the [K9s](https://github.com/derailed/k9s) tool.

Now that everything is deployed, it is time to [Configure Single Sign-On](sso.md)
