# One JSON secret holding the backend's runtime config (expected keys: local.backend_secret_keys in locals.tf).
# Terraform creates the container only; set the value by hand so it stays out of state:
#   aws secretsmanager put-secret-value --secret-id wild-workout-<workspace>/backend --secret-string file://backend.json
module "secrets" {
  source = "../modules/ssm"

  name         = local.name
  secret_names = ["backend", "grafana", "ghcr"]
}

# Set it with `make set-grafana-secret`. "grafana" holds the Grafana Cloud Loki credentials, read by Alloy on every start (see modules/ec2/scripts/40-alloy.sh.tftpl).
# Either JSON key/value (the console's Key/value editor) or plain env lines work; the env-line form is:
#   LOKI_URL=https://logs-prod-XXX.grafana.net/loki/api/v1/push
#   LOKI_USER=<numeric user id>
#   LOKI_TOKEN=<access token with logs:write>
#
# "ghcr" holds a GitHub PAT (read:packages only) under the key GHCR_TOKEN; .github/workflows/deploy.yml logs docker into ghcr.io with it:
#   aws secretsmanager put-secret-value --secret-id wild-workout-<workspace>/ghcr --secret-string '{"GHCR_TOKEN":"ghp_xxx"}'
