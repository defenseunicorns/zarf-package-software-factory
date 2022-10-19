# 2. Switch to Authless HA Redis

Date: 2022-10-19

## Status

Accepted

## Context

Due to a customer concern about the stability of the default Redis deployment GitLab offers and GitLab's own docs that the embedded Redis deployment should not be used for production deployments, the decision was made to deploy an HA enabled Redis with Sentinel.

During the implementation of HA Redis, an issue was encountered with the ability for Gitlab to use HA Redis with authentication enabled. [Link to GitLab issue here](https://gitlab.com/gitlab-org/charts/gitlab/-/issues/2902)

## Decision

Three options were discussed:

1. Deploy HA Redis without authentication required
   - AuthPolicies and/or NetworkPolicies can be used to limit the inbound/outbound traffic to the redis namespace.
   - Add a new issue to revisit the ability to enable authentication on the HA Redis deployment.
2. Keep the existing deployment of Redis that GitLab deploys, with the caveat that it is single node and not recommended for production use.

## Consequences

What becomes easier or more difficult to do and any risks introduced by the change that will need to be mitigated.

With the deployment of HA Redis with Sentinal, we will have a more fault-tolerant Redis deployment that can withstand node failures and other cluster issues.
