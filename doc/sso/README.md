# Manual SSO configuration

This folder contains the steps for manual configuration of SSO, which will be provided while we work on automating the declarative configuration.

## Procedure

NOTE: `bigbang.dev` is the default domain. If you are using a different domain, substitute `bigbang.dev` for your domain in all URLs

### Deploy the software factory

```shell
zarf init --components k3s,gitops-service --confirm
./zarf package deploy zarf-package-flux-amd64.tar.zst --confirm
./zarf package deploy zarf-package-software-factory-amd64.tar.zst --confirm
```

### Configure GitLab

1. Navigate to [https://gitlab.bigbang.dev](https://gitlab.bigbang.dev))
2. Log in using username `root` with password `Ch@ngeMe!`
3. Navigate to [https://gitlab.bigbang.dev/-/profile/password/edit](https://gitlab.bigbang.dev/-/profile/password/edit) and change the root password. Save the new password as you will need it in disaster recovery scenarios.
4. 
