# memcached-operator

// TODO(user): Add simple overview of use/purpose

## Description

// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

### Prerequisites

- cargo version 1.85.1
- docker version 27.5.0+
- kubectl version v1.32.1+.
- Access to a Kubernetes v1.25.3+ cluster.

### To Run locally

**Build your operator:**

```sh
make build
```

**Run your operator:**

```sh
make run
```

### To Deploy on the cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make image-build image-push IMG=<some-registry>/memcached-operator:tag
```

> **NOTE:** This image ought to be published in the personal registry you specified.
> And it is required to have access to pull the image from the working environment.
> Make sure you have the proper permission to the registry if the above commands don’t work.

**Generate the CRDs:**

```sh
make generate-crds
```

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the operator to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/memcached-operator:tag
```

> **IMPORTANT**: You will face API access errors, as this script only creates the deployment.
> You will need to create the required role bindings for your deployment for now.

**Create instances of your solution**
You can apply your example CRs:

```sh
kubectl apply -k path/to/your/samples/
```

> **IMPORTANT**: Ensure that the samples has default values to test it out.

### To Uninstall

**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k path/to/your/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

// TODO(user): Add a license
