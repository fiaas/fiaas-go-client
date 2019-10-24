#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME=fiaas-go-client-e2e

__cleanup() {
    kind delete cluster --name "$CLUSTER_NAME"
}
trap __cleanup EXIT

kind create cluster --name "$CLUSTER_NAME" --wait 30s
KIND_KUBECONFIG="$(kind get kubeconfig-path --name "$CLUSTER_NAME")"
kubectl --kubeconfig="$KIND_KUBECONFIG" apply -f hack/fiaas-crds.yml

export KIND_KUBECONFIG
go test --tags=e2e ./...
