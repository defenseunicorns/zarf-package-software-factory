#!/bin/bash

echo "Installing MetalLB..."

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.9/config/manifests/metallb-native.yaml 2>&1 >/dev/null

echo "Waiting for MetalLB to be ready..."

kubectl rollout status deployment/controller -n metallb-system --watch --timeout=90s 2>&1 >/dev/null

kubectl rollout status daemonset/speaker -n metallb-system --watch --timeout=90s 2>&1 >/dev/null

NETWORK_BASE=$(docker network inspect -f '{{(index .IPAM.Config 0).Subnet}}' kind | cut -d '.' --fields '1 2 3')

NETWORK_RANGE_START="${NETWORK_BASE}.200"
NETWORK_RANGE_END="${NETWORK_BASE}.210"

METALLB_CONFIG=$(sed "s/#RANGE_START#/${NETWORK_RANGE_START}/g" test/metallb/config.yaml | sed "s/#RANGE_END#/${NETWORK_RANGE_END}/g")

echo "Configuring MetalLB..."

echo "$METALLB_CONFIG" | kubectl apply --wait=true -f - 2>&1 >/dev/null
