

provider "nops" {}

resource "nops_project" "project" {
  name                        = "nops-provider"
  account_number              = data.aws_caller_identity.current.account_id
  master_payer_account_number = data.aws_organizations_organization.current.master_account_id
}

resource "nops_notification" "notification" {
  role_arn       = aws_iam_role.nops_integration_role.arn
  external_id    = local.external_id
  aws_account_id = local.account_id
  bucket_name    = local.is_master_account ? local.system_bucket_name : "na"
  depends_on = [
    time_sleep.wait_for_resources
  ]
}

data "nops_projects" "this" {}
