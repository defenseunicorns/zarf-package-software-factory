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
  ami                    = data.aws_ami.default.id
  instance_type          = var.instance_type
  vpc_security_group_ids = [aws_security_group.public.id]
  key_name               = var.key_pair_name

  root_block_device {
    volume_size = 200
  }

  # This EC2 Instance has a public IP and will be accessible directly from the public Internet
  associate_public_ip_address = true

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

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# LOOK UP THE AMI
# ---------------------------------------------------------------------------------------------------------------------

data "aws_ami" "default" {
  most_recent = true
  owners      = ["309956199498"]

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
    values = ["RHEL-8.*_HVM-*-x86_64-2-Hourly2-GP2"]
  }
}
