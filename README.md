# istio-upgrade-consumer

Consumer for Istio Upgrade tooling to receive events related to Istio upgrade

![Version: 1.0.1](https://img.shields.io/badge/Version-1.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0.1](https://img.shields.io/badge/AppVersion-1.0.1-informational?style=flat-square) [![made with Go](https://img.shields.io/badge/made%20with-Go-brightgreen)](http://golang.org) [![Github main branch build](https://img.shields.io/github/workflow/status/gopaytech/istio-upgrade-consumer/Main)](https://github.com/gopaytech/istio-upgrade-consumer/actions/workflows/main.yml) [![GitHub issues](https://img.shields.io/github/issues/gopaytech/istio-upgrade-consumer)](https://github.com/gopaytech/istio-upgrade-consumer/issues) [![GitHub pull requests](https://img.shields.io/github/issues-pr/gopaytech/istio-upgrade-consumer)](https://github.com/gopaytech/istio-upgrade-consumer/pulls)[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/istio-upgrade-consumer)](https://artifacthub.io/packages/search?repo=istio-upgrade-consumer)

## Installing

To install the chart with the release name `my-release`:

```console
helm repo add istio-upgrade-consumer https://gopaytech.github.io/istio-upgrade-consumer/charts/releases/
helm install my-istio-upgrade-consumer istio-upgrade-consumer/istio-upgrade-consumer --values values.yaml
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| configuration.clusterEnvironment | string | `"production"` |  |
| configuration.clusterName | string | `"my-cluster"` |  |
| configuration.nonProductionWaitingDay | string | `"7"` |  |
| configuration.productionWaitingDay | string | `"28"` |  |
| configuration.receiverHTTPPort | string | `"8080"` |  |
| configuration.receiverMode | string | `"http"` |  |
| configuration.storageConfigMapName | string | `"istio-auto-upgrade-config"` |  |
| configuration.storageConfigMapNamespace | string | `"istio-system"` |  |
| configuration.storageMode | string | `"configmap"` |  |
| configuration.timeFormat | string | `"2006-01-02"` |  |
| configuration.timeLocation | string | `"Asia/Jakarta"` |  |
| deployment.image | string | `"ghcr.io/gopaytech/istio-upgrade-consumer"` |  |
| deployment.replicas | int | `2` |  |
| deployment.tag | string | `"v1.0.1"` |  |
| podLabels | object | `{}` |  |
| resources.limits.cpu | string | `"1024m"` |  |
| resources.limits.memory | string | `"1024Mi"` |  |
| resources.requests.cpu | string | `"256m"` |  |
| resources.requests.memory | string | `"256Mi"` |  |
| serviceAccount.automountServiceAccountToken | bool | `true` |  |
| serviceAccount.imagePullSecrets | list | `[]` |  |

