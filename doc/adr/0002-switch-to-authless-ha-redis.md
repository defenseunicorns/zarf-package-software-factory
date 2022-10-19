# 2. Switch to Authless HA Redis

Date: 2022-10-19

## Status

Accepted

## Context

Due to a customer concern about the stability of the default redis deployment gitlab offers and gitlab's own docs that the embedded redis deployment should not be used for production deployments, the decision was made to deploy an HA enabled redis with sentinel.

During the implementation of HA redis, an issue was encountered with the ability for Gitlab to use HA redis with authication enabled. [link to gitlab issue here](https://gitlab.com/gitlab-org/charts/gitlab/-/issues/2902)

## Decision

The change that we're proposing or have agreed to implement.

## Consequences

What becomes easier or more difficult to do and any risks introduced by the change that will need to be mitigated.
