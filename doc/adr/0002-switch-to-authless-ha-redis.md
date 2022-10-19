# 2. Switch to Authless HA Redis

Date: 2022-10-19

## Status

Accepted

## Context

[Due to a customer concern](https://github.com/defenseunicorns/zarf-package-software-factory/issues/318) about the stability of the default Redis deployment that GitLab offers and GitLab's own docs that the embedded Redis deployment should not be used for production deployments, the decision was made to deploy an HA enabled Redis with Sentinel.

During the implementation of HA Redis, an issue was encountered with the ability for Gitlab to use HA Redis with authentication enabled. [Link to GitLab issue here](https://gitlab.com/gitlab-org/charts/gitlab/-/issues/2902)

Three options were discussed:

1. Deploy HA Redis without authentication required
   - AuthPolicies and/or NetworkPolicies can be used to limit the inbound/outbound traffic to the redis namespace.
   - Add a new issue to revisit the ability to enable authentication on the HA Redis deployment.
2. Keep the existing deployment of Redis that GitLab deploys, with the caveat that it is single node and not recommended for production use.
3. Continue to beat our heads against the wall with little to show for it.

## Decision

We went with Option 1, with some additional action items to take in the future.

Action Items:
1. New issue for limiting network access to Redis such that only Gitlab is able to connect to it.
2. New issue to revisit the ability to enable authentication on the HA Redis deployment.

## Consequences

1. Redis is more fault tolerant and can withstand node failures and other cluster issues.
2. There is slightly less defense in depth being used to secure Redis due to the elimination of the password. However, we believe this is an acceptable course of action because there are still several layers of security in front of Redis such as limiting its exposure to only internal traffic and creation of AuthPolicies/NetworkPolicies such that Redis will only be accessible by Gitlab.
