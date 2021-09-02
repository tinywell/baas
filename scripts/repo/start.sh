#!/bin/sh

docker rm -f chartrepo

docker run -d \
  -p 8000:8080 \
  -e DEBUG=1 \
  -e STORAGE=local \
  -e STORAGE_LOCAL_ROOTDIR=/charts \
  -v /tmp/baas/repo:/charts \
  --name chartrepo \
  chartmuseum/chartmuseum:latest