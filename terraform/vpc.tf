resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  instance_tenancy     = "default"
  enable_dns_hostnames = "true"
  enable_dns_support   = "true"

  tags = {
    Name = var.vpc
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id
}

resource "aws_subnet" "private_b" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.10.0/25"
  availability_zone = "us-east-1b"

  tags = {
    "Name" = "${var.app} | private | us-east-1b"
  }
}

resource "aws_subnet" "private_c" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.11.0/25"
  availability_zone = "us-east-1c"

  tags = {
    "Name" = "${var.app} | private | us-east-1c"
  }
}

resource "aws_db_subnet_group" "private_group" {
  name       = "${var.app}-private-subnet-group"
  subnet_ids = [aws_subnet.private_b.id, aws_subnet.private_c.id]

  tags = {
    "Name" = "${var.app}"
  }
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id
  tags = {
    "Name" = "${var.app} | private"
  }
}

resource "aws_route_table_association" "private_d_subnet" {
  subnet_id      = aws_subnet.private_b.id
  route_table_id = aws_route_table.private.id
}

resource "aws_route_table_association" "private_e_subnet" {
  subnet_id      = aws_subnet.private_c.id
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

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block]
  }
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${var.region}.s3"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = [aws_route_table.private.id]

  tags = {
    Name        = "s3-endpoint"
    Environment = "demo"
  }
}

resource "aws_vpc_endpoint" "dkr" {
  vpc_id              = aws_vpc.main.id
  private_dns_enabled = true
  service_name        = "com.amazonaws.${var.region}.ecr.dkr"
  vpc_endpoint_type   = "Interface"
  security_group_ids = [
    aws_security_group.egress_all.id,
    aws_security_group.nsg_task.id,
  ]
  subnet_ids = [aws_subnet.private_b.id, aws_subnet.private_c.id]

  tags = {
    Name        = "dkr-endpoint"
    Environment = "demo"
  }
}

resource "aws_vpc_endpoint" "ecr-api" {
  vpc_id              = aws_vpc.main.id
  private_dns_enabled = true
  service_name        = "com.amazonaws.${var.region}.ecr.api"
  vpc_endpoint_type   = "Interface"
  security_group_ids = [
    aws_security_group.egress_all.id,
    aws_security_group.nsg_task.id,
  ]
  subnet_ids = [aws_subnet.private_b.id, aws_subnet.private_c.id]

  tags = {
    Name        = "ecr-api-endpoint"
    Environment = "demo"
  }
}

resource "aws_vpc_endpoint" "logs" {
  vpc_id              = aws_vpc.main.id
  private_dns_enabled = true
  service_name        = "com.amazonaws.${var.region}.logs"
  vpc_endpoint_type   = "Interface"
  security_group_ids = [
    aws_security_group.egress_all.id,
    aws_security_group.nsg_task.id,
  ]
  subnet_ids = [aws_subnet.private_b.id, aws_subnet.private_c.id]

  tags = {
    Name        = "logs-endpoint"
    Environment = "demo"
  }
}
