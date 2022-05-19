# Configure Single Sign-On

NOTE: `bigbang.dev` is the default domain. If you are using a different domain, substitute `bigbang.dev` for your domain in all URLs

## Configure GitLab

1. Retrieve the initial root password for GitLab:

    ```shell
    kubectl get secret gitlab-gitlab-initial-root-password -n gitlab -o jsonpath='{.data.password}' | base64 --decode
    ```

1. Navigate to [https://gitlab.bigbang.dev](https://gitlab.bigbang.dev)

1. Log in using username `root` with the password retrieved from the previous step

1. Navigate to [https://gitlab.bigbang.dev/-/profile/password/edit](https://gitlab.bigbang.dev/-/profile/password/edit) and change the root password. Save the new password as you will need it in disaster recovery scenarios.

1. [OPTIONAL] Disable Sign-up in the Sign-up restrictions section on [https://gitlab.bigbang.dev/admin/application_settings/general](https://gitlab.bigbang.dev/admin/application_settings/general). If you disable it you will need to manually create all new users. It may be advantageous to leave it on, since you can require admin approval for new sign-ups. Click "Save Changes" at the bottom of the section.

1. Enable "Enforce two-factor authentication" in the Sign-in restrictions section on [https://gitlab.bigbang.dev/admin/application_settings/general](https://gitlab.bigbang.dev/admin/application_settings/general). Click "Save Changes" at the bottom of the section.

1. Configure two-factor authentication on the root account. Make sure this gets done right away. If you wait past the grace period the root account will be locked out.

### Configure Jenkins

**Warnings and notes**

- :warning: WARNING: This section involves committing and pushing a secret to your config repo. The file needs to be encrypted with SOPS before you push it to your repo. See the [SOPS guide](sops.md) for instructions.

- NOTE: Configuring Jenkins to use GitLab as the SSO provider will not work when using the `bigbang.dev` domain. During the OAuth ping-pong auth flow Jenkins tries to run a POST http request to get a login token, but it isn't able to hit `https://gitlab.bigbang.dev` since that resolves to `127.0.0.1`. This may be fixable through the use of an Istio [ServiceEntry](https://istio.io/latest/docs/reference/config/networking/service-entry/). More investigation needed.

#### Instructions

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

### Configure Jira

#### Notes

- A limitation of Jira's SSO configuration is that you still need to create Jira users for each GitLab user that is going to use Jira. This may be able to be resolved using Jira's "Just-In-Time" user creation, but so far we haven't been able to get that working. The Jira user's username must be the same as the GitLab user that is logging in.
- The SSO login for Jira is weird. When you go to the Jira website you then have to click the "Login" button in the top right of the page.

#### Instructions

1. Create a new Application in GitLab: Admin > Applications > New application
   1. Name: `Jira`
   1. Redirect URI: `https://jira.bigbang.dev/plugins/servlet/oidc/callback`
   1. Trusted: Yes (checked)
   1. Confidential: Yes (checked)
   1. Expire access tokens: Yes (checked)
   1. Scopes: "openid" checked, all others unchecked

1. Navigate to [https://jira.bigbang.dev](https://jira.bigbang.dev)

1. Set up Jira. This guide will use the "I'll set it up myself" option.

   1. Click "I'll set it up myself" then click Next

   1. Use "Built In" for database connection. (Note: For eval purposes only. In production, Jira needs to use an external database. We're working on adding Zalando's Postgres Operator to the package.)

   1. Set up application properties
      1. Application Title: "Jira"
      1. Mode: Private
      1. Base URL: Should be prefilled. `https://jira.bigbang.dev`

   1. Specify your license key -- If you have one apply it. If not, click "generate a Jira trial license" to generate one

   1. Set up administrator account -- Use your own information, but make the username `root`. Keep the password that you set. You'll need to use it to do administrative stuff.

   1. Set up email notifications -- Choose "Later", unless you have your own SMTP server that you can use that is external to the kubernetes cluster

1. Create an initial project

1. Navigate to Settings > System > Authentication methods, click "Add configuration"
   1. Name: `GitLab`
   1. Authentication method: OpenID Connect single sign-on
   1. Issuer URL: `https://gitlab.bigbang.dev`
   1. Client ID: The Client ID from the GitLab Application you created earlier
   1. Client Secret: The Client Secret from the GitLab Application you created earlier
   1. Username mapping: `${nickname}`
   1. Additional scopes: none
   1. Additional settings / Fill the data automatically from my chosen identity provider: checked
   1. JIT provisioning / Create users on login to the application: unchecked
   1. Remember user logins: user's preference
   1. Show IdP on the login page: checked
   1. Login button text: user's preference. Recommend `Continue with GitLab`

### Configure Confluence

#### Notes

- A limitation of Confluence's SSO configuration is that you still need to create Confluence users for each GitLab user that is going to use Confluence. This may be able to be resolved using Confluence's "Just-In-Time" user creation, but so far we haven't been able to get that working. The Confluence user's username must be the same as the GitLab user that is logging in.

#### Instructions

1. Create a new Application in GitLab: Admin > Applications > New application
   1. Name: `Confluence`
   1. Redirect URI: `https://confluence.bigbang.dev/plugins/servlet/oidc/callback`
   1. Trusted: Yes (checked)
   1. Confidential: Yes (checked)
   1. Expire access tokens: Yes (checked)
   1. Scopes: "openid" checked, all others unchecked

1. Navigate to [https://confluence.bigbang.dev](https://confluence.bigbang.dev)

1. Apply your confluence license. If you don't have one it will let you create a trial license (requires internet access)

1. Choose your deployment type. Right now we only support "Standalone".

1. Load Content. This tutorial will use "Example Site", but you may want to choose to restore from a backup or create an empty site

1. Configure User Management. Choose "Manage Users and Groups within Confluence"

1. Configure System Administrator Account
   1. Username: `root`
   1. Name: User's preference
   1. Email: User's preference
   1. Password: User's preference.
   1. Confirm: Same as Password

1. Click "Further Configuration", then go to the "SSO 2.0" page
   1. Authentication method: OpenID Connect single sign-on
   1. Issuer URL: `https://gitlab.bigbang.dev`
   1. Client ID: The Client ID from the Application you made earlier in GitLab
   1. Client Secret: The Client Secret from the Application you made earlier in GitLab
   1. Additional scopes: none
   1. Username mapping: `${nickname}`
   1. Additional settings / Fill the data automatically from my chosen identity provider: checked
   1. JIT provisioning / Create users on login to the application: unchecked
   1. Remember user logins: user's preference
   1. Login mode: Use OpenID Connect as primary authentication
