# LCP Server

Lightweight License Content Protection (LCP) server that exposes GraphQL APIs for uploading encrypted publications and issuing licenses. This repository now includes the DevOps assets needed to run the service in production with Docker, Kubernetes, GitLab CI/CD, and ArgoCD.

## Features
- GraphQL endpoint at `/graphql` for managing publications and licenses.
- Pluggable encryption interface (default file copy encrypter for development) to integrate with a full LCP DRM backend.
- In-memory repositories that keep the service stateless for easy containerization.
- Download endpoint at `/publications/{id}/content` for clients to retrieve encrypted assets using the URLs returned on licenses.
- Deployment assets for Docker, Kubernetes (with Kustomize), and ArgoCD GitOps flows.
- GitLab pipeline that lints, tests, builds, and deploys the container image.

## Configuration
Set the following environment variables (see `.env.example` for defaults):

- `DB_DSN`: Database connection string (used by adapters that expect persistent storage).
- `LCP_PROFILE`: `basic` or `production`.
- `LCP_CERTIFICATE` / `LCP_PRIVATE_KEY`: Paths to DRM keys/certificates.
- `LCP_STORAGE_MODE`: `fs` (default) or `s3`.
- `LCP_STORAGE_FS_DIR`: Target directory for encrypted assets.
- `LCP_S3_REGION`, `LCP_S3_BUCKET`, `LCP_S3_ACCESS_KEY`, `LCP_S3_SECRET_KEY`: S3 storage settings when `LCP_STORAGE_MODE=s3`.
- `JWT_SECRET`: Secret for future JWT-protected endpoints.
- `SERVER_PORT`: Listen address (defaults to `:8080`).
- `PUBLIC_BASE_URL`: Public base URL used to generate download links (defaults to `http://localhost:PORT`).

## Local development
```bash
# Install dependencies and start the API
cp .env.example .env  # then edit values
export $(grep -v '^#' .env | xargs)
go run ./cmd/server
```
The GraphQL playground will be available at `http://localhost:8080/graphql`.

### GraphQL upload notes
The `uploadPublication` mutation expects the `file` argument as a **base64-encoded string**. Example variables:
```json
{
  "title": "My Book",
  "file": "<base64 encoded content>"
}
```

## Docker
```bash
docker build -t lcp-server:local .
docker run -p 8080:8080 --env-file .env lcp-server:local
```
The multi-stage Dockerfile compiles the Go binary and ships a minimal distroless runtime image suitable for production.

## Kubernetes
Apply the manifests with Kustomize:
```bash
kubectl apply -k deploy/k8s
```
The deployment uses two replicas, resource requests/limits, and a writable volume for encrypted assets (`/var/lib/lcp/storage`). Update `deploy/k8s/deployment.yaml` with your container registry image and storage class as needed.

## ArgoCD
The ArgoCD application manifest at `deploy/argocd/application.yaml` continuously syncs the Kustomize manifests from this repo:
```bash
kubectl apply -n argocd -f deploy/argocd/application.yaml
```
Namespace `lcp` will be created automatically, and changes merged to the repo will propagate to the cluster.

## GitLab CI/CD
`.gitlab-ci.yml` defines four stages:
1. **lint**: runs `gofmt -l` and `go vet`.
2. **test**: executes `go test ./...`.
3. **build**: builds the Docker image and optionally pushes it when registry credentials exist.
4. **deploy**: applies the Kustomize overlays to the cluster (expects `kubectl` credentials in CI variables).

Set `CI_REGISTRY`, `CI_REGISTRY_USER`, `CI_REGISTRY_PASSWORD`, and `KUBECONFIG` (or in-cluster service account variables) in GitLab to enable full automation.

## Repository layout
- `cmd/server`: HTTP server wiring, GraphQL handler, and LCP use cases.
- `internal/usecase/lcp`: Business logic for publications and licenses.
- `internal/adapter/graphql`: GraphQL schema and resolvers.
- `deploy/k8s`: Production manifests with Kustomize.
- `deploy/argocd`: GitOps application definition.
- `.gitlab-ci.yml`: Pipeline definition for GitLab.

