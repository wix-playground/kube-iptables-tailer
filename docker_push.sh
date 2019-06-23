#!/bin/bash
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push alexshemesh/kube-iptables-tailer:$TRAVIS_BUILD_NUMBER
