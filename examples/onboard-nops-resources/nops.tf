

terraform {
  required_providers {
    nops = {
      source = "registry.terraform.io/nops-io/nops"
    }
  }
}

provider "nops" {}

# resource "nops_project" "project" {
#   name                        = "nops-provider"
#   account_number              = ""
#   master_payer_account_number = ""
# }

# resource "nops_notification" "notification" {
#   role_arn       = ""
#   external_id    = ""
#   aws_account_id = ""
#   bucket_name    = ""
#   depends_on = [
#     nops_project.project
#   ]
# }

data "nops_projects" "this" {}

output "projects" {
  value = data.nops_projects.this.project

}
