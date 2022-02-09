# zarf-package-software-factory
Pre-built Zarf Package of a Software Factory (a.k.a. "DI2-ME")

Deploys the components of a software factory with the following services, all running on top of Big Bang Core:

- SonarQube*
- GitLab*
- GitLab Runner*
- Minio Operator*
- Mattermost Operator*
- Mattermost*
- Nexus*
- Jira
- Confluence
- Jenkins


**Deployed using Big Bang Umbrella*

![warning](img/warning.png)

## Zarf Compatibility
All version of this package will not be compatible with all versions of Zarf. Here's a compatibility matrix to use to determine which versions match up

| Package Version                                                                    | Zarf Version                                                            |
| ---------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| [main](https://github.com/defenseunicorns/zarf-package-software-factory/tree/main) | [v0.14.0](https://github.com/defenseunicorns/zarf/releases/tag/v0.14.0) |

## Known Issues

<!-- Not relevant until we update to Zarf v0.15.0 -->
<!-- - Due to issues with Elasticsearch this package doesn't work yet in some distros. It does work in the Vagrant VM detailed below. Upcoming work to update to the latest version of Big Bang and swap the EFK stack out for the PLG stack (Promtail, Loki, Grafana) should resolve this issue. Keep in mind the big note above about the package being huge. Unless you have the Mother of All Laptops you'll need to turn a bunch of stuff off before you try to deploy it locally using Vagrant. -->

<!-- Should be removed after we update to Zarf v0.15.0 -->
- Since Traefik and Istio are both running in the cluster, Istio has been set to expose its services on port 9443. The Vagrant VM forwards that port, so you can pull up services in your browser. For example, https://grafana.bigbang.dev:9443

<!-- Not relevant until we update to Zarf v0.15.0 -->
<!-- - Inside the Vagrant VM the services are available on the standard port 443. Outside the VM if you want to pull something up in your browser that traffic is being routed to port **8443** to avoid needing to be root when running the Vagrant box. -->

- Currently this package does the equivalent of `kustomize build | kubectl apply -f -`, which means Flux will be used to deploy everything, but it won't be watching a Git repository for changes. Upcoming work is planned to update the example so that you will be able to open up a Git repo in the private Gitea server inside the cluster, commit and push a change, and see that change get reflected in the deployment.

## Prerequisites

- [Zarf installed](https://github.com/defenseunicorns/zarf/blob/master/docs/workstation.md#just-gimmie-zarf)
- `zarf` present in PATH
- [Logged into registry1.dso.mil](https://github.com/defenseunicorns/zarf/blob/master/docs/ironbank.md#2-configure-zarf-the-use-em)
- Clone this repo
- `make` present in PATH
- `kustomize` present in PATH
- `sha256sum` present in PATH
- TONS of CPU and RAM. Our testing shows the EC2 instance type m6i.8xlarge works pretty well at about $1.50/hour, which can be reduced further if you do a spot instance.
- [Vagrant](https://www.vagrantup.com/) and [VirtualBox](https://www.virtualbox.org/), only if you are going to use a Vagrant VM, which is incompatible when using an EC2 instance.

## Instructions

```sh
# Go into the cloned repo
cd zarf-package-software-factory

# Build the package, it will get put in the 'build' folder
make build
```

## Customize

There's lots of different customization you can do, but the most important one will be to update `template/bigbang/values.yaml` to change the `domain: bigbang.dev` to your actual domain, and update the TLS Cert and Key in the `istio:` section to your actual cert and key. _**YOUR TLS CERT KEY MUST BE TREATED AS A SECRET**_. It is present here because `bigbang.dev` has its A record configured to point at `127.0.0.1` to make development and testing easier. Further development is planned to move the istio cert configuration to a Kubernetes Secret instead of a ConfigMap.


## Next Steps

Depending on where you want to run the package you just created, there are a few different paths

### Airgap

1. Burn to removable media
    - `zarf`
    - `zarf-init.tar.zst`
    - `zarf-package-software-factory.tar.zst`

1. Use [Sneakernet](https://en.wikipedia.org/wiki/Sneakernet) or whatever other method you want to get it where it needs to go

1. Deploy

### Vagrant

1. Edit `template/bigbang/values.yaml` and make the deployment much smaller than it is by default by changing a bunch of the `enabled: true` parameters to `enabled: false`. You can disable the Atlassian stack or Jenkins from `zarf.yaml`. Change `required: true` to `required:false` then press `N` when asked whether you want to deploy them.

1. Rebuild the package -- `make build`

1. Spin up the VM with `make vm-init`. The `build` folder will be mounted into the VM

1. Copy `zarf` and `zarf-init-tar.zst` into the `build` folder. `zarf-package-software-factory.tar.zst` should already be there.

1. SSH into the VM with `make ssh`

1. Deploy

1. When you're done, bring everything down with `make vm-destroy`

### Cloud

1. Copy files to the machine that will be running `zarf`. Could be your local computer, or could be a Bastion Host, or something else entirely.
    - `zarf`
    - `zarf-init.tar.zst`
    - `zarf-package-software-factory.tar.zst`

<!-- Not compatible until v0.15.0 -->
<!-- 1. Have a Kubernetes cluster ready that you'll be deploying to. Have your KUBECONFIG be configured such that running something like `kubectl get nodes` will connect to the right cluster. -->

1. Deploy

## Troubleshooting

### Elasticsearch is unhealthy

Make sure `sysctl -w vm.max_map_count=262144` got run. Elasticsearch needs it to function properly.
