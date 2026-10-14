data "aws_region" "current" {}

data "aws_ssm_parameter" "ami" {
  name = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-arm64"
}

# Instance bootstrap, one script per concern in scripts/, rendered and run in filename order.
locals {
  nginx_conf = templatefile("${path.module}/files/nginx.conf.tftpl", {
    app_host = var.app_host
    api_host = var.api_host
    app_port = var.app_port
    api_port = var.api_port
  })

  alloy_conf = templatefile("${path.module}/files/alloy.config.tftpl", {})

  user_data = join("\n", [
    templatefile("${path.module}/scripts/00-base.sh.tftpl", {}),
    templatefile("${path.module}/scripts/10-docker.sh.tftpl", {}),
    templatefile("${path.module}/scripts/20-nginx.sh.tftpl", {
      nginx_conf = local.nginx_conf
    }),
    templatefile("${path.module}/scripts/30-certbot.sh.tftpl", {
      app_host      = var.app_host
      api_host      = var.api_host
      certbot_email = var.certbot_email
      eip           = aws_eip.this.public_ip
    }),
    templatefile("${path.module}/scripts/40-alloy.sh.tftpl", {
      alloy_conf        = local.alloy_conf
      region            = data.aws_region.current.name
      grafana_secret_id = var.grafana_secret_id
    }),
  ])
}

resource "aws_security_group" "this" {
  name_prefix = "${var.name}-ec2-"
  description = "Public web ingress for ${var.name}"
  vpc_id      = var.vpc_id

  tags = merge(var.tags, { Name = "${var.name}-ec2-sg" })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_vpc_security_group_ingress_rule" "http" {
  security_group_id = aws_security_group.this.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "tcp"
  from_port         = 80
  to_port           = 80
}

resource "aws_vpc_security_group_ingress_rule" "https" {
  security_group_id = aws_security_group.this.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "tcp"
  from_port         = 443
  to_port           = 443
}

resource "aws_vpc_security_group_ingress_rule" "ssh" {
  for_each = toset(var.ssh_cidr_blocks)

  security_group_id = aws_security_group.this.id
  cidr_ipv4         = each.value
  ip_protocol       = "tcp"
  from_port         = 22
  to_port           = 22
}

resource "aws_vpc_security_group_egress_rule" "all" {
  security_group_id = aws_security_group.this.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}

data "aws_iam_policy_document" "assume" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "this" {
  name               = "${var.name}-ec2"
  assume_role_policy = data.aws_iam_policy_document.assume.json

  tags = var.tags
}

# Session Manager access instead of SSH.
resource "aws_iam_role_policy_attachment" "ssm" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

data "aws_iam_policy_document" "secrets" {
  count = length(var.secret_arns) > 0 ? 1 : 0

  statement {
    actions   = ["secretsmanager:GetSecretValue"]
    resources = var.secret_arns
  }
}

resource "aws_iam_role_policy" "secrets" {
  count = length(var.secret_arns) > 0 ? 1 : 0

  name   = "read-secrets"
  role   = aws_iam_role.this.id
  policy = data.aws_iam_policy_document.secrets[0].json
}

resource "aws_iam_instance_profile" "this" {
  name = "${var.name}-ec2"
  role = aws_iam_role.this.name

  tags = var.tags
}

resource "aws_eip" "this" {
  domain = "vpc"

  tags = merge(var.tags, { Name = "${var.name}-ec2-eip" })
}

resource "aws_instance" "this" {
  ami                    = data.aws_ssm_parameter.ami.value
  instance_type          = var.instance_type
  subnet_id              = var.subnet_id
  vpc_security_group_ids = [aws_security_group.this.id]
  iam_instance_profile   = aws_iam_instance_profile.this.name

  user_data                   = local.user_data
  user_data_replace_on_change = true

  metadata_options {
    http_tokens = "required"
  }

  root_block_device {
    volume_type = "gp3"
    volume_size = var.root_volume_size
    encrypted   = true
  }

  tags = merge(var.tags, { Name = "${var.name}-ec2" })
}

resource "aws_eip_association" "this" {
  instance_id   = aws_instance.this.id
  allocation_id = aws_eip.this.id
}
