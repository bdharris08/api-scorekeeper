resource "aws_vpc" "main" {
  cidr_block       = "10.0.0.0/16"
  instance_tenancy = "default"
  enable_dns_hostnames = "true"
  enable_dns_support   = "true"

  tags = {
    Name = var.vpc
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id
}

resource "aws_subnet" "private_d" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.10.0/25"
  availability_zone = "us-east-1d"

  tags = {
    "Name" = "${var.app} | private | us-east-1d"
  }
}

resource "aws_subnet" "private_e" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.11.0/25"
  availability_zone = "us-east-1e"

  tags = {
    "Name" = "${var.app} | private | us-east-1d"
  }
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id
  tags = {
    "Name" = "${var.app} | private"
  }
}

resource "aws_route_table_association" "private_d_subnet" {
  subnet_id      = aws_subnet.private_d.id
  route_table_id = aws_route_table.private.id
}

resource "aws_security_group" "egress_all" {
  name        = "${var.app}-egress-all"
  description = "Allow all outbound traffic for ${var.app}"
  vpc_id      = aws_vpc.main.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
