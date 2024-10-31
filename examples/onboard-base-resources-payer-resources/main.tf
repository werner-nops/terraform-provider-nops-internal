

terraform {
  required_providers {
    nops = {
      source  = "terraform.local/custom/nops"
      version = "1.0.1"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "5.73.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "0.12.1"
    }
  }
}

provider "aws" {
  profile = "nops-root"
  region  = "us-west-2"
}

resource "time_sleep" "wait_for_resources" {
  depends_on = [
    aws_iam_role.nops_integration_role,
    aws_iam_role_policy.nops_integration_policy,
    aws_iam_role_policy.nops_system_bucket_policy,
    aws_iam_role_policy.nops_essentials_policy,
    aws_iam_role_policy.nops_compute_copilot_policy,
    aws_iam_role_policy.nops_wafr_policy,
    aws_s3_bucket.nops_system_bucket,
    aws_s3_bucket_policy.nops_bucket_policy,
    aws_s3_bucket_server_side_encryption_configuration.nops_bucket_encryption
  ]

  create_duration = "10s"
}

