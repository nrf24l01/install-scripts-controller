# install-scripts-controller

Helm chart for [Install Scripts Controller](https://github.com/nrf24l01/install-scripts-controller): a web app that hosts shell install scripts behind a password-protected UI, exposing a single `curl | bash` command per script.

## Install

From the OCI registry (published by CI):

```sh
helm upgrade --install install-scripts-controller \
  oci://ghcr.io/nrf24l01/charts/install-scripts-controller \
  --version 0.2.1 \
  --namespace install-scripts-controller \
  --create-namespace \
  --set secrets.password='your-password'
```

From the local chart (development):

```sh
helm upgrade --install install-scripts-controller \
  ./helm/install-scripts-controller \
  --namespace install-scripts-controller \
  --create-namespace \
  --set secrets.password='your-password'
```

Check pods and get the URL:

```sh
kubectl get pods -l app.kubernetes.io/instance=install-scripts-controller
kubectl port-forward -n install-scripts-controller svc/install-scripts-controller 1325:80
# open http://localhost:1325
```

## Parameters

| Value                  | Default                                              | Description                                |
| ---------------------- | ---------------------------------------------------- | ------------------------------------------ |
| `replicaCount`         | `1`                                                  | Pod replicas (SQLite is single-writer, keep `1`) |
| `image.repository`     | `ghcr.io/nrf24l01/install-scripts-controller/install-scripts-controller` | App image      |
| `image.tag`            | `latest`                                             | App image tag                              |
| `image.pullPolicy`     | `IfNotPresent`                                       | Image pull policy                          |
| `secrets.password`     | `""` (required)                                      | UI sign-in password                        |
| `secrets.create`       | `true`                                               | Create the config Secret; set `false` + `secrets.name` to reuse an existing one |
| `secrets.name`         | `""`                                                 | Secret name override                       |
| `config.installKeyTTL` | `24h`                                                | How long each random install key stays valid |
| `config.publicUrl`     | `""`                                                 | Public base URL for install links          |
| `config.serverAddr`    | `:8080`                                              | Address the app listens on inside the pod  |
| `config.databasePath`  | `/data/app.db`                                       | SQLite DB path (keep under `/data`)        |
| `service.type`         | `ClusterIP`                                          | Service type                               |
| `service.port`         | `80`                                                 | Service port                               |
| `service.targetPort`   | `http`                                               | Container port to route to                  |
| `service.nodePort`     | `""`                                                 | Fixed nodePort when `type=NodePort` (auto-assigned if empty) |
| `persistence.enabled`  | `true`                                               | Persistent volume for the database         |
| `persistence.size`     | `1Gi`                                                | PVC size                                   |
| `persistence.storageClass` | `""`                                             | PVC storage class                          |
| `persistence.accessMode` | `ReadWriteOnce`                                    | PVC access mode                            |
| `ingress.enabled`      | `false`                                              | Expose via Ingress                         |
| `ingress.className`    | `""`                                                 | Ingress class name                         |
| `ingress.host`         | `""`                                                 | Ingress host (e.g. `scripts.example.com`)  |
| `ingress.tls`          | `[]`                                                 | Ingress TLS entries (`secretName` + `hosts`) |
| `resources`            | `requests: 50m/64Mi`, `limits: 500m/256Mi`           | Container resources                        |
| `podSecurityContext`   | `fsGroup: 65532`                                     | Pod security context                       |
| `securityContext`      | non-root `65532`, RO filesystem                      | Container security context                 |

## Configuration

The app reads a single `config.yml`. The chart renders it from the values above
and stores it in a Kubernetes Secret (`<release>-config`), mounted read-only at
`/app/config.yml`. `site.password`, `site.install_key_ttl`, `site.public_url`,
`server.addr` and `database.path` all come from chart values.

Since the Secret is mounted with a `subPath`, the deployment carries a
`checksum/secret` annotation, so changing `secrets.password` (or other config)
triggers a rolling restart automatically.

## Uninstall

```sh
helm uninstall install-scripts-controller -n install-scripts-controller
```

To also delete the database volume, remove the PVC:

```sh
kubectl delete pvc install-scripts-controller-data -n install-scripts-controller
```
