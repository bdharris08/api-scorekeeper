/*
 * ecr.tf
 * Creates a Amazon Elastic Container Registry (ECR) for the application
 * https://aws.amazon.com/ecr/
 */

# create an ECR repo at the app/image level
resource "aws_ecr_repository" "app" {
  name                 = var.app
  image_tag_mutability = "MUTABLE"
}

data "aws_caller_identity" "current" {
}
