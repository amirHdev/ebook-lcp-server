#!/usr/bin/env sh
set -eu

KUSTOMIZE="${KUSTOMIZE:-kustomize}"

render_kustomization() {
  if command -v "${KUSTOMIZE}" >/dev/null 2>&1; then
    "${KUSTOMIZE}" build "$1"
    return
  fi

  if command -v kubectl >/dev/null 2>&1; then
    kubectl kustomize "$1"
    return
  fi

  echo "kustomize or kubectl is required to render Kubernetes resources" >&2
  exit 1
}

go test ./...
go vet ./...
go build -buildvcs=false ./...

(
  cd frontend
  if [ ! -d node_modules ]; then
    npm ci
  fi
  npm run build
)

yamllint -c .yamllint.yml \
  deploy/k8s \
  deploy/overlays \
  deploy/argocd \
  docker-compose.yml \
  deploy/monitoring/prometheus.yml \
  .github/workflows/go.yml \
  .gitlab-ci.yml

for env in dev staging prod; do
  render_kustomization "deploy/overlays/${env}" >/dev/null
done

render_kustomization deploy/argocd/apps >/dev/null
