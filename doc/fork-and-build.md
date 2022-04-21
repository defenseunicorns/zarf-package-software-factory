# Fork the repo and build the packages

Since you will need to make environment-specific changes to the system's configuration, you should fork this repository, and update the package configuration to look at your fork. Here's the steps to take:

1. Fork the repo. On GitHub that can be done by clicking the "Fork" button in the top right of the page. For any other Git system you'll want to create a bare clone and do a mirror push. Like this:

   ```shell
   # Assuming you have created a brand new completely empty repo located at https://gitsite.com/yourusername/new-repository.git
   git clone --bare https://github.com/defenseunicorns/zarf-package-software-factory.git
   cd zarf-package-software-factory.git
   git push --mirror https://gitsite.com/yourusername/new-repository.git
   cd ..
   rm -rf ./zarf-package-software-factory.git
   ```

> Note: If you want to make the repo private don't use the "Fork" feature on GitHub, since forks can't be made private unless you first submit a support request to have them detach the fork. Note that if you are using SOPS encryption of your secrets (highly recommended) then it is okay for your config repo to be public since the files that contain secrets will be committed to the repository encrypted, and decrypted inside the cluster using Flux.

2. Clone your new repo and add this repo as an "upstream" remote, so you can pull down updates later

   ```shell
   git clone https://gitsite.com/yourusername/new-repository.git
   cd new-repository
   # If you forked on GitHub they already did this for you
   git remote add upstream https://github.com/defenseunicorns/zarf-package-software-factory.git
   ```

3. Customize `zarf.yaml` -- Change the first entry in the "big-bang" component from `https://github.com/defenseunicorns/zarf-package-software-factory.git` to the git URL of your config repo that you created by forking the upstream

4. Customize `day2/zarf.yaml` -- Change the referenced git URL from `https://github.com/defenseunicorns/zarf-package-software-factory.git` to the git URL of your config repo that you created by forking the upstream

5. Customize `manifests/big-bang.yaml` -- Change the url `http://zarf-gitea-http.zarf.svc.cluster.local:3000/zarf-git-user/mirror__github.com__defenseunicorns__zarf-package-software-factory.git` to the "Zarf-ified" version of your config repo that you created by forking the upstream. The easiest way to do that is to change it to the regular URL, then run this command on that file:

   ```shell
   zarf prepare patch-git http://zarf-gitea-http.zarf.svc.cluster.local:3000 manifests/big-bang.yaml
   ```

> Note: If you need to install Zarf, you can run either `make build/zarf-mac-intel` or `make build/zarf` (depending on what OS distro you are using). Zarf will be installed in the `build` folder in this repo.

6. Customize `manifests/softwarefactoryaddons.yaml` -- Change the url `http://zarf-gitea-http.zarf.svc.cluster.local:3000/zarf-git-user/mirror__github.com__defenseunicorns__zarf-package-software-factory.git` to the "Zarf-ified" version of your config repo that you created by forking the upstream. The easiest way to do that is to change it to the regular URL, then run this command on that file:

```shell
zarf prepare patch-git http://zarf-gitea-http.zarf.svc.cluster.local:3000 manifests/softwarefactoryaddons.yaml
``

7. Customize `kustomizations/bigbang/environment-bb/values.yaml` -- Replace `bigbang.dev` with your real domain, and change the TLS key and cert to your own key and cert, then SOPS encrypt the file. Click [HERE](sops.md) for instructions on how to set up SOPS encryption.

7. Generate `zarf.yaml`, `manifests/big-bang.yaml`, and `manifests/softwarefactoryaddons.yaml` from the provided templates:

   ```shell
   # First download Zarf. Assuming you are on MacOS, otherwise on Linux switch the target to `build/zarf` and the calls to `build/zarf` instead of `build/zarf-mac-intel`
   make build/zarf-mac-intel
   # Then use the provided template files to generate the real one
   export CONFIG_REPO_URL="https://gitsite.com/yourusername/new-repository.git"
   envsubst '$CONFIG_REPO_URL' < zarf.tmpl.yaml > zarf.yaml
   envsubst '$CONFIG_REPO_URL' < day2/zarf.tmpl.yaml > day2/zarf.yaml
   envsubst '$CONFIG_REPO_URL' < manifests/big-bang.tmpl.yaml > manifests/big-bang.yaml
   envsubst '$CONFIG_REPO_URL' < manifests/softwarefactoryaddons.tmpl.yaml > manifests/softwarefactoryaddons.yaml
   # These ones will require you to confirm that you want to perform this action by typing "y"
   build/zarf-mac-intel prepare patch-git "http://zarf-gitea-http.zarf.svc.cluster.local:3000" manifests/big-bang.yaml
   build/zarf-mac-intel prepare patch-git "http://zarf-gitea-http.zarf.svc.cluster.local:3000" manifests/softwarefactoryaddons.yaml
   ```

8. Modify the package to use your DNS domain and your TLS certificate and key:

   **TODO:** write this stuff. It will need to use SOPS since it will have a real TLS key, which isn't so bad on EKS but will be a challenge to get working in the airgap. For now the bigbang.dev domain and TLS cert/key can be used for dev/test.

   > IMPORTANT NOTE: _**YOUR TLS CERT KEY MUST BE TREATED AS A SECRET**_. Never commit the actual secret to a git repository.

9.  Commit the changes to the repo

   ```shell
   git add .
   git commit -m "Add environment-specific configuration"
   git push
   ```

11. Build the packages

   ```shell
   make all
   ```

Now that the necessary packages are created, it is time to [Initialize the cluster](initialize.md).
