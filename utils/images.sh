#!/bin/bash

touch images.txt

NEWIMAGES=$(kubectl get pods -A -o jsonpath='{range .items[*]}{.spec.containers[*].image}{"\n"}{end}')
IMAGELIST=$(cat images.txt)

echo -e "${NEWIMAGES}\n${IMAGELIST}" | sed 's/127.0.0.1:[0-9]*\///g' | sed 's/ /\n/g' | sort -u > images.txt

cat images.txt