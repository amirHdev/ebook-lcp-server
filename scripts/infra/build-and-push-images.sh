#!/usr/bin/env sh
set -eu

REGISTRY="${REGISTRY:-your-registry.example.com}"
TAG="${TAG:-latest}"
BACKEND_IMAGE="${REGISTRY}/lcp-server"
FRONTEND_IMAGE="${REGISTRY}/lcp-admin-ui"

export BUILDAH_ISOLATION="${BUILDAH_ISOLATION:-chroot}"

buildah bud -t "${BACKEND_IMAGE}:${TAG}" .
buildah bud -t "${FRONTEND_IMAGE}:${TAG}" frontend
buildah push --tls-verify=false "${BACKEND_IMAGE}:${TAG}" "docker://${BACKEND_IMAGE}:${TAG}"
buildah push --tls-verify=false "${FRONTEND_IMAGE}:${TAG}" "docker://${FRONTEND_IMAGE}:${TAG}"
