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

> Note: If you want to make the repo private don't use the "Fork" feature on GitHub, since GitHub won't let you change the visibility of forks to private.

2. Clone your new repo and add this repo as an "upstream" remote, so you can pull down updates later

   ```shell
   git clone https://gitsite.com/yourusername/new-repository.git
   cd new-repository
   # If you forked on GitHub they already did this for you
   git remote add upstream https://github.com/defenseunicorns/zarf-package-software-factory.git
   ```

3. Generate `zarf.yaml`, `manifests/big-bang.yaml`, and `manifests/softwarefactoryaddons.yaml` from the provided templates:

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

4. Modify the package to use your DNS domain and your TLS certificate and key:

   **TODO:** write this stuff. It will need to use SOPS since it will have a real TLS key, which isn't so bad on EKS but will be a challenge to get working in the airgap. For now the bigbang.dev domain and TLS cert/key can be used for dev/test.

> IMPORTANT NOTE: _**YOUR TLS CERT KEY MUST BE TREATED AS A SECRET**_. Never commit the actual secret to a git repository.

5. Commit the changes to the repo

   ```shell
   git add .
   git commit -m "Add environment-specific configuration"
   git push
   ```

6. Build the packages

   ```shell
   make all
   ```

Now that the necessary packages are created, it is time to [Deploy the packages](deploy.md)