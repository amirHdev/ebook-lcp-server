#!/bin/sh
set -eu

secret_file="/srv/readium/source/readium-core-secret.yaml"
config_dir="/srv/readium/config"

mkdir -p "$config_dir"

awk '
  /^  htpasswd:/ {
    sub(/^  htpasswd: /, "")
    print
  }
' "$secret_file" > "$config_dir/htpasswd"

awk '
  /^  cert-edrlab-test.pem: \|/ { in_block = 1; next }
  /^  privkey-edrlab-test.pem: \|/ { in_block = 0 }
  in_block {
    sub(/^    /, "")
    print
  }
' "$secret_file" > "$config_dir/cert-edrlab-test.pem"

awk '
  /^  privkey-edrlab-test.pem: \|/ { in_block = 1; next }
  in_block {
    sub(/^    /, "")
    print
  }
' "$secret_file" > "$config_dir/privkey-edrlab-test.pem"
