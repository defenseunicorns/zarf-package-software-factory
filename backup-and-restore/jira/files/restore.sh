#!/bin/bash

echo "Currently Jira is scaled to zero replicas and the jira-database helmrelease/clutser has been deleted"
echo "This script will walk you through restoring data from this package to the jira database and bringing Jira back up"

yq '.spec.values.resources[0].spec.clone = {"cluster": "test", "uid": "testuid"}' kustomizations/softwarefactoryaddons/base/databases/jira.yaml