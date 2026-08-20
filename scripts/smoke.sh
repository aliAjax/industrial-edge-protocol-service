#!/bin/sh
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
PORT=${EDGE_SMOKE_PORT:-18099}
DATA_DIR=$(mktemp -d)
EDGE_HTTP_ADDR=:$PORT EDGE_DATA_DIR="$DATA_DIR" "$ROOT/edge-api" >/tmp/edge-smoke.log 2>&1 &
PID=$!
trap 'kill "$PID" 2>/dev/null || true; rm -rf "$DATA_DIR"' EXIT
for i in $(seq 1 30); do curl -fsS "http://127.0.0.1:$PORT/healthz" >/dev/null && break; sleep .1; done
curl -fsS -X POST "http://127.0.0.1:$PORT/api/v1/sites" -H 'content-type: application/json' -d '{"name":"demo","location":"lab"}'
curl -fsS "http://127.0.0.1:$PORT/api/v1/sites"
