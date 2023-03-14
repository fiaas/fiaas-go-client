# FIAAS Custom Resource Definition

This library contains a Custom Resource Definition and an auto-generated Go client used to manipulate the Application FIAAS resource.

Code responsible for generating this CRD is present in `pkg/apis/fiaas.schibsted.io/v1/`:

- `types.go` was modeled after https://github.com/fiaas/fiaas-deploy-daemon/blob/master/fiaas_deploy_daemon/crd/types.py.
- `register.go` contains... necessary boilerplate code.
- `doc.go` contains some magic comments.

Everything in `pkg/client` and in `pkg/apis/fiaas.schibsted.io/v1/zz_generated.deepcopy.go` is completely auto generated. **Please do not edit those files manually.** Update any of the files listed above and run the code generation using `make generate-code`. Then, commit your changes and push.

`make generate-code` uses an `update-codegen` script adapted from the one Kubernetes maintainers use.

## Type structure

In `types.go`, we made a deliberate choice to treat the `Config` struct field of the `ApplicationSpec` as a "black box" in the form of `map[string]interface`. The fiaas-deploy-daemon client will accept any config structure as long as it's valid YAML. This gives us the benefit of being backwards- and forwards compatible.

## Usage

<!-- TODO: Add example code -->


## Tests

To run end-to-end tests, [Kind](https://github.com/kubernetes-sigs/kind#installation-and-usage) and [Docker](https://docs.docker.com/install/) must be installed. The end-to-end tests can be run with `make e2e`. To runs e2e tests against a specific Kubernetes version, set `$K8S_VERSION`, e.g. `K8S_VERSION=vX.Y.Z make e2e`.
