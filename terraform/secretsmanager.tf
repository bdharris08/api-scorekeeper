# Allow for code reuse...
locals {
  # KMS write actions
  kms_write_actions = [
    "kms:CancelKeyDeletion",
    "kms:CreateAlias",
    "kms:CreateGrant",
    "kms:CreateKey",
    "kms:DeleteAlias",
    "kms:DeleteImportedKeyMaterial",
    "kms:DisableKey",
    "kms:DisableKeyRotation",
    "kms:EnableKey",
    "kms:EnableKeyRotation",
    "kms:Encrypt",
    "kms:GenerateDataKey",
    "kms:GenerateDataKeyWithoutPlaintext",
    "kms:GenerateRandom",
    "kms:GetKeyPolicy",
    "kms:GetKeyRotationStatus",
    "kms:GetParametersForImport",
    "kms:ImportKeyMaterial",
    "kms:PutKeyPolicy",
    "kms:ReEncryptFrom",
    "kms:ReEncryptTo",
    "kms:RetireGrant",
    "kms:RevokeGrant",
    "kms:ScheduleKeyDeletion",
    "kms:TagResource",
    "kms:UntagResource",
    "kms:UpdateAlias",
    "kms:UpdateKeyDescription",
  ]

  # KMS read actions
  kms_read_actions = [
    "kms:Decrypt",
    "kms:DescribeKey",
    "kms:List*",
  ]

  # secretsmanager write actions
  sm_write_actions = [
    "secretsmanager:CancelRotateSecret",
    "secretsmanager:CreateSecret",
    "secretsmanager:DeleteSecret",
    "secretsmanager:PutSecretValue",
    "secretsmanager:RestoreSecret",
    "secretsmanager:RotateSecret",
    "secretsmanager:TagResource",
    "secretsmanager:UntagResource",
    "secretsmanager:UpdateSecret",
    "secretsmanager:UpdateSecretVersionStage",
  ]

  # secretsmanager read actions
  sm_read_actions = [
    "secretsmanager:DescribeSecret",
    "secretsmanager:List*",
    "secretsmanager:GetRandomPassword",
    "secretsmanager:GetSecretValue",
  ]

  sm_arn = "arn:aws:secretsmanager:${var.region}:${data.aws_caller_identity.current.account_id}:secret:${var.app}-${var.environment}-??????"
}

# create the KMS key for this secret
resource "aws_kms_key" "sm_kms_key" {
  description = "${var.app}-${var.environment}"
  policy      = data.aws_iam_policy_document.kms_resource_policy_doc.json
  tags = merge(
    var.tags,
    {
      "Name" = format("%s-%s", var.app, var.environment)
    },
  )
}

# alias for the key
resource "aws_kms_alias" "sm_kms_alias" {
  name          = "alias/${var.app}-${var.environment}"
  target_key_id = aws_kms_key.sm_kms_key.key_id
}

# the kms key policy
data "aws_iam_policy_document" "kms_resource_policy_doc" {
  statement {
    sid    = "AllowWrite"
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    actions   = local.kms_write_actions
    resources = ["*"]
  }

  statement {
    sid    = "AllowReadRole"
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    actions   = local.kms_read_actions
    resources = ["*"]
  }
}

# create the secretsmanager secret
resource "aws_secretsmanager_secret" "sm_secret" {
  name       = "${var.app}-${var.environment}-1"
  kms_key_id = aws_kms_key.sm_kms_key.key_id
  tags       = var.tags
  policy     = data.aws_iam_policy_document.sm_resource_policy_doc.json
}

# create the placeholder secret json
resource "aws_secretsmanager_secret_version" "initial" {
  secret_id     = aws_secretsmanager_secret.sm_secret.id
  secret_string = "{}"
}

# resource policy doc that limits access to secret
data "aws_iam_policy_document" "sm_resource_policy_doc" {
  statement {
    sid    = "AllowWrite"
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    actions   = local.sm_write_actions
    resources = [local.sm_arn]

    condition {
      test     = "StringLike"
      variable = "aws:userId"
      values   = ["ben"]
    }
  }

  statement {
    sid    = "AllowReadRole"
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    actions   = local.sm_read_actions
    resources = [local.sm_arn]

    condition {
      test     = "StringLike"
      variable = "aws:userId"
      values   = ["ben"]
    }
  }
}
