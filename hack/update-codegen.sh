#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

TEMP_DIR=`mktemp -d`
SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PACKAGE=$TEMP_DIR/code-generator

git clone -b kubernetes-1.13.1 https://github.com/kubernetes/code-generator $CODEGEN_PACKAGE

export GOPATH=~/go
export GO111MODULE=on

${CODEGEN_PACKAGE}/generate-groups.sh all \
  github.schibsted.io/finn/fiaas-skipshund/pkg/client \
  github.schibsted.io/finn/fiaas-skipshund/pkg/apis \
  "fiaas.schibsted.io:v1" \
  --output-base ${TEMP_DIR} \
  --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt

echo "Syncing files back to repository root..."
rsync -av ${TEMP_DIR}/github.schibsted.io/finn/fiaas-skipshund/ $SCRIPT_ROOT/

rm -rf $TEMP_DIR
