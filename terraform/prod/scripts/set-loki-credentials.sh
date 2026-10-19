#!/usr/bin/env bash
# Merges the Grafana Cloud Loki credentials (LOKI_URL, LOKI_USER, LOKI_TOKEN) into the JSON of the
# <NAME_PREFIX>/backend Secrets Manager secret, leaving its other keys untouched. Needs aws and jq.
#
# Inputs (env vars, prompted when missing): LOKI_URL, LOKI_USER, LOKI_TOKEN (prompted silently).
# Optional: NAME_PREFIX (default: `terraform output -raw name_prefix`), AWS_REGION (default: ap-southeast-2).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export AWS_REGION="${AWS_REGION:-ap-southeast-2}"

NAME_PREFIX="${NAME_PREFIX:-$(terraform -chdir="$SCRIPT_DIR/.." output -raw name_prefix)}"
SECRET_ID="$NAME_PREFIX/backend"

if [ -z "${LOKI_URL:-}" ]; then
  read -r -p "LOKI_URL (https://logs-prod-XXX.grafana.net/loki/api/v1/push): " LOKI_URL
fi
if [ -z "${LOKI_USER:-}" ]; then
  read -r -p "LOKI_USER (numeric user id): " LOKI_USER
fi
if [ -z "${LOKI_TOKEN:-}" ]; then
  read -r -s -p "LOKI_TOKEN (input hidden): " LOKI_TOKEN
  echo
fi

for var in LOKI_URL LOKI_USER LOKI_TOKEN; do
  [ -n "${!var}" ] || { echo "$var is empty" >&2; exit 1; }
done

SECRET_FILE="$(mktemp)"
trap 'rm -f "$SECRET_FILE"' EXIT
chmod 600 "$SECRET_FILE"

CURRENT="$(aws secretsmanager get-secret-value --secret-id "$SECRET_ID" --query SecretString --output text)"
jq -e 'type == "object"' <<<"$CURRENT" >/dev/null || { echo "$SECRET_ID is not a JSON object" >&2; exit 1; }
jq --arg url "$LOKI_URL" --arg user "$LOKI_USER" --arg token "$LOKI_TOKEN" \
  '. + {LOKI_URL: $url, LOKI_USER: $user, LOKI_TOKEN: $token}' <<<"$CURRENT" > "$SECRET_FILE"

aws secretsmanager put-secret-value --secret-id "$SECRET_ID" --secret-string "file://$SECRET_FILE" \
  --query VersionId --output text

echo "Updated $SECRET_ID. Restart alloy on the instance (systemctl restart alloy) so it picks up the new value."
