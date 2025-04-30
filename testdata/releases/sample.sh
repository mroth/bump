#!/bin/bash
if [ "$#" -ne 2 ]; then
    echo "usage: sample.sh <owner> <repo>"
    exit 1
fi

TS=$(date +"%Y%m%d")
OWNER=$1
REPO=$2

curl "https://api.github.com/repos/${OWNER}/${REPO}/releases" \
    > "${OWNER}-${REPO}-${TS}.sample.json"
