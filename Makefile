.PHONY: e2e

# This version should be kept in sync with the most recent version in the CI matrix in .semaphore/semaphore.yaml
K8S_VERSION ?= v1.24.7
ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

generate-code:
	${ROOT_DIR}/hack/update-codegen.sh

verify:
	${ROOT_DIR}/hack/verify-codegen.sh

e2e:
	${ROOT_DIR}/hack/e2e-test $(K8S_VERSION)

test:
	go test ./... -v -cover --vet=all
