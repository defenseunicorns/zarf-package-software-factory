#!/bin/bash

kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.9/config/manifests/metallb-native.yaml

kubectl wait --namespace metallb-system \
                --for=condition=ready pod \
                --selector=app=metallb \
                --timeout=90s

NETWORK_BASE=$(docker network inspect -f '{{(index .IPAM.Config 0).Subnet}}' kind | cut -d '.' --fields '1 2 3')

NETWORK_RANGE_START="${NETWORK_BASE}.200"
NETWORK_RANGE_END="${NETWORK_BASE}.210"

sed -i "s/#RANGE_START#/${NETWORK_RANGE_START}/g" test/metallb/config.yaml
sed -i "s/#RANGE_END#/${NETWORK_RANGE_END}/g" test/metallb/config.yaml

kubectl apply -f test/metallb/config.yaml