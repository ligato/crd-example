[![Build Status](https://travis-ci.org/ligato/crd-example.svg?branch=master)](https://travis-ci.org/ligato/crd-example)
[![GitHub license](https://img.shields.io/badge/license-Apache%20license%202.0-blue.svg)](https://github.com/ligato/crd-example/blob/master/LICENSE)

Custom Resource Definition Example
==================================

This file shows an example of creating Custom Resource Definitions (CRDs) to
integrate inside of a Ligato plugin. At their essence, CRDs are an extension
of the Kubernetes API. [This](1) page describes CRDs in detail, while [this](2)
page describes extending Kubernetes with CRDs.

To programmatically extend the Kubernetes API with CRDs, we will show an
example comprising three distinct parts:

* A protobuf file, in this case `api/crdexample.proto`, which defines the API
  resources you want to expose.
* A types.go file, in this case `pkg/apis/crdexample.io/v1/types.go`, which
  defines the CRD API structures and uses the generated Go code from the
  protobuf file as it's schema.
* The Kubernetes code generators to generate all the structural pieces
  necessary to build an application to handle the CRDs. This includes listers,
  watchers, and informers.

This combination is powerful and shows how extending the Kubernetes API is
an effective way to introduce new resources types almost natively.

Running the Example
-------------------

The example requires a working Kubernetes install. A Minikube install will work
just fine.

To build and install the example using go v1.10.x:

```
make all
```

To build and install the example using go v1.11.x:

```
GO111MODULE=on make all
```

To run the example:

```
kubectl label --overwrite nodes docker-for-desktop app=crdservice-node
kubectl create -f conf/crd-daemonset.yaml
```

If you want to stop the example:

```
kubectl delete -f conf/crd-daemonset.yaml
```

Updating the Generated Code
---------------------------

To update the generated code:

```
PATH=$GOPATH/bin:$PATH go generate ./...
CODEGEN_PKG=./vendor/k8s.io/code-generator ./scripts/update-codegen.sh
```

[1]: https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/
[2]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/
