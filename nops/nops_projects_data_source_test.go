package nops

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProjectsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
data "nops_projects" "test" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of projects returned
					// resource.TestCheckResourceAttr("data.nops_projects.test", "projects.#", "9"),
					// Verify the second project to ensure all attributes are set, this test runs on a mock client created in nOPS
					resource.TestCheckResourceAttr("data.nops_projects.test", "projects.1.client", "8549"),
					resource.TestCheckResourceAttr("data.nops_projects.test", "projects.1.arn", "arn:aws:iam::580010171808:role/na"),
					resource.TestCheckResourceAttr("data.nops_projects.test", "projects.1.bucket", ""),
				),
			},
		},
	})
}
