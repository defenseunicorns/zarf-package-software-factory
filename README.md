# zarf-package-software-factory
Pre-built Zarf Package of a Software Factory (a.k.a. "DI2-ME")

This example deploys the components of a software factory with the following services, all running on top of Big Bang Core:

- SonarQube*
- GitLab*
- GitLab Runner*
- Minio Operator*
- Mattermost Operator*
- Mattermost*
- Nexus*
- Keycloak*
- Jira
- Confluence
- Jenkins

**Deployed using Big Bang Umbrella*

<span style="color:red; font-size:2em">This package is huge. We recommend not trying to run it on a developer laptop without disabling lots of stuff first.</span>

> Note: Right now the intention is to show that all of these services can be deployed easily using a single Zarf package. They are not configured (yet). You can't take this demo and deploy it expecting to have a fully operational software factory at the push of a button, though that is the end goal. There's a lot of work to do between what is here now and that end goal, some of which might just not make very much sense in the context of a demo/example.

## Prerequisites

- Logged into registry1.dso.mil
- `make`
- `kustomize`
- `sha256sum`
- TONS of CPU and RAM. Our testing shows the EC2 instance type m6i.8xlarge works pretty well at about $1.50/hour, which can be reduced further if you do a spot instance.
- [Vagrant](https://www.vagrantup.com/) and [VirtualBox](https://www.virtualbox.org/), only if you are going to use a Vagrant VM, which is incompatible when using an EC2 instance.

Note: Vagrant and VirtualBox aren't required for Zarf to function, but this example's Makefile uses them to create a VM which everything will run in. In production you'll likely just run Zarf on the machine itself.

## Instructions

1. `cd examples/software-factory`
1. Run one of these two commands:
   - `make all` - Download the latest version of Zarf, build the deploy package, and start a VM with Vagrant
   - `make all-dev` - Build Zarf locally, build the deploy package, and start a VM with Vagrant. Requires Golang.

     > Note: If you are in an EC2 instance you should skip the `vm-init` make target, so run `make clean fetch-release package-example-software-factory && cd ../sync && sudo su` instead, then move on to the next step.
1. Run: `./zarf init --confirm --components management,gitops-service --host 127.0.0.1` - Initialize Zarf, telling it to install the management component and gitops service and skip logging component (since BB has logging already) and tells Zarf to use `127.0.0.1` as the cluster's address. If you want to use interactive mode instead just run `./zarf init`.
1. Wait a bit, run `k9s` to see pods come up. Don't move on until everything is running
1. Run: `./zarf package deploy zarf-package-software-factory-demo.tar.zst --confirm` - Deploy the software factory package. If you want interactive mode instead just run `./zarf package deploy`, it will give you a picker to choose the package.
1. Wait several minutes. Run `k9s` to watch progress
1. :warning: `kubectl delete -n istio-system envoyfilter/misdirected-request` (due to [this bug](https://repo1.dso.mil/platform-one/big-bang/bigbang/-/issues/802))
1. Use a browser to visit the various services, available at https://*.bigbang.dev:9443
1. When you're done, run `exit` to leave the VM then `make vm-destroy` to bring everything down

## Notes

- If you are not running in a Vagrant box created with the Vagrantfile in ./examples you will have to run `sysctl -w vm.max_map_count=262144` to get ElasticSearch to start correctly.
- If you want to turn off certain services to help the package run on smaller machines go into `template/bigbang/values.yaml` and change `enabled: true` to `enabled: false` for each service you want to disable. You can disable the Atlassian stack or Jenkins from `zarf.yaml`. Change `required: true` to `required:false` then press `N` when asked whether you want to deploy them.

## Services

| URL                                                   | Username  | Password                                                                                                                                                                                   | Notes           |
| ----------------------------------------------------- | --------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------- |
| [AlertManager](https://alertmanager.bigbang.dev:9443) | n/a       | n/a                                                                                                                                                                                        | Unauthenticated |
| [Grafana](https://grafana.bigbang.dev:9443)           | `admin`   | `prom-operator`                                                                                                                                                                            |                 |
| [Kiali](https://kiali.bigbang.dev:9443)               | n/a       | `kubectl get secret -n kiali -o=json \| jq -r '.items[] \| select(.metadata.annotations."kubernetes.io/service-account.name"=="kiali-service-account") \| .data.token' \| base64 -d; echo` |                 |
| [Kibana](https://kibana.bigbang.dev:9443)             | `elastic` | `kubectl get secret -n logging logging-ek-es-elastic-user -o=jsonpath='{.data.elastic}' \| base64 -d; echo`                                                                                |                 |
| [Prometheus](https://prometheus.bigbang.dev:9443)     | n/a       | n/a                                                                                                                                                                                        | Unauthenticated |
| [Jaeger](https://tracing.bigbang.dev:9443)            | n/a       | n/a                                                                                                                                                                                        | Unauthenticated |
| [Twistlock](https://twistlock.bigbang.dev:9443)       | n/a       | n/a                                                                                                                                                                                        |                 |
| [Jira](https://jira.bigbang.dev:9443)                 | n/a       | n/a                                                                                                                                                                                        |                 |
| [Confluence](https://confluence.bigbang.dev:9443)     | n/a       | n/a                                                                                                                                                                                        |                 |
| [GitLab](https://gitlab.bigbang.dev:9443)             | n/a       | n/a                                                                                                                                                                                        |                 |
| [Nexus](https://nexus.bigbang.dev:9443)               | n/a       | n/a                                                                                                                                                                                        |                 |
| [Mattermost](https://chat.bigbang.dev:9443)           | n/a       | n/a                                                                                                                                                                                        |                 |
| [Sonarqube](https://sonarqube.bigbang.dev:9443)       | n/a       | n/a                                                                                                                                                                                        |                 |
| [Jenkins](https://jenkins.bigbang.dev:9443)           | `admin`   | `admin`                                                                                                                                                                                    |                 |
