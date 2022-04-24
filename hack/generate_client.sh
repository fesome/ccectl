#!/usr/bin/env bash

go mod vendor
retVal=$?
if [ $retVal -ne 0 ]; then
    exit $retVal
fi

set -e

GO111MODULE=on /bin/bash vendor/k8s.io/code-generator/generate-groups.sh all \
github.com/xiaohuifirst/ccectl/client github.com/xiaohuifirst/ccectl/apis "calico:v1" -h ./hack/boilerplate.go.txt

GO111MODULE=on /bin/bash vendor/k8s.io/code-generator/generate-internal-groups.sh defaulter \
github.com/xiaohuifirst/ccectl/client github.com/xiaohuifirst/ccectl/apis github.com/xiaohuifirst/ccectl/apis "calico:v1" -h ./hack/boilerplate.go.txt

rm -rf ./vendor
