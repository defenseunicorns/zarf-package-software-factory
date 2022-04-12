# Manual SSO configuration

This folder contains the steps for manual configuration of SSO, which will be provided while we work on automating the declarative configuration.

## Procedure

### Deploy the software factory

```shell
zarf init --components k3s,gitops-service --confirm
./zarf package deploy zarf-package-flux-amd64.tar.zst --confirm
./zarf package deploy zarf-package-software-factory-amd64.tar.zst --confirm
```

### Configure GitLab

1. Get the initial root password by running
    ```shell
    
    ```
3. Navigate to GitLab (default is [https://gitlab.bigbang.dev](https://gitlab.bigbang.dev))
4. Log in using 
