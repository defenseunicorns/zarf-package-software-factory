terraform {
  # Follow best practice for root module version constraining
  # See https://www.terraform.io/docs/language/expressions/version-constraints.html
  required_version = ">= 1.2.0, < 2.0.0"
}

locals {
  fullname = "${var.namespace}-${var.stage}-${var.name}"
}

provider "aws" {
  region = var.aws_region
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE VPC
# ---------------------------------------------------------------------------------------------------------------------
resource "aws_vpc" "terratest_vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "terratest-vpc"
  }
}

# Create an Internet Gateway for the VPC
resource "aws_internet_gateway" "terratest_igw" {
  vpc_id = aws_vpc.terratest_vpc.id
}

# Create a public subnet in the VPC
resource "aws_subnet" "terratest_public_subnet" {
  vpc_id                  = aws_vpc.terratest_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = var.aws_availability_zone
  map_public_ip_on_launch = true

  tags = {
    Name = "terratest-public-subnet"
  }
}

# Create a route table associated with the public subnet
resource "aws_route_table" "terratest_public_rt" {
  vpc_id = aws_vpc.terratest_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.terratest_igw.id
  }

  tags = {
    Name = "terratest-public-route-table"
  }
}

# Associate the route table with the public subnet
resource "aws_route_table_association" "terratest_public_rt_assoc" {
  subnet_id      = aws_subnet.terratest_public_subnet.id
  route_table_id = aws_route_table.terratest_public_rt.id
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A PUBLIC EC2 INSTANCE
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_instance" "public" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = var.instance_type
  vpc_security_group_ids = [aws_security_group.public.id]
  key_name               = var.key_pair_name
  subnet_id              = aws_subnet.terratest_public_subnet.id

  root_block_device {
    volume_size = 400
  }

  user_data = <<_EOF_
#!/bin/bash
echo "PubkeyAcceptedKeyTypes=+ssh-rsa" >> /etc/ssh/sshd_config
service ssh reload
sysctl fs.inotify.max_user_instances=512
sysctl -p
_EOF_

  # This EC2 Instance has a public IP and will be accessible directly from the public Internet
  associate_public_ip_address = true

  tags = {
    Name      = "${local.fullname}-public"
    CreatedBy = "${data.aws_caller_identity.whoami.arn}"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A SECURITY GROUP TO CONTROL WHAT REQUESTS CAN GO IN AND OUT OF THE EC2 INSTANCES
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_security_group" "public" {
  name = local.fullname

  vpc_id = aws_vpc.terratest_vpc.id

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
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# LOOK UP THE CALLER ID
# ---------------------------------------------------------------------------------------------------------------------

data "aws_caller_identity" "whoami" {}
