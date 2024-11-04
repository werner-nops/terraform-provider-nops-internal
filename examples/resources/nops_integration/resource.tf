data "aws_caller_identity" "current" {}

data "aws_organizations_organization" "current" {}

resource "nops_project" "project" {
  name                        = "project"
  account_number              = data.aws_caller_identity.current.account_id
  master_payer_account_number = data.aws_organizations_organization.current.master_account_id
}

# Notifies the nOps platform a new account has linked to a project with the required input values. 
# This resource is mostly used only for secure connection with nOps APIs.
resource "nops_integration" "integration" {
  # Role with sufficient permissions for nOps to get cost and metadata
  role_arn       = aws_iam_role.nops_integration_role.arn
  external_id    = nops_project.project.external_id
  aws_account_id = data.aws_caller_identity.current.account_id
  # If being deployed in a management account set S3 bucket name, if not value should be "na"
  bucket_name = aws_s3_bucket.nops_system_bucket.id
  depends_on = [
    nops_project.project
  ]
}