locals {
  name = "wild-workout-${terraform.workspace}"

  # Keys that must exist in the "backend" secret's JSON before the backend task can start.
  backend_secret_keys = [
    "APP_ENV",
    "CLERK_AUTHORIZED_PARTIES",
    "CLERK_SECRET_KEY",
    "CLERK_WEBHOOK_SIGNING_SECRET",
    "CORS_ALLOWED_ORIGINS",
    "GHCR_TOKEN",
    "LOKI_TOKEN",
    "LOKI_URL",
    "LOKI_USER",
  ]
}
