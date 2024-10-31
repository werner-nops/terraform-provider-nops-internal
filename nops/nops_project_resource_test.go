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
			// Update and Read testing - not implemented
			// 			{
			// 				Config: providerConfig + `
			// resource "hashicups_order" "test" {
			//   items = [
			//     {
			//       coffee = {
			//         id = 2
			//       }
			//       quantity = 2
			//     },
			//   ]
			// }
			// `,
			// 				Check: resource.ComposeAggregateTestCheckFunc(
			// 					// Verify first order item updated
			// 					resource.TestCheckResourceAttr("hashicups_order.test", "items.0.quantity", "2"),
			// 					resource.TestCheckResourceAttr("hashicups_order.test", "items.0.coffee.id", "2"),
			// 					// Verify first coffee item has Computed attributes updated.
			// 					resource.TestCheckResourceAttr("hashicups_order.test", "items.0.coffee.description", ""),
			// 					resource.TestCheckResourceAttr("hashicups_order.test", "items.0.coffee.image", "/packer.png"),
			// 					resource.TestCheckResourceAttr("hashicups_order.test", "items.0.coffee.name", "Packer Spiced Latte"),
			// 					resource.TestCheckResourceAttr("hashicups_order.test", "items.0.coffee.price", "350"),
			// 					resource.TestCheckResourceAttr("hashicups_order.test", "items.0.coffee.teaser", "Packed with goodness to spice up your images"),
			// 				),
			// 			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
