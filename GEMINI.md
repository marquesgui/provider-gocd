# GEMINI.md

## Project Overview

This project is a [Crossplane](https://crossplane.io/) provider for [GoCD](https://www.gocd.org/). It allows you to manage GoCD resources, such as `PipelineConfig`, `Role`, and `AuthorizationConfiguration`, using the Crossplane declarative API.

The provider is built using the Crossplane provider-runtime and follows the standard controller pattern. It defines Custom Resource Definitions (CRDs) for the GoCD resources it manages and includes controllers that reconcile the state of these resources with the GoCD API.

This repository serves as a template for creating new Crossplane providers. It includes scaffolding for new types and controllers, as well as a `Makefile` with common development tasks.

## Building and Running

The `Makefile` provides a set of targets for building, running, and testing the provider.

### Building

To build the provider, run:

```shell
make build
```

This will compile the Go code and create the provider binary in the `_output` directory.

### Running Locally

To run the provider locally against a Kubernetes cluster, you can use the `run` target:

```shell
make run
```

This will run the provider out-of-cluster and connect to the Kubernetes cluster configured in your `kubeconfig` file.

For a more integrated development experience, you can use the `dev` target:

```shell
make dev
```

This will create a local `kind` cluster, install the provider's CRDs, and run the provider in-cluster.

To clean up the development environment, run:

```shell
make dev-clean
```

### Testing

To run the integration tests, use the `test-integration` target:

```shell
make test-integration
```

This will run the integration tests against a local `kind` cluster.

To run linters, code generation, and tests, use the `reviewable` target:

```shell
make reviewable
```

## Development Conventions

### Adding New Types

To add a new managed resource type, you can use the `provider.addtype` make target:

```shell
export provider_name=MyProvider # Camel case, e.g. GitHub
export group=sample # lower case e.g. core, cache, database, storage, etc.
export type=MyType # Camel casee.g. Bucket, Database, CacheCluster, etc.
make provider.addtype provider=${provider_name} group=${group} kind=${type}
```

This will generate the necessary API and controller files for the new type.

### Code Style

The project follows the standard Go code style. You can use `make reviewable` to run the linters and ensure your code conforms to the style guidelines.

### Commits

Commit messages should follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.
