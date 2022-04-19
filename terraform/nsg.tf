resource "aws_security_group" "nsg_task" {
  name        = "${var.app}-${var.environment}-task"
  description = "Limit connections from internal resources while allowing ${var.app}-${var.environment}-task to connect to all external resources"
  vpc_id      = aws_vpc.main.id

  tags = var.tags
}

# Rules for the TASK (Targets the LB's IPs)
resource "aws_security_group_rule" "nsg_task_ingress_rule" {
  description = "Only allow connections from the NLB on port ${var.container_port}"
  type        = "ingress"
  from_port   = var.container_port
  to_port     = var.container_port
  protocol    = "tcp"
  cidr_blocks = formatlist("%s/32", [for eni in data.aws_network_interface.nlb : eni.private_ip])

  security_group_id = aws_security_group.nsg_task.id
}

resource "aws_security_group_rule" "nsg_task_egress_rule" {
  description     = "Allows task to establish connections to all resources"
  type            = "egress"
  from_port       = "0"
  to_port         = "0"
  protocol        = "-1"
  cidr_blocks     = ["0.0.0.0/0"]
  prefix_list_ids = [aws_vpc_endpoint.s3.prefix_list_id]

  security_group_id = aws_security_group.nsg_task.id
}

resource "aws_security_group_rule" "nsg_task_egress_rule_vpc_link" {
  description     = "Allows task to establish connections through VPC Link"
  type            = "egress"
  from_port       = "443"
  to_port         = "443"
  protocol        = "tcp"
  cidr_blocks     = [aws_vpc.main.cidr_block]
  prefix_list_ids = [aws_vpc_endpoint.s3.prefix_list_id]

  security_group_id = aws_security_group.nsg_task.id
}

# lookup the ENIs associated with the NLB
data "aws_network_interface" "nlb" {
  for_each = aws_lb.main.subnets

  filter {
    name   = "description"
    values = ["ELB ${aws_lb.main.arn_suffix}"]
  }

  filter {
    name   = "subnet-id"
    values = [each.value]
  }
}

