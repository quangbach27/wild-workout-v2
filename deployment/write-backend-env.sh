#!/bin/bash
# Builds /opt/app/backend.env on the instance, shipped and run by .github/workflows/deploy.yml.
# Sources: the "<NAME_PREFIX>/backend" secret (JSON of env vars), the RDS-managed secret
# (JSON username/password) and the DB_* values passed in by the workflow.
# Expects: REGION, NAME_PREFIX, DB_SECRET_STORE, DB_HOST, DB_NAME in the environment.
set -euo pipefail
umask 077

fetch() {
  aws secretsmanager get-secret-value --region "$REGION" --secret-id "$1" \
    --query SecretString --output text
}

# Prints KEY=value lines from a JSON object on stdin, quoted so docker compose reads them literally.
to_env() {
  python3 -c '
import json, sys

def quote(v):
    v = str(v)
    if "\x27" not in v:
        return "\x27" + v + "\x27"
    return "\"" + v.replace("\\", "\\\\").replace("\"", "\\\"").replace("$", "\\$") + "\""

# GHCR_TOKEN and LOKI_* share this secret but are for the host (docker login, Alloy), not the backend.
for k, v in json.load(sys.stdin).items():
    if k == "GHCR_TOKEN" or k.startswith("LOKI_"):
        continue
    print(f"{k}={quote(v)}")
'
}

fetch_db() {
  fetch "$DB_SECRET_STORE" | python3 -c '
import json, sys

d = json.load(sys.stdin)
print(json.dumps({"DB_USERNAME": d["username"], "DB_PASSWORD": d["password"]}))
' | to_env
}

tmp=$(mktemp backend.env.XXXXXX)
{
  fetch "$NAME_PREFIX/backend" | to_env
  fetch_db
  echo "DB_HOST=$DB_HOST"
  echo "DB_NAME=$DB_NAME"
  echo "DB_PORT=5432"
  echo "DB_SSLMODE=require"
} > "$tmp"
mv "$tmp" backend.env
