data "nops_projects" "this" {}

output "projects" {
  value = data.nops_projects.this.projects

}
