data "aws_caller_identity" "current" {}

data "aws_organizations_organization" "current" {}

# Resource intended to be used for the initial onboarding of an account to the nOps platform, 
# used for communication with nOps APIs.
resource "nops_project" "project" {
  name                        = "project"
  account_number              = data.aws_caller_identity.current.account_id
  master_payer_account_number = data.aws_organizations_organization.current.master_account_id
}
