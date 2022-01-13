.PHONY: e2e

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

generate-code:
	${ROOT_DIR}/hack/update-codegen.sh

verify:
	${ROOT_DIR}/hack/verify-codegen.sh

e2e:
	${ROOT_DIR}/hack/e2e-test v1.20.7
	${ROOT_DIR}/hack/e2e-test v1.21.2
	${ROOT_DIR}/hack/e2e-test v1.22.4
