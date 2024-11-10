package nops

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &projectIntegrationResource{}
	_ resource.ResourceWithConfigure = &projectIntegrationResource{}
)

// projectIntegrationResource is the resource implementation.
type projectIntegrationResource struct {
	client *Client
}

type newProjectIntegrationModel struct {
	ID           types.Int64  `tfsdk:"id"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	ExternalID   types.String `tfsdk:"external_id"`
	AwsAccountID types.String `tfsdk:"aws_account_id"`
	RoleArn      types.String `tfsdk:"role_arn"`
	BucketName   types.String `tfsdk:"bucket_name"`
}

// NewprojectIntegrationResource is a helper function to simplify the provider implementation.
func NewProjectIntegrationResource() resource.Resource {
	return &projectIntegrationResource{}
}

// Configure adds the provider configured client to the resource.
func (r *projectIntegrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *projectIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

// Schema defines the schema for the resource.
func (r *projectIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Notifies the nOps platform a new account has linked to a project with the required input values." +
			" This resource is mostly used only for secure connection with nOps APIs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Integration identifier",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the resource was last updated",
			},
			"role_arn": schema.StringAttribute{
				Required:    true,
				Description: "AWS IAM role to create/update account integration to nOps",
			},
			"bucket_name": schema.StringAttribute{
				Required:    true,
				Description: "AWS S3 bucket name to be used for CUR reports",
			},
			"external_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier to be used by nOps in order to securely assume a role in the target account",
			},
			"aws_account_id": schema.StringAttribute{
				Required:    true,
				Description: "Target AWS account id to integrate with nOps",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *projectIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan newProjectIntegrationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Notify nOps with new values
	var integration Integration
	integration.RoleArn = plan.RoleArn.ValueString()
	integration.BucketName = plan.BucketName.ValueString()
	integration.AccountNumber = plan.AwsAccountID.ValueString()
	integration.ExternalID = plan.ExternalID.ValueString()
	integration.RequestType = "Create"
	integration.ResourceProperties = ResourceProperties{
		ServiceBucket: plan.BucketName.ValueString(),
		AWSAccountID:  plan.AwsAccountID.ValueString(),
		RoleArn:       plan.RoleArn.ValueString(),
		ExternalID:    plan.ExternalID.ValueString(),
	}
	_, err := r.client.NotifyNops(integration)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error notifying nOps",
			"Failed to notify, unexpected error: "+err.Error(),
		)
		return
	}

	// Get updated project values from nOps
	projects, err := r.client.GetProjects()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting remote project data",
			err.Error(),
		)
		return
	}

	for _, project := range projects {
		if types.StringValue(project.AccountNumber) == plan.AwsAccountID {
			// Map response body to schema and populate Computed attribute values
			tflog.Debug(ctx, "Upstream integration project data received for project "+strconv.Itoa(project.ID)+" name: "+project.Name)
			plan.ID = types.Int64Value(int64(project.ID))
			plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Created nOps integration resource", map[string]any{"ID": plan.ID, "AwsAccountID": plan.AwsAccountID})

}

// Read refreshes the Terraform state with the latest data.
func (r *projectIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state newProjectIntegrationModel
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
		if types.StringValue(project.AccountNumber) == state.AwsAccountID {
			// Map response body to schema and populate Computed attribute values
			tflog.Debug(ctx, "Upstream integration project data received for project "+strconv.Itoa(project.ID)+" name: "+project.Name)
			state.ID = types.Int64Value(int64(project.ID))
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
func (r *projectIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan newProjectIntegrationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Notify nOps with updated values
	var integration Integration
	integration.RoleArn = plan.RoleArn.ValueString()
	integration.BucketName = plan.BucketName.ValueString()
	integration.AccountNumber = plan.AwsAccountID.ValueString()
	integration.ExternalID = plan.ExternalID.ValueString()
	integration.RequestType = "Update"
	integration.ResourceProperties = ResourceProperties{
		ServiceBucket: plan.BucketName.ValueString(),
		AWSAccountID:  plan.AwsAccountID.ValueString(),
		RoleArn:       plan.RoleArn.ValueString(),
		ExternalID:    plan.ExternalID.ValueString(),
	}
	_, err := r.client.NotifyNops(integration)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating nOps project",
			"Failed to notify update, unexpected error: "+err.Error(),
		)
		return
	}

	// Get updated project values from nOps
	projects, err := r.client.GetProjects()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting remote project data",
			err.Error(),
		)
		return
	}

	for _, project := range projects {
		if types.StringValue(project.AccountNumber) == plan.AwsAccountID {
			// Map response body to schema and populate Computed attribute values
			tflog.Debug(ctx, "Upstream integration project data received for project "+strconv.Itoa(project.ID)+" name: "+project.Name)
			plan.ID = types.Int64Value(int64(project.ID))
			plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Updated nOps integration resource", map[string]any{"ID": plan.ID, "ExternalID": plan.ExternalID, "LastUpdated": plan.LastUpdated})
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *projectIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No current project delete API on the nOps platform, this is a manual process done in the nOps UI.
	// Framework automatically removes resource from state, no action to be taken on that side.
	var state newProjectIntegrationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Capability to import existing project already integrated in the nOps platform into the TF state without recreation.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("aws_account_id"), req.ID)...)
}
