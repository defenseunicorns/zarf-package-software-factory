# Jira And Confluence Default To RWO For Shared Home Volumes

Date: 2023-05-30

## Status

Accepted

## Context

By default Jira and Confluence require the "shared-home" volume be a RWX volume to facilitate high availability across multiple nodes. However, kubernetes does not come with a RWX storageclass and standing an RWX storageclass up is out of scope for the software factory package. This means that without configuration Jira and Confluence will run in a single replica config without persistence for the shared-home volume as there is no RWX storageclass.

From Jira and Confluence docs:

```
The shared-home folder in Confluence Data Center stores the Confluence attachments, backups, rendered previews, among other content.
```

```
The shared-home folder in Jira Data Center stores the Jira attachments, backups, and other content. 
```

When this volume is not persisted whenever the pod restarts the data mentioned above is lost.

## Decision

The default configuration in this package is that Jira and Confluence use a custom PVC that specifies an RWO volume instead. This value in `kustomizations/softwarefactoryaddons/<jira or confluence>/values.yaml` can be overridden or the PVC in `kustomizations/softwarefactoryaddons/base/pvcs/<jira or confluence>-shared-home-single-node.yaml` can be modified.

## Consequences

By default neither Jira or Confluence support HA now. However data is persisted by default so that users of the package don't accidentally delete user data if they didn't configure Confluence and Jira.
