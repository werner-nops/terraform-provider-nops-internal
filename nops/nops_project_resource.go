package nops

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &projectResource{}
	_ resource.ResourceWithConfigure = &projectResource{}
)

// projectResource is the resource implementation.
type projectResource struct {
	client *Client
}

type newProjectModel struct {
	ID                       types.String `tfsdk:"id"`
	LastUpdated              types.String `tfsdk:"last_updated"`
	Name                     types.String `tfsdk:"name"`
	AccountNumber            types.String `tfsdk:"account_number"`
	MasterPayerAccountNumber types.String `tfsdk:"master_payer_account_number"`
	Arn                      types.String `tfsdk:"arn"`
	Bucket                   types.String `tfsdk:"bucket"`
	Client                   types.Int64  `tfsdk:"client"`
	ExternalID               types.String `tfsdk:"external_id"`
}

// NewProjectResource is a helper function to simplify the provider implementation.
func NewProjectResource() resource.Resource {
	return &projectResource{}
}

// Configure adds the provider configured client to the resource.
func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the resource.
func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "nOps project identifier.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the resource was last updated",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "nOps project name",
			},
			"account_number": schema.StringAttribute{
				Required:    true,
				Description: "Target AWS account id to integrate with nOps",
			},
			"master_payer_account_number": schema.StringAttribute{
				Required:    true,
				Description: "Master payer AWS account id used to conditionally create resources",
			},
			"client": schema.Int64Attribute{
				Computed:    true,
				Description: "nOps client ID",
			},
			"arn": schema.StringAttribute{
				Computed:    true,
				Description: "AWS IAM role to create/update account integration to nOps",
			},
			"bucket": schema.StringAttribute{
				Computed:    true,
				Description: "AWS S3 bucket name to be used for CUR reports",
			},
			"external_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier to be used by nOps in order to securely assume a role in the target account",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan newProjectModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new project
	var newProject NewProject
	newProject.Name = plan.Name.ValueString()
	newProject.AccountNumber = plan.AccountNumber.ValueString()
	newProject.MasterPayerAccountNumber = plan.MasterPayerAccountNumber.ValueString()
	project, err := r.client.CreateProject(newProject)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	tflog.Debug(ctx, "Upstream project data received for project "+strconv.Itoa(project.ID)+" name: "+project.Name)
	plan.ID = types.StringValue(strconv.Itoa(project.ID))
	plan.Client = types.Int64Value(int64(project.Client))
	plan.Arn = types.StringValue(project.Arn)
	plan.Bucket = types.StringValue(project.Bucket)
	plan.Name = types.StringValue(project.Name)
	plan.AccountNumber = types.StringValue(project.AccountNumber)
	plan.ExternalID = types.StringValue(project.ExternalID)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state newProjectModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projects, err := r.client.GetProjects()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting remote project data",
			err.Error(),
		)
		return
	}

	for _, project := range projects {
		if types.StringValue(project.AccountNumber) == state.AccountNumber {
			ctx = tflog.SetField(ctx, "project", project)
			tflog.Debug(ctx, "Upstream project data received for account number "+project.AccountNumber+" name: "+project.Name)
			state.ID = types.StringValue(strconv.Itoa(project.ID))
			state.Client = types.Int64Value(int64(project.Client))
			state.Arn = types.StringValue(project.Arn)
			state.Bucket = types.StringValue(project.Bucket)
			state.Name = types.StringValue(project.Name)
			state.AccountNumber = types.StringValue(project.AccountNumber)
			state.ExternalID = types.StringValue(project.ExternalID)
			state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No current nOps API available to update project data.
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No current nOps API available to delete project
}
