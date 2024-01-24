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

4. Wait for everything to come up. Use `./zarf tools monitor` to monitor using the [K9s](https://github.com/derailed/k9s) tool
