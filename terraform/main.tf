terraform {
  required_version = ">= 0.14.9"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.8"
    }
  }

  cloud {
    organization = "testbuildpleaseignore"

    workspaces {
      name = "api-scorekeeper"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

output "aws_account" {
  value = "621387225812"
}

output "region" {
  value = "us-east-1"
}
