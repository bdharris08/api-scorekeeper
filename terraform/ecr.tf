resource "aws_ecr_repository" "ecr" {
  name                 = "api-scorekeeper"
  image_tag_mutability = "MUTABLE"
}
