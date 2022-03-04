# 1. Use GitLab as the SSO provider (but without CAC support)

Date: 2022-02-28

## Status

Accepted

## Context

To bring this Software Factory together as one solution, we need to use Single Sign-On so that a user can use the same account on all of the software factory's services.

Two options were discussed to accomplish this:

1. Use [Keycloak](https://repo1.dso.mil/platform-one/big-bang/apps/security-tools/keycloak) from [Big Bang Umbrella](https://repo1.dso.mil/platform-one/big-bang/bigbang)

1. Use [GitLab](https://repo1.dso.mil/platform-one/big-bang/apps/developer-tools/gitlab) from [Big Bang Umbrella](https://repo1.dso.mil/platform-one/big-bang/bigbang)

Keycloak is the incumbant. Platform One uses it, and a lot of customization work has been done to it to suit the needs that are common in a DoD environment, like support for Common Access Card (CAC). However, due to its architecture Keycloak is not able to run in the same Kubernetes cluster as the rest of the software factory. To use Keycloak we will need to run two clusters, one for Keycloak and one for the Software Factory.

The proposal was made that perhaps GitLab could act as the SSO provider. It is already part of the software factory, it is capable of acting as an SSO provider, and it would mean one less deployment to manage and only having to run one Kubernetes cluster. However, while GitLab documentation says it is capable of doing the things we expect to do, nobody that we know of has put together a system of this magnitude in a DoD environment, so we would be trailblazing this work with nothing to reference.

When faced with the option between GitLab and Keycloak as the SSO provider, the primary sticking point is CAC authorization. If CAC auth were not a hard requirement, the decision would be made much easier, since GitLab is the superior choice when it comes to operational complexity, but presents significantly elevated risk when figuring out CAC auth.

## Decision

- We will implement GitLab as the SSO provider, but without CAC auth, up until our "1.0" release.

- After the 1.0 release, we will again look at CAC auth as a requirement and adjust as necessary, by either figuring out how to do it in GitLab, or switching to Keycloak or some other SSO solution. Since SSO uses Open Standards OAuth 2.0 and OIDC, it actually is not very hard to swap out the SSO provider if we need to.

## Consequences

- We will be immediately unblocked from being able to wire up all the services with SSO, without having to first spend a lot of time investigating how to do CAC auth successfully.

- We will not have to worry about running a 2nd Kubernetes cluster for Keycloak, or deal with its operational complexity and brittleness.

- CAC auth will not be a supported feature in the near term, though it will still be on the roadmap. In the meantime the authentication method of Username+Password+2FA TOTP Token will be supported.
