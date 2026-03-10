# Istio Helm Chart

> [Istio](https://istio.io/) extends Kubernetes to establish a programmable, application-aware network. Working with both Kubernetes and
traditional workloads, Istio brings standard, universal traffic management, telemetry, and security to complex deployments.

## Introduction

This Helm chart bootstraps an [Istio](https://istio.io/) deployment on a [Kubernetes](http://kubernetes.io) cluster using the
[Helm](https://helm.sh) package manager.

## Quick Start

```console
helm install my-release oci://dp.apps.rancher.io/charts/istio \
    --set global.imagePullSecrets={application-collection} \
    -n istio-system --create-namespace
```

### Prerequisites

* Helm 3.8.0 or later.
* Kubernetes 1.24 or later.
* PV provisioner support in the underlying infrastructure.

## Install Chart

To install the Helm chart with the release name *my-release*:

```console
helm install my-release oci://dp.apps.rancher.io/charts/istio \
    --set global.imagePullSecrets={application-collection} \
    -n istio-system --create-namespace
```

This deploys the application to the Kubernetes cluster using the default configuration provided by the Helm chart, in the *istio-system*
namespace.

> NOTE: Follow [these steps](https://docs.apps.rancher.io/get-started/authentication/#kubernetes) to create the image pull secret named
> *application-collection*, if you don't have it already.

## Uninstall Chart

To uninstall the Helm chart with the release name *my-release*:

```console
helm uninstall my-release -n istio-system
```

Due to Helm's [design](https://helm.sh/docs/chart_best_practices/custom_resource_definitions/#some-caveats-and-explanations),
CustomResourceDefinitions (CRDs) are not removed when uninstalling the related chart. To do so, you will need to explicitly remove them:

```console
kubectl delete $(kubectl get CustomResourceDefinition -l='release=istio' -o name -A) -n istio-system
```

## Configuration

To view support configuration options and documentation, run:

```console
helm show values oci://dp.apps.rancher.io/charts/istio
```

### Profiles

The Istio Helm chart has a concept of a `profile`, which is a bundled collection of value presets.
These can be set with `--set <subchart>.profile=<profile>`.
For example, the `demo` profile offers a preset configuration to try out Istio in a test environment,
with additional features enabled and lowered resource requirements.

For consistency, the same profiles are used across each subchart, even if they do not impact a given chart.

Explicitly set values have highest priority, then profile settings, then chart defaults.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| base.* | string | `nil` | base sub-chart configurations (refer to the sub-chart to obtain the list of supported values) |
| base.enabled | bool | `true` | Whether to enable deploy the base sub-chart |
| base.profile | string | `""` | [Configuration profile](https://istio.io/latest/docs/setup/additional-setup/config-profiles/) to be enabled for the base sub-chart |
| cni.* | string | `nil` | cni sub-chart configurations (refer to the sub-chart to obtain the list of supported values) |
| cni.enabled | bool | `false` | Whether to enable deploy the cni sub-chart |
| cni.profile | string | `""` | [Configuration profile](https://istio.io/latest/docs/setup/additional-setup/config-profiles/) to be enabled for the base sub-chart |
| gateway.* | string | `nil` | gateway sub-chart configurations (refer to the sub-chart to obtain the list of supported values) |
| gateway.enabled | bool | `false` |  |
| gateway.profile | string | `""` | [Configuration profile](https://istio.io/latest/docs/setup/additional-setup/config-profiles/) to be enabled for the base sub-chart |
| global.* | string | `nil` | Global override values for all Istio sub-charts (refer to each sub-chart to obtain the list of all supported values) |
| global.imagePullSecrets | list | `[]` | Global override for container image registry pull secrets |
| global.imageRegistry | string | `""` | Global override for container image registry |
| global.platform | string | `""` | Platform where Istio is deployed. Possible values are: "openshift", "gcp". An empty value means it is a vanilla Kubernetes distribution, therefore no special treatment will be considered. |
| istiod.* | string | `nil` | istiod sub-chart configurations (refer to the sub-chart to obtain the list of supported values) |
| istiod.enabled | bool | `true` | Whether to deploy the istiod sub-chart |
| istiod.profile | string | `""` | [Configuration profile](https://istio.io/latest/docs/setup/additional-setup/config-profiles/) to be enabled for the base sub-chart |
| ztunnel.* | string | `nil` | ztunnel sub-chart configurations (refer to the sub-chart to obtain the list of supported values) |
| ztunnel.enabled | bool | `false` | Whether to deploy the ztunnel sub-chart |
| ztunnel.profile | string | `""` | [Configuration profile](https://istio.io/latest/docs/setup/additional-setup/config-profiles/) to be enabled for the base sub-chart |
