# Infrastructure Runbook

This project is prepared for self-hosted K3s with ArgoCD GitOps deployment across `dev`, `staging`, and `prod`.

## Required Cluster Add-ons

Install these before syncing the LCP applications:

- ArgoCD
- cert-manager with a `letsencrypt-prod` ClusterIssuer
- kube-prometheus-stack or another Prometheus Operator installation
- Optional: External Secrets or Vault for production-grade secret delivery

K3s already provides the core platform pieces:

- Traefik ingress controller
- local-path storage provisioner
- ServiceLB
- CoreDNS
- Flannel CNI
- container runtime

## Repository Layout

- `deploy/k8s`: reusable Kubernetes base.
- `deploy/overlays/dev`: low-resource development environment.
- `deploy/overlays/staging`: production-like staging environment.
- `deploy/overlays/prod`: production environment.
- `deploy/registry`: self-hosted image registry.
- `deploy/k3s`: bootstrap notes and K3s cluster config.
- `deploy/argocd/root-application.yaml`: app-of-apps entrypoint.
- `deploy/argocd/apps`: ArgoCD project and environment applications.

## Bootstrap

1. Install the in-cluster registry from `deploy/registry`.
2. Push backend and frontend images to `registry.yourdomain.com:5000`.
3. Replace hosts in overlay ingress patches:
   - `dev.yourdomain.com`
   - `staging.yourdomain.com`
   - `yourdomain.com`
   - `status.yourdomain.com`
   - `argocd.yourdomain.com`
   - `registry.yourdomain.com`
4. Create real secrets for each namespace:

```bash
cp deploy/secrets/lcp-secrets.example.env /tmp/lcp-dev.env
scripts/infra/create-secret.sh lcp-dev /tmp/lcp-dev.env
scripts/infra/create-secret.sh lcp-staging /tmp/lcp-staging.env
scripts/infra/create-secret.sh lcp-prod /tmp/lcp-prod.env
```

5. Install ArgoCD root application:

```bash
kubectl apply -n argocd -f deploy/argocd/root-application.yaml
```

ArgoCD will create and sync the dev, staging, and prod applications.
The K3s cluster itself is expected to be installed separately using the scripts in `scripts/infra` or the upstream K3s installer.

Useful scripts:

- `scripts/infra/bootstrap-k3s-stack.sh`
- `scripts/infra/install-k3s-config.sh`
- `scripts/infra/install-k3s-registries.sh`
- `scripts/infra/install-k3s-server.sh`
- `scripts/infra/install-k3s-agent.sh`
- `scripts/infra/build-and-push-images.sh`

## Validation

Render all manifests:

```bash
KUSTOMIZE=/root/go/bin/kustomize scripts/infra/render-all.sh
```

Run local validation:

```bash
KUSTOMIZE=/root/go/bin/kustomize scripts/infra/validate.sh
```

Check deployment:

```bash
kubectl -n lcp-prod rollout status deploy/lcp-server
kubectl -n lcp-prod rollout status deploy/lcp-admin-ui
kubectl -n lcp-prod rollout status deploy/lcp-core
kubectl -n lcp-prod rollout status deploy/lsd-core
kubectl -n lcp-prod get hpa,ingress,cronjob,servicemonitor
```

## Production Notes

The included PostgreSQL StatefulSet is acceptable for a fully self-hosted deployment. If you want higher durability, swap it for an external PostgreSQL cluster that you still operate yourself.

The included backup CronJob archives app metadata and publication storage. For production, mount durable block storage or object storage instead of `emptyDir`.
