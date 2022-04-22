# Configure Single Sign-On

NOTE: `bigbang.dev` is the default domain. If you are using a different domain, substitute `bigbang.dev` for your domain in all URLs

## Configure GitLab

1. Retrieve the initial root password for GitLab:

    ```shell
    kubectl get secret gitlab-gitlab-initial-root-password -n gitlab -o jsonpath='{.data.password}' | base64 --decode
    ```

3. Navigate to [https://gitlab.bigbang.dev](https://gitlab.bigbang.dev)

4. Log in using username `root` with the password retrieved from the previous step

5. Navigate to [https://gitlab.bigbang.dev/-/profile/password/edit](https://gitlab.bigbang.dev/-/profile/password/edit) and change the root password. Save the new password as you will need it in disaster recovery scenarios.

6. [OPTIONAL] Disable Sign-up in the Sign-up restrictions section on [https://gitlab.bigbang.dev/admin/application_settings/general](https://gitlab.bigbang.dev/admin/application_settings/general). If you disable it you will need to manually create all new users. It may be advantageous to leave it on, since you can require admin approval for new sign-ups. Click "Save Changes" at the bottom of the section.

7. Enable "Enforce two-factor authentication" in the Sign-in restrictions section on [https://gitlab.bigbang.dev/admin/application_settings/general](https://gitlab.bigbang.dev/admin/application_settings/general). Click "Save Changes" at the bottom of the section.

8. Configure two-factor authentication on the root account. Make sure this gets done right away. If you wait past the grace period the root account will be locked out.

### Configure Jenkins

**Warnings and notes**

- :warning: WARNING: This section involves committing and pushing a secret to your config repo. The file needs to be encrypted with SOPS before you push it to your repo. See the [SOPS guide](sops.md) for instructions.

- NOTE: Configuring Jenkins to use GitLab as the SSO provider will not work when using the `bigbang.dev` domain. During the OAuth ping-pong auth flow Jenkins tries to run a POST http request to get a login token, but it isn't able to hit `https://gitlab.bigbang.dev` since that resolves to `127.0.0.1`. This may be fixable through the use of an Istio [ServiceEntry](https://istio.io/latest/docs/reference/config/networking/service-entry/). More investigation needed.

### Instructions

1. Navigate to [https://gitlab.bigbang.dev/admin/applications/new](https://gitlab.bigbang.dev/admin/applications/new) and create a new Application for Jenkins. Click "Save application" when finished.
   1. Name: `Jenkins`
   2. Redirect URI: `https://jenkins.bigbang.dev/securityRealm/finishLogin`
   3. Trusted: Yes (checked)
   4. Confidential: Yes (checked)
   5. Expire access tokens: Yes (checked)
   6. Scopes: "api" checked, all others unchecked

1. Copy/Paste the Application ID and Secret from Gitlab into your config repo in the file `kustomizations/softwarefactoryaddons/jenkins/environment-bb-values.yaml` in the parameters that say `YOUR_CLIENT_ID_HERE` AND `YOUR_CLIENT_SECRET_HERE`

1. :warning: Encrypt `kustomizations/softwarefactoryaddons/jenkins/environment-bb-values.yaml` with SOPS. See the [SOPS guide](sops.md) for instructions.

    ```shell
    sops -e -i kustomizations/softwarefactoryaddons/jenkins/environment-bb-values.yaml
    ```

1. Commit and push the changes to your config repo

1. Create a "Day 2" package and deploy it. This package contains nothing but your config repo, so that Gitea will receive the new commit that you just pushed. For convenience, there is a Makefile in that repo

    ```shell
    cd day2
    zarf package create --confirm
    
    # Sneakernet the package if you need to

    zarf package deploy zarf-package-day-two-update-amd64.tar.zst --confirm
    ```

After Flux reconciles the change, Jenkins should now be using GitLab as the SSO provider.
