package nops

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &nopsIntegrationProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &nopsIntegrationProvider{
			version: version,
		}
	}
}

// nopsIntegrationProviderModel maps provider schema data to a Go type.
type nopsIntegrationProviderModel struct {
	ApiKey types.String `tfsdk:"nops_api_key"`
	Host   types.String `tfsdk:"nops_host"`
}

// nopsIntegrationProvider is the provider implementation.
type nopsIntegrationProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *nopsIntegrationProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nops"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *nopsIntegrationProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"nops_api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "nOps API key that will be used for secure communication with the platform APIs, may also be provided with an environment variable NOPS_API_KEY.",
			},
			"nops_host": schema.StringAttribute{
				Optional:    true,
				Description: "nOps API URL, may also be provided with an environment variable NOPS_HOST.",
			},
		},
	}
}

// Configure prepares a nopsIntegration API client for data sources and resources.
func (p *nopsIntegrationProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config nopsIntegrationProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apiKey"),
			"Unknown nOps API key",
			"The provider cannot create the nOps API client as there is an unknown configuration value for the nOps API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NOPS_API_KEY environment variable.",
		)
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown nOps API host",
			"The provider cannot create the nOps API client as there is an unknown configuration value for the nOps API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NOPS_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	apiKey := os.Getenv("NOPS_API_KEY")
	host := os.Getenv("NOPS_HOST")

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apiKey"),
			"Missing nOps API Key",
			"The provider cannot create the nOps API client as there is a missing or empty value for the nOps API key. "+
				"Set the API key value in the configuration or use the NOPS_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if host == "" {
		host = HostURL
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "nops_host", host)
	ctx = tflog.SetField(ctx, "nops_api_key", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "nops_api_key")
	tflog.Debug(ctx, "Creating nops client")

	// Create a new nOps client using the configuration values
	client, err := NewClient(&host, &apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create nOps API Client",
			"An unexpected error occurred when creating the nOps API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"nOps Client Error: "+err.Error(),
		)
		return
	}

	// Make the nOps client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured nOps client", map[string]any{"success": true})

}

// DataSources defines the data sources implemented in the provider.
func (p *nopsIntegrationProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *nopsIntegrationProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewProjectNotificationResource,
	}
}
