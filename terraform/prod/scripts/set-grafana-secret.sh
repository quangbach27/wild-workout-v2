#!/usr/bin/env bash
# Stores the Grafana Cloud Loki credentials in the <NAME_PREFIX>/grafana Secrets Manager secret.
# Terraform only creates the empty secret (values never enter state), so run this once after `terraform apply`.
#
# Inputs (env vars, prompted when missing): LOKI_URL, LOKI_USER, LOKI_TOKEN (prompted silently).
# Optional: NAME_PREFIX (default: `terraform output -raw name_prefix`), AWS_REGION (default: ap-southeast-2).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export AWS_REGION="${AWS_REGION:-ap-southeast-2}"

NAME_PREFIX="${NAME_PREFIX:-$(terraform -chdir="$SCRIPT_DIR/.." output -raw name_prefix)}"
SECRET_ID="$NAME_PREFIX/grafana"

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
printf 'LOKI_URL=%s\nLOKI_USER=%s\nLOKI_TOKEN=%s\n' "$LOKI_URL" "$LOKI_USER" "$LOKI_TOKEN" > "$SECRET_FILE"

aws secretsmanager put-secret-value --secret-id "$SECRET_ID" --secret-string "file://$SECRET_FILE" \
  --query VersionId --output text

echo "Updated $SECRET_ID. Run the deploy workflow (or 'Rotate Grafana token') so Alloy picks up the new value."
