# One JSON secret holding the backend's runtime config (expected keys: local.backend_secret_keys in locals.tf).
# Terraform creates the container only; set the value by hand so it stays out of state:
#   aws secretsmanager put-secret-value --secret-id wild-workout-<workspace>/backend --secret-string file://backend.json
module "secrets" {
  source = "../modules/ssm"

  name         = local.name
  secret_names = ["backend", "grafana", "ghcr"]
}

# Set it with `make set-grafana-secret`. "grafana" holds the Grafana Cloud Loki credentials as plain env lines (read by .github/workflows/deploy.yml
# into /opt/app/alloy.env for the Alloy log shipper):
#   LOKI_URL=https://logs-prod-XXX.grafana.net/loki/api/v1/push
#   LOKI_USER=<numeric user id>
#   LOKI_TOKEN=<access token with logs:write>
#
# "ghcr" holds a GitHub PAT (read:packages only) under the key GHCR_TOKEN; .github/workflows/deploy.yml logs docker into ghcr.io with it:
#   aws secretsmanager put-secret-value --secret-id wild-workout-<workspace>/ghcr --secret-string '{"GHCR_TOKEN":"ghp_xxx"}'
