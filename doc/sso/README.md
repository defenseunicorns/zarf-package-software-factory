# Manual SSO configuration

This folder contains the steps for manual configuration of SSO, which will be provided while we work on automating the declarative configuration.

## Procedure

NOTE: `bigbang.dev` is the default domain. If you are using a different domain, substitute `bigbang.dev` for your domain in all URLs

### Deploy the software factory

```shell
zarf init --components k3s,gitops-service --confirm
zarf package deploy zarf-package-flux-amd64.tar.zst --confirm
zarf package deploy zarf-package-software-factory-amd64.tar.zst --confirm
```

### Configure GitLab

1. Navigate to [https://gitlab.bigbang.dev](https://gitlab.bigbang.dev))
2. Retrieve the initial root password for GitLab:
    ```shell
    kubectl get secret gitlab-gitlab-initial-root-password -n gitlab -o jsonpath='{.data.password}' | base64 --decode
    ```
3. Log in using username `root` with the password retrieved from the previous step
4. Navigate to [https://gitlab.bigbang.dev/-/profile/password/edit](https://gitlab.bigbang.dev/-/profile/password/edit) and change the root password. Save the new password as you will need it in disaster recovery scenarios.
5. [OPTIONAL] Disable Sign-up in the Sign-up restrictions section on [https://gitlab.bigbang.dev/admin/application_settings/general](https://gitlab.bigbang.dev/admin/application_settings/general). If you disable it you will need to manually create all new users. It may be advantageous to leave it on, since you can require admin approval for new sign-ups. Click "Save Changes" at the bottom of the section.
6. Enable "Enforce two-factor authentication" in the Sign-in restrictions section on [https://gitlab.bigbang.dev/admin/application_settings/general](https://gitlab.bigbang.dev/admin/application_settings/general). Click "Save Changes" at the bottom of the section.
7. Configure two-factor authentication on the root account. Make sure this gets done right away. If you wait past the grace period the root account will be locked out.

### Configure Jenkins

1. 