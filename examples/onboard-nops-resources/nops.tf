

terraform {
  required_providers {
    nops = {
      source = "hashicorp.com/edu/nops"
    }
  }
}

provider "nops" {
  nops_api_key = "8650.55076cfa3656fe9cbd0671eef9e89666"
}

resource "nops_project" "project" {
  name                        = "nops-provider"
  account_number              = "202279780353"
  master_payer_account_number = "728471903238"
}

# resource "nops_integration" "integration" {
#   role_arn       = "arn:aws:iam::202279780353:role/na"
#   external_id    = "NOPS-24317961A6AFCA4BA9AB32BAF77759"
#   aws_account_id = "202279780353"
#   bucket_name    = "na"
#   depends_on = [
#     nops_project.project
#   ]
# }

# data "nops_projects" "this" {}

# output "projects" {
#   value = data.nops_projects.this.projects

# }
