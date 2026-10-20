terraform {
  required_version = ">= 1.10"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

module "vpc" {
  source = "../modules/vpc"

  name                 = local.name
  availability_zones   = ["${var.region}a", "${var.region}b"]
  public_subnet_cidr   = "10.0.1.0/24"
  private_subnet_cidrs = ["10.0.11.0/24", "10.0.12.0/24"]
}

module "ec2" {
  source = "../modules/ec2"

  name          = local.name
  vpc_id        = module.vpc.vpc_id
  subnet_id     = module.vpc.public_subnet_id
  app_host      = var.app_host
  api_host      = var.api_host
  certbot_email = var.certbot_email

  grafana_secret_id = module.secrets.secret_arns["backend"]

  secret_arns = concat(
    values(module.secrets.secret_arns),
    [module.rds.master_user_secret_arn],
  )
}

module "rds" {
  source = "../modules/rds"

  name                       = local.name
  vpc_id                     = module.vpc.vpc_id
  subnet_ids                 = module.vpc.private_subnet_ids
  allowed_security_group_ids = [module.ec2.security_group_id]
}
