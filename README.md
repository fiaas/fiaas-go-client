# FIAAS Custom Resource Definition

This library contains a Custom Resource Definition and an auto-generated Go client used to manipulating the Application FIAAS resource.

Code responsible for generating this CRD is present in `pkg/apis/fiaas.schibsted.io/v1/`:

- `types.go` was modeled after https://github.schibsted.io/finn/fiaas-deploy-daemon/blob/master/fiaas_deploy_daemon/crd/types.py
- `register.go` contains... necessary boilerplate code.
- `doc.go` contains some magic comments.

Everything in `pkg/client` and in `pkg/apis/fiaas.schibsted.io/v1/zz_generated.deepcopy.go` is completely auto generated. **Please do not edit those files manually.** Update any of the files listed above and run the code generation using `make generate-code`. Then, commit your changes and push.

`make generate-code` uses an `update-codegen` script adapted from the one Kubernetes maintainers use.

## Usage

<!-- TODO: Add example code -->
