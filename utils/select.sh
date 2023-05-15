#!/bin/bash

IMAGES=$(cat currentImages.txt)

while read -r p; do
  IMAGES=$(echo "${IMAGES}" | sed "s~\b${p}\b~& match~")
done <images.txt

echo "Matching images"
echo "${IMAGES}" | grep ' match"' | sed 's/ match//g'
echo "Non-matching images"
echo "${IMAGES}" | grep -v ' match"'