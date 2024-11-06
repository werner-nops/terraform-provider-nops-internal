

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
  account_number              = "xxxxxx"
  master_payer_account_number = "xxxxxx"
}

resource "nops_integration" "integration" {
  role_arn       = "arn:aws:iam::xxxxx:role/na"
  external_id    = "NOPS-xxxxxx"
  aws_account_id = "xxxx"
  bucket_name    = "na"
  depends_on = [
    nops_project.project
  ]
}

data "nops_projects" "this" {}

output "projects" {
  value = data.nops_projects.this.projects

}
