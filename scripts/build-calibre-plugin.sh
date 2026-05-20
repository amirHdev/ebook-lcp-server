#!/bin/sh
set -eu

out="${1:-dist/lcp-send-calibre-plugin.zip}"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

mkdir -p "$(dirname "$out")" "$tmp/lcp_send"
cp integrations/calibre_plugin/__init__.py "$tmp/__init__.py"
cp integrations/calibre_plugin/config.py "$tmp/lcp_send/config.py"
cp integrations/calibre_plugin/ui.py "$tmp/lcp_send/ui.py"

(cd "$tmp" && zip -qr "$OLDPWD/$out" __init__.py lcp_send)
echo "$out"
