

terraform {
  required_providers {
    nops = {
      source = "hashicorp.com/edu/nops"
    }
  }
}

provider "nops" {}

resource "nops_project" "project" {
  name                        = "nops-provider"
  account_number              = ""
  master_payer_account_number = ""
}

resource "nops_integration" "integration" {
  role_arn       = "arn:aws:iam::580010171808:role/na"
  external_id    = "NOPS-24317961A6AFCA4BA9AB32BAF77759"
  aws_account_id = "580010171808"
  bucket_name    = "na"
  depends_on = [
    nops_project.project
  ]
}

# data "nops_projects" "this" {}

# output "projects" {
#   value = data.nops_projects.this.projects

# }
