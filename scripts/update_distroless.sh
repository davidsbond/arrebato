#!/bin/bash

set -e

docker pull gcr.io/distroless/static

LATEST_DIGEST=$(docker images --format '{{json .}}' --digests gcr.io/distroless/static | jq -r 'select (.Tag=="latest") | .Digest')
CURRENT_DIGEST=$(cat Dockerfile | grep FROM | awk -F'@' '{print $2}')

if [ "$LATEST_DIGEST" != "$CURRENT_DIGEST" ]; then
  sed -i "s/$CURRENT_DIGEST/$LATEST_DIGEST/" Dockerfile
fi
