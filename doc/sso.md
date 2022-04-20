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

:warning: WARNING: This section involves committing and pushing a secret to your config repo. Eventually we plan to support SOPS encryption of secrets, but that's not been documented or configured yet. If you are committing secrets to your git repo make sure the repo is private and the entire repo (and resulting zarf package) is treated as a secret. We highly recommend not committing unencrypted operational secrets to git repos, even when they are private. This project will stay in version "0.0.X" until this is addressed, signifying that it is not ready for production and should only be used in dev/test/kick-the-tires types of environments.

1. Navigate to [https://gitlab.bigbang.dev/admin/applications/new](https://gitlab.bigbang.dev/admin/applications/new) and create a new Application for Jenkins. Click "Save application" when finished.
   1. Name: `Jenkins`
   2. Redirect URI: `https://jenkins.bigbang.dev/securityRealm/finishLogin`
   3. Trusted: Yes (checked)
   4. Confidential: Yes (checked)
   5. Expire access tokens: Yes (checked)
   6. Scopes: "api" checked, all others unchecked

2. Copy/Paste the Application ID and Secret from Gitlab into your config repo in the file `kustomizations/softwarefactoryaddons/jenkins-common-values.yaml` in the parameters that say `YOUR_CLIENT_ID_HERE` AND `YOUR_CLIENT_SECRET_HERE`

3. Uncomment the two blocks that are labeled with `WHEN SWITCHING TO GITLAB SSO UNCOMMENT THIS SECTION`

4. Commit and push the changes to your config repo

5. Create a "Day 2" package and deploy it. This package contains nothing but your config repo, so that Gitea will receive the new commit that you just pushed. For convenience, there is a Makefile in that repo

```shell
cd day2
make build-and-deploy
```

After Flux reconciles the change, Jenkins should now be using GitLab as the SSO provider.