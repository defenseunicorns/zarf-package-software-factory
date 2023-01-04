# zarf-package-software-factory
Pre-built Zarf Package of a Software Factory (a.k.a. "DI2-ME")

:construction: **This project is still early in its development.** Early adopters are expected to have a good understanding of how Kubernetes works as well as how to operate and maintain the various services present in this package. This package is intended to significantly speed up the initial deployment and configuration of a software factory, as well as provide a simple upgrade path for updates and security fixes. It will not absolve the user from understanding what it takes to operate a software factory in a production setting. We plan to transition to an open beta soon™ (maybe _sooner_ if any smarties out there want to come help out :smile:), where we will anticipate people starting to use this package operationally, with the caveat that it will require a lot more care and feeding than we intend it to require by the time version 1.0 releases.

Deploys the components of a software factory with the following services, all running on top of Big Bang Core:

- GitLab*
- GitLab Runner*
- Minio Operator*
- Jira
- Confluence
- Jenkins
- Artifactory
- HA Redis with Sentinel

Coming Soon:

- SonarQube*
- Mattermost*

**Deployed using Big Bang Umbrella*

![warning](img/warning.png)

## Zarf Compatibility
All versions of this package will not be compatible with all versions of Zarf. For maximum repeatability each version of this package expects a specific version of Zarf. This is taken care of by using `make` targets that download the correct version of Zarf to the `/build` directory. When you run `make all` you're downloading the right version of Zarf automatically.

## Known Issues

- :warning: We do not currently test compatibility from one version to the next. The user of this package is expected to first deploy to a test environment when doing upgrades. We will start testing upgrade paths as we get closer to a v1.0 release.

- Twistlock is disabled for now while we determine how we'll automatically enable it and test it when it requires a license key to be present in the values yaml.

- Several services have Istio sidecar injection disabled until we can determine how to enable it without breaking functionality with Postgres Operator.

- Jenkins won't work in disconnected environments due to its dependency on plugins pulled from the internet. Work is needed to figure out and implement a method of doing locally sourced plugins. Jira and Confluence plugins will have the same issue.

- Due to issues with Elasticsearch this package doesn't work yet in some k8s distros. It does work in K3s (using `zarf init --components k3s,gitops-service`). Upcoming work to swap the EFK stack out for the PLG stack (Promtail, Loki, Grafana) should resolve this issue. Keep in mind the big note above about the package being huge. Unless you have the Mother of All Laptops you'll need to turn a bunch of stuff off before you try to deploy it locally using Vagrant.

- If you are using Vagrant, inside the Vagrant VM the services are available on the standard port 443. Outside the VM if you want to pull something up in your browser that traffic is being routed to port **8443** to avoid needing to be root when running the Vagrant box.

- The current deployment of HA Redis with Sentinel is deployed with no authentication. Please reference [ADR-002](doc/adr/0002-switch-to-authless-ha-redis.md) for additional details.

## Prerequisites

- [Logged into registry1.dso.mil](https://github.com/defenseunicorns/zarf/blob/master/docs/ironbank.md)
- `make` present in PATH
- `sha256sum` present in PATH
- TONS of CPU and RAM. Our testing shows the AWS EC2 instance type m6i.8xlarge works pretty well at about $1.50/hour, which can be reduced further if you do a spot instance.
- [Vagrant](https://www.vagrantup.com/) and [VirtualBox](https://www.virtualbox.org/), only if you are going to use a Vagrant VM, which is incompatible when using an EC2 instance.

Note that having Zarf installed is not a prerequisite. This repo pulls in its own version of Zarf so that it can be versioned separately from whatever your system has installed. To change the version of Zarf used modify the ZARF_VERSION variable in the Makefile

## Documentation

To run this package in "self-contained" mode, follow these steps.

> "self-contained" mode means the package will deploy its own PostgreSQL databases and S3 buckets, using Zalando's Postgres Operator and MinIO, respectively.

1. [Fork the repo, customize, and build the packages](doc/fork-and-build.md)
1. [Deploy](doc/deploy.md)
1. [Configure SOPS encryption](doc/sops.md)
1. [Configure Single Sign-On](doc/sso.md)
1. [Day-2 Ops/Maintenance/Upgrades](doc/day2.md)
1. [Backup and Restore](doc/backup-and-restore/README.md)
1. [Troubleshooting](doc/troubleshooting.md)

If you want to customize your deployment by utilizing external PostgreSQL databases and S3 buckets, you can additionally follow along with these docs:

1. [Disable Postgres Operator and MinIO](doc/disable-postgres-operator-and-minio.md)
1. [Configure GitLab to use an external database and S3 buckets](doc/configure-gitlab-to-use-an-external-database-and-s3-buckets.md)
1. [Configure Jira to use an external database](doc/configure-jira-to-use-an-external-database.md)
1. [Configure Confluence to use an external database](doc/configure-confluence-to-use-an-external-database.md)
