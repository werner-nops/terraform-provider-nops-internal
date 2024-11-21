package nops

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "nops_project" "test" {
  name                        = "automated-testing"
	account_number = "580010171808"
	master_payer_account_number = "580010171808"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// resource.TestCheckResourceAttr("nops_project.test", "items.#", "1"),
					// Verify first order item
					resource.TestCheckResourceAttr("nops_project.test", "name", "automated-testing"),
					resource.TestCheckResourceAttr("nops_project.test", "account_number", "580010171808"),
					resource.TestCheckResourceAttr("nops_project.test", "master_payer_account_number", "580010171808"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "nops_project" "test" {
  name                        = "automated-testing-updated"
	account_number = "471112641702"
	master_payer_account_number = "580010171808"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nops_project.test", "name", "automated-testing-updated"),
					resource.TestCheckResourceAttr("nops_project.test", "account_number", "471112641702"),
				),
			},
		},
	})
}
