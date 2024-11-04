---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nops_integration Resource - nops"
subcategory: ""
description: |-
  Notifies the nOps platform a new account has linked to a project with the required input values. This resource is mostly used only for secure connection with nOps APIs.
---

# nops_integration (Resource)

Notifies the nOps platform a new account has linked to a project with the required input values. This resource is mostly used only for secure connection with nOps APIs.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `aws_account_id` (String) Target AWS account id to integrate with nOps
- `bucket_name` (String) AWS S3 bucket name to be used for CUR reports
- `external_id` (String) Identifier to be used by nOps in order to securely assume a role in the target account
- `role_arn` (String) AWS IAM role to create/update account integration to nOps

### Read-Only

- `id` (Number) Integration identifier
- `last_updated` (String) Timestamp when the resource was last updated