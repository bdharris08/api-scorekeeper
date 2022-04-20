/*
 * variables.tf
 * Common variables to use in various Terraform files (*.tf)
 */

# The AWS region to use for the bucket and registry; typically `us-east-1`.
# Other possible values: `us-east-2`, `us-west-1`, or `us-west-2`.
# Currently, Fargate is only available in `us-east-1`.
variable "region" {
  default = "us-east-1"
}

# Name of the application. This value should usually match the application tag below.
variable "app" {
}

# The environment that is being built
variable "environment" {
}

# A map of the tags to apply to various resources. The required tags are:
# `application`, name of the app;
# `environment`, the environment being created;
# `team`, team responsible for the application;
# `contact-email`, contact email for the _team_;
# and `customer`, who the application was create for.
variable "tags" {
  type = map(string)
}

# The port the container will listen on, used for load balancer health check
# Best practice is that this value is higher than 1024 so the container processes
# isn't running at root.
variable "container_port" {
}

# The port the load balancer will listen on
variable "lb_port" {
  default = "80"
}

# The load balancer protocol
variable "lb_protocol" {
  default = "TCP"
}

# Database URI secret
variable "database_url" {
}

# Network configuration

# The VPC to use for the Fargate cluster
variable "vpc" {
}

locals {
  namespace = "${var.app}-${var.environment}"
}
