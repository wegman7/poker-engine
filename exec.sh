#!/usr/bin/env bash

# build for amd64
docker buildx build --platform linux/amd64 \
  -t gcr.io/poker-451119/engine:v1 \
  --push .
docker run -d -p 8080:8080 --env-file=.env.prod gcr.io/poker-451119/engine:v1 --env=prod
