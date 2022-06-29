# Fork the repo, customize, and build the packages

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

1. Clone your new repo and add this repo as an "upstream" remote, so you can pull down updates later

    ```shell
    git clone https://gitsite.com/yourusername/new-repository.git
    cd new-repository
    # If you forked on GitHub they already did this for you
    git remote add upstream https://github.com/defenseunicorns/zarf-package-software-factory.git
    ```

1. Customize `zarf.yaml` -- Change the repo URL in the "setup" component from `https://github.com/defenseunicorns/zarf-package-software-factory.git` to the repo URL of your config repo that you created by forking the upstream

1. Customize `day2/zarf.yaml` -- Change the repo URL from `https://github.com/defenseunicorns/zarf-package-software-factory.git` to the repo URL of your config repo that you created by forking the upstream

1. Customize `manifests/setup.yaml` -- Change the repo URL `https://github.com/defenseunicorns/zarf-package-software-factory.git` to the repo URL of your config repo that you created by forking the upstream. Also change the branch name specified from `not-the-real-branch-name` to a real branch name. Our recommendation is to create a new branch off of `main` called `main_airgap` and use that. That gives you the ability to easily pull in upstream changes to `main` and then do pull requests from `main` to `main_airgap` as you are able to.

1. Customize `kustomizations/bigbang/environment-bb/values.yaml` -- Replace `bigbang.dev` with your real domain, and change the TLS key and cert to your own key and cert, then SOPS encrypt the file. Click [HERE](sops.md) for instructions on how to set up SOPS encryption.

    ```shell
    sops -e -i kustomizations/bigbang/environment-bb/values.yaml
    ```

1. Customize `kustomizations/softwarefactoryaddons/jenkins/environment-bb-values.yaml` -- Replace `bigbang.dev` with your real domain. Do a find and replace on the whole file, it appears in multiple places. Later on in the [SSO](sso.md) step you'll also update the `clientID` and `clientSecret` parameters but we can't do that until after GitLab is deployed. Encrypt the file with SOPS if you want at this point, though the only things in the file that are likely to be considered secrets are the client ID and client secret, which won't have been added yet.

    ```shell
    sops -e -i kustomizations/softwarefactoryaddons/jenkins/environment-bb-values.yaml
    ```

1. Customize `kustomizations/softwarefactoryaddons/base/virtualservice.yaml` -- Replace `bigbang.dev` with your real domain.

1. Customize `kustomizations/softwarefactoryaddons/jira/values.yaml` -- Replace `bigbang.dev` with your real domain

1. Customize `kustomizations/softwarefactoryaddons/confluence/environment-bb-values.yaml` -- Replace `bigbang.dev` with your real domain

1. Customize `kustomizations/softwarefactoryaddons/confluence/environment-bb-values.yaml` -- Replace `bigbang.dev` with your real domain

1.  Commit the changes to the repo

    ```shell
    git add .
    git commit -m "Add environment-specific configuration"
    git push
    ```

1. Build the packages

    ```shell
    make all
    ```

Now that the necessary packages are created, it is time to [Deploy](deploy.md).
