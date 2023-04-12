#!/bin/bash

HOST_LIST=$(kubectl get vs -A -o=jsonpath='{range .items[*]}{.spec.hosts[*]}{"\n"}{end}' | sort -u)

LB_IP=$(kubectl get svc -n istio-system public-ingressgateway -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')

echo >> /etc/hosts
echo "# Following entries are from metallb dns.sh" >> /etc/hosts

for host in $HOST_LIST; do
    echo "${LB_IP} ${host}" >> /etc/hosts
done

echo "# End of metallb dns.sh" >> /etc/hosts
