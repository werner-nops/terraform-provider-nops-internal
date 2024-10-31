package nops

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &projectsDataSource{}
	_ datasource.DataSourceWithConfigure = &projectsDataSource{}
)

func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

// Data source implementation.
type projectsDataSource struct {
	client *Client
}

type projectsDataSourceModel struct {
	Projects []projectsModel `tfsdk:"projects"`
}

type projectsModel struct {
	ID     types.Int64  `tfsdk:"id"`
	Client types.Int64  `tfsdk:"client"`
	Arn    types.String `tfsdk:"arn"`
	Bucket types.String `tfsdk:"bucket"`
}

// Metadata returns the data source type name.
func (d *projectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

// Schema defines the schema for the data source.
func (d *projectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "nOps project identifier",
						},
						"client": schema.Int64Attribute{
							Computed:    true,
							Description: "nOps client identifier",
						},
						"arn": schema.StringAttribute{
							Computed:    true,
							Description: "AWS IAM role to create/update account integration to nOps",
						},
						"bucket": schema.StringAttribute{
							Computed:    true,
							Description: "AWS S3 bucket name to be used for CUR reports",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state projectsDataSourceModel

	projects, err := d.client.GetProjects()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting remote project data",
			err.Error(),
		)
		return
	}

	for _, project := range projects {
		ctx = tflog.SetField(ctx, "project", project)
		tflog.Debug(ctx, "Got project data")
		projectState := projectsModel{
			ID:     types.Int64Value(int64(project.ID)),
			Client: types.Int64Value(int64(project.Client)),
			Arn:    types.StringValue(project.Arn),
			Bucket: types.StringValue(project.Bucket),
		}

		state.Projects = append(state.Projects, projectState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *projectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
