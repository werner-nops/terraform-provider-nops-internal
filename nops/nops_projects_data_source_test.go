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
					resource.TestCheckResourceAttr("data.nops_projects.test", "projects.#", "1"),
					// Verify the test project to ensure all attributes are set, this test runs on a mock client created in nOPS UAT account (tf-automated-testing)
					resource.TestCheckResourceAttr("data.nops_projects.test", "projects.0.client", "15418"),
					resource.TestCheckResourceAttr("data.nops_projects.test", "projects.0.arn", "arn:aws:iam::471112641702:role/na"),
					resource.TestCheckResourceAttr("data.nops_projects.test", "projects.0.bucket", ""),
				),
			},
		},
	})
}
