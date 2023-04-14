#!/bin/bash

HOST_LIST=$(kubectl get vs -A -o=jsonpath='{range .items[*]}{.spec.gateways[*]}{" "}{.spec.hosts[*]}{"\n"}{end}' | sort -u)

PUBLIC_HOSTS=$(echo "${HOST_LIST}" | grep public | cut -d ' ' -f2)
PRIVATE_HOSTS=$(echo "${HOST_LIST}" | grep private | cut -d ' ' -f2)

PUBLIC_LB_IP=$(kubectl get svc -n istio-system public-ingressgateway -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')
PRIVATE_LB_IP=$(kubectl get svc -n istio-system private-ingressgateway -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')

sed -i -z 's/\n# Following entries are from metallb dns.sh.*# End of metallb dns.sh//' /etc/hosts

echo "# Following entries are from metallb dns.sh" >> /etc/hosts

echo "# Public hostnames" >> /etc/hosts

for host in $PUBLIC_HOSTS; do
    echo "${LB_IP} ${host}" >> /etc/hosts
done

echo "# Private hostnames" >> /etc/hosts

for host in $PRIVATE_HOSTS; do
    echo "${LB_IP} ${host}" >> /etc/hosts
done

echo "# End of metallb dns.sh" >> /etc/hosts
