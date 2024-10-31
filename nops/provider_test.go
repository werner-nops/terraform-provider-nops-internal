package nops

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"nops": providerserver.NewProtocol6WithError(New("test")()),
	}
	nops_api_key = os.Getenv("NOPS_API_KEY")
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the HashiCups client is properly configured.
	// It is also possible to use the environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = fmt.Sprintf(`
provider "nops" {
  nops_api_key="%s"
}`, nops_api_key,
	)
)
