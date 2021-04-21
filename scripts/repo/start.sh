#!/bin/sh

docker run -d \
  -p 8080:8080 \
  -e DEBUG=1 \
  -e STORAGE=local \
  -e STORAGE_LOCAL_ROOTDIR=/charts \
  -v /tmp/baas/repo:/charts \
  --name chartrepo \
  chartmuseum/chartmuseum:latest