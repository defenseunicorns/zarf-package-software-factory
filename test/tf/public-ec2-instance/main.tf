terraform {
  # Follow best practice for root module version constraing
  # See https://www.terraform.io/docs/language/expressions/version-constraints.html
  required_version = "~> 1.2.0"
}

locals {
  fullname = "${var.namespace}-${var.stage}-${var.name}"
}

provider "aws" {
  region = var.aws_region
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A PUBLIC EC2 INSTANCE
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_instance" "public" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = var.instance_type
  vpc_security_group_ids = [aws_security_group.public.id]
  key_name               = var.key_pair_name


  # This EC2 Instance has a public IP and will be accessible directly from the public Internet
  associate_public_ip_address = true

  user_data = <<EOF
#!/bin/bash

# install deps
apt-get update
apt-get install -y jq
apt-get install -y git
apt-get install -y make
apt-get install -y wget

# elasticsearch needs this
sysctl -w vm.max_map_count=262144

EOF

  tags = {
    Name = "${local.fullname}-public"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A SECURITY GROUP TO CONTROL WHAT REQUESTS CAN GO IN AND OUT OF THE EC2 INSTANCES
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_security_group" "public" {
  name = local.fullname

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port = 22
    to_port   = 22
    protocol  = "tcp"

    # To keep this example simple, we allow incoming SSH requests from any IP. In real-world usage, you should only
    # allow SSH requests from trusted servers, such as a bastion host or VPN server.
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# LOOK UP THE LATEST UBUNTU AMI
# ---------------------------------------------------------------------------------------------------------------------

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "image-type"
    values = ["machine"]
  }

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }
}
