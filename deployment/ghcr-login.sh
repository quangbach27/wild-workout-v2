#!/bin/bash
# Logs docker into ghcr.io on the instance so private images can be pulled. Shipped and run by
# .github/workflows/deploy.yml. The token (a PAT with read:packages only) is the GHCR_TOKEN key
# of the JSON secret "<NAME_PREFIX>/ghcr".
# Expects: REGION, NAME_PREFIX, GHCR_USER in the environment.
set -euo pipefail

aws secretsmanager get-secret-value --region "$REGION" --secret-id "$NAME_PREFIX/ghcr" \
  --query SecretString --output text \
  | python3 -c 'import json, sys; sys.stdout.write(json.load(sys.stdin)["GHCR_TOKEN"])' \
  | docker login ghcr.io -u "$GHCR_USER" --password-stdin
