# One JSON secret holding the backend's runtime config (expected keys: local.backend_secret_keys in locals.tf).
# It is the only secret we keep, to save the per-secret Secrets Manager fee. Besides the backend env vars it holds:
#   LOKI_URL, LOKI_USER, LOKI_TOKEN  Grafana Cloud Loki credentials, read by Alloy on every start (see modules/ec2/scripts/40-alloy.sh.tftpl). Set them with `terraform/prod/scripts/set-loki-credentials.sh`.
#   GHCR_TOKEN                       GitHub PAT (read:packages only); .github/workflows/deploy.yml logs docker into ghcr.io with it.
# Neither group is passed to the backend container (see deployment/write-backend-env.sh).
# Terraform creates the container only; set the value by hand so it stays out of state:
#   aws secretsmanager put-secret-value --secret-id wild-workout-<workspace>/backend --secret-string file://backend.json
module "secrets" {
  source = "../modules/ssm"

  name         = local.name
  secret_names = ["backend"]
}
