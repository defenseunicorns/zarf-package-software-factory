#!/bin/bash
############
# This script can be used to modify an existing CoreDNS configmap Corefile by adding a 'rewrite' rule to the main domain's configuration.
# Expectations:
#   * This should only be run on a default installation in dev/test environments
#   * Kubectl must be installed and configured
#   * Basic bash utils (wc, grep...) must be available

#######
# VARS
#######
REWRITE_CMD='rewrite name gitlab.bigbang.dev gitlab-webservice-default.gitlab.svc.cluster.local'
TMP_FILE=tmp_cm.yaml
NAMESPACE=kube-system
CONFIGMAP=coredns
DEPLOYMENT=coredns
MATCHSTRING="/loadbalance"

########
# FUNCS
########
gracefulExit(){
  echo "Something isn't quite right.."
  exit 0
}

#########
# SCRIPT
#########
# check if kubectl exists
if ! command -v kubectl &> /dev/null
then
  echo "kubectl not found - bailing"
  gracefulExit
fi

# check if coredns configmap exists
if [ ! $(kubectl get cm -n $NAMESPACE | grep $CONFIGMAP | wc -l) -gt 0 ]; then
  printf "Expected CoreDNS configmap '%s' not found - bailing\n" $CONFIG_MAP
  gracefulExit
fi

# check if the rewrite cmd already exists
if [ $(kubectl get cm -n $NAMESPACE $CONFIGMAP -o yaml | grep "$REWRITE_CMD" | wc -l) -gt 0 ]; then
  echo "Target rewrite plugin already exists in the configmap - bailing"
  gracefulExit
fi

# cleanup temp file if it exists
rm -f $TMP_FILE

# get existing data values from the configmap
_NODEHOSTS=$(kubectl get cm -n $NAMESPACE $CONFIGMAP -o jsonpath='{ .data.NodeHosts }')
_COREFILE=$(kubectl get cm -n $NAMESPACE $CONFIGMAP -o jsonpath='{ .data.Corefile }')

# modify the Corefile with the rewrite string
_COREFILE=$(echo "$_COREFILE" | sed "$MATCHSTRING/a \ \ \ \ $REWRITE_CMD")

# build a new configmap (start with the boilerplate)
cat << EOF > $TMP_FILE
apiVersion: v1
kind: ConfigMap
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: ""
  name: coredns
  namespace: kube-system
data:
  NodeHosts: |
$(while IFS= read -r line; do printf '%4s%s\n' '' "$line"; done <<< "$_NODEHOSTS")
  Corefile: |
$(while IFS= read -r line; do printf '%4s%s\n' '' "$line"; done <<< "$_COREFILE")
EOF

# apply the configmap
echo "Attempting to apply the following ConfigMap:"
echo "$TMP_FILE"
kubectl apply -f $TMP_FILE

# restart coredns
kubectl rollout restart -n $NAMESPACE deployment/$DEPLOYMENT

# cleanup the tmp file
rm -f $TMP_FILE