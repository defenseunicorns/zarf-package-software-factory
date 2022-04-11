# zarf-package-software-factory
Pre-built Zarf Package of a Software Factory (a.k.a. "DI2-ME")

Deploys the components of a software factory with the following services, all running on top of Big Bang Core:

- GitLab*
- GitLab Runner*
- Minio Operator*
- Nexus*
- Jira
- Confluence
- Jenkins

Coming Soon:

- SonarQube*
- Mattermost*

**Deployed using Big Bang Umbrella*

![warning](img/warning.png)

## Zarf Compatibility
All versions of this package will not be compatible with all versions of Zarf. Here's a compatibility matrix to use to determine which versions match up

| Package Version                                                                                                                                                                   | Zarf Version                                                             |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |--------------------------------------------------------------------------|
| [0.0.2](https://github.com/defenseunicorns/zarf-package-software-factory/releases/tag/0.0.2) - [main](https://github.com/defenseunicorns/zarf-package-software-factory/tree/main) | [v0.15.0](https://github.com/defenseunicorns/zarf/releases/tag/v0.15.0)+ |
| [0.0.1](https://github.com/defenseunicorns/zarf-package-software-factory/releases/tag/0.0.1)                                                                                      | [v0.14.0](https://github.com/defenseunicorns/zarf/releases/tag/v0.14.0)  |

## Known Issues

- Jenkins won't work in disconnected environments due to its dependency on plugins pulled from the internet. Work is needed to figure out and implement a method of doing locally sourced plugins.

- Due to issues with Elasticsearch this package doesn't work yet in some k8s distros. It does work in K3s (using `zarf init --components k3s,gitops-service`). Upcoming work to swap the EFK stack out for the PLG stack (Promtail, Loki, Grafana) should resolve this issue. Keep in mind the big note above about the package being huge. Unless you have the Mother of All Laptops you'll need to turn a bunch of stuff off before you try to deploy it locally using Vagrant.

- Inside the Vagrant VM the services are available on the standard port 443. Outside the VM if you want to pull something up in your browser that traffic is being routed to port **8443** to avoid needing to be root when running the Vagrant box.

- Secrets like the TLS cert are currently stored in a ConfigMap, they need to be moved to a Secret.

## Prerequisites

- [Logged into registry1.dso.mil](https://github.com/defenseunicorns/zarf/blob/master/docs/ironbank.md)
- Clone this repo
- `make` present in PATH
- `sha256sum` present in PATH
- TONS of CPU and RAM. Our testing shows the AWS EC2 instance type m6i.8xlarge works pretty well at about $1.50/hour, which can be reduced further if you do a spot instance.
- [Vagrant](https://www.vagrantup.com/) and [VirtualBox](https://www.virtualbox.org/), only if you are going to use a Vagrant VM, which is incompatible when using an EC2 instance.

Note that having Zarf installed is not a prerequisite. This repo pulls in its own version of Zarf so that it can be versioned separately from whatever your system has installed. To change the version of Zarf used modify the ZARF_VERSION variable in the Makefile

## Instructions

```sh
# Go into the cloned repo
cd zarf-package-software-factory

# Build everything, it all will get put in the 'build' folder
make all
```

## Customize

There's lots of different customization you can do, but the most important one will be to update `kustomizations/bigbang/values.yaml` to change the `domain: bigbang.dev` to your actual domain, and update the TLS Cert and Key in the `istio:` section to your actual cert and key. _**YOUR TLS CERT KEY MUST BE TREATED AS A SECRET**_. It is present here because `bigbang.dev` has its A record configured to point at `127.0.0.1` to make development and testing easier. Further development is planned to move the istio cert configuration to a Kubernetes Secret instead of a ConfigMap.

## Next Steps

Depending on where you want to run the package you just created, there are a few different paths

### Airgap

1. Burn to removable media
    - `build/zarf`
    - `build/zarf-init-amd64.tar.zst`
    - `build/zarf-package-flux-amd64.tar.zst`
    - `build/zarf-package-software-factory-amd64.tar.zst`

1. Use [Sneakernet](https://en.wikipedia.org/wiki/Sneakernet) or whatever other method you want to get it where it needs to go

1. Deploy (init, then deploy Flux, then deploy software factory package)

### Vagrant

1. Update the package to disable enough services so that it will fit on your computer. Big Bang deployments can be disabled by changing `enabled: true` to `enabled: false` in `kustomizations/bigbang/values.yaml`. Other addons like Jira, Confluence, and Jenkins can be disabled by modifying the Kustomization in the `kustomizations/softwarefactoryaddons` directory.

2. Rebuild the package -- `make all`

3. Spin up the VM with `make vm-init`. The `build` folder will be mounted into the VM

4. SSH into the VM with `make vm-ssh`

5. Deploy (init, then deploy Flux, then deploy software factory package)

6. When you're done, bring everything down with `make vm-destroy`

### Cloud

1. Copy files to the machine that will be running `zarf`. Could be your local computer, or could be a Bastion Host, or something else entirely.
   - `build/zarf`
   - `build/zarf-init-amd64.tar.zst`
   - `build/zarf-package-flux-amd64.tar.zst`
   - `build/zarf-package-software-factory-amd64.tar.zst`

1. (optional, unstable) Have a Kubernetes cluster ready that you'll be deploying to. Have your KUBECONFIG be configured such that running something like `kubectl get nodes` will connect to the right cluster. OR use the built-in K3s that Zarf comes bundled with.

1. Deploy

## Troubleshooting

### Elasticsearch is unhealthy

Make sure `sysctl -w vm.max_map_count=262144` got run. Elasticsearch needs it to function properly.
