#!/bin/bash

docker build -t alexshemesh/kube-iptables-tailer .
docker tag alexshemesh/kube-iptables-tailer alexshemesh/kube-iptables-tailer:$TRAVIS_BUILD_NUMBER
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push alexshemesh/kube-iptables-tailer:$TRAVIS_BUILD_NUMBER
