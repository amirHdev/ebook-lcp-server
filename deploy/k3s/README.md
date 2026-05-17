# K3s Bootstrap

This project is intended to run on self-hosted K3s.

## Cluster Shape

Recommended minimum for production-like use:

- 3 server nodes for HA control plane
- 2 worker nodes for workloads
- Persistent storage available on each node or via a shared storage backend

K3s already includes the pieces this project needs by default:

- container runtime
- Flannel CNI
- CoreDNS
- Traefik ingress controller
- ServiceLB
- local-path storage provisioner

## Install

Use the install script from this repo or the upstream K3s installer.

First place the cluster config on the server:

```bash
scripts/infra/install-k3s-config.sh
```

Example single-node install:

```bash
curl -sfL https://get.k3s.io | sh -
```

Example HA server with embedded etcd:

```bash
curl -sfL https://get.k3s.io | K3S_TOKEN=replace-with-token sh -s - server \
  --cluster-init \
  --tls-san yourdomain.com \
  --tls-san status.yourdomain.com
```

Join a worker:

```bash
curl -sfL https://get.k3s.io | K3S_TOKEN=replace-with-token sh -s - agent \
  --server https://server-ip:6443
```

## Add-ons

Install these after K3s:

- ArgoCD
- cert-manager
- Prometheus Operator stack or equivalent
- optional: external secrets or Vault

The repo includes ArgoCD Applications for `dev`, `staging`, and `prod`.
