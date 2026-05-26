#!/bin/sh
set -eu

out_dir="${1:-artifacts/e2e-readiness}"
mkdir -p "$out_dir"

services="postgres redis minio lcp-core lsd-core lcp-server"
project_name="${COMPOSE_PROJECT_NAME:-lcp-e2e}"
compose_cmd="docker compose -p $project_name"
keep_up="${KEEP_STACK_UP:-0}"

cleanup() {
  if [ "$keep_up" = "1" ]; then
    return
  fi
  $compose_cmd logs --no-color > "$out_dir/compose.log" 2>&1 || true
  $compose_cmd down -v --remove-orphans >/dev/null 2>&1 || true
}

trap cleanup EXIT INT TERM

$compose_cmd down -v --remove-orphans >/dev/null 2>&1 || true

attempt=1
while :; do
  if $compose_cmd up --build -d $services; then
    break
  fi
  if [ "$attempt" -ge 3 ]; then
    exit 1
  fi
  attempt=$((attempt + 1))
  sleep 5
done

for _ in $(seq 1 90); do
  if curl -fsS http://localhost:8080/readyz >/dev/null 2>&1; then
    break
  fi
  sleep 2
done

curl -fsS http://localhost:8080/readyz > "$out_dir/readyz.json"

sh scripts/acceptance-smoke.sh "$out_dir/acceptance-report.json" >/dev/null
sh scripts/demo-local.sh > "$out_dir/demo.txt"
sh scripts/certification-smoke.sh "$out_dir/certification-report.json" >/dev/null
sh scripts/generate-certification-packet.sh "$out_dir/certification-packet" >/dev/null

$compose_cmd ps > "$out_dir/compose-ps.txt"
$compose_cmd logs --no-color > "$out_dir/compose.log"

if command -v docker >/dev/null 2>&1; then
  postgres_id="$($compose_cmd ps -q postgres 2>/dev/null || true)"
  if [ -n "$postgres_id" ]; then
    docker exec "$postgres_id" psql -U lcp -d lcp -Atqc \
      "select table_name from information_schema.tables where table_schema='public' order by table_name;" \
      > "$out_dir/postgres-tables.txt"
  fi
fi

echo "$out_dir"
