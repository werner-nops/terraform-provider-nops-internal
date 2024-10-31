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
	_ resource.Resource              = &projectNotificationResource{}
	_ resource.ResourceWithConfigure = &projectNotificationResource{}
)

// projectNotificationResource is the resource implementation.
type projectNotificationResource struct {
	client *Client
}

type newProjectNotificationModel struct {
	ID           types.String `tfsdk:"id"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	ExternalID   types.String `tfsdk:"external_id"`
	AwsAccountID types.String `tfsdk:"aws_account_id"`
	RoleArn      types.String `tfsdk:"role_arn"`
	BucketName   types.String `tfsdk:"bucket_name"`
}

// NewprojectNotificationResource is a helper function to simplify the provider implementation.
func NewProjectNotificationResource() resource.Resource {
	return &projectNotificationResource{}
}

// Configure adds the provider configured client to the resource.
func (r *projectNotificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *projectNotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

// Schema defines the schema for the resource.
func (r *projectNotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Notification identifier",
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
func (r *projectNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan newProjectNotificationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Notify nOps with new values
	var notification Notification
	notification.RoleArn = plan.RoleArn.ValueString()
	notification.BucketName = plan.BucketName.ValueString()
	notification.AccountNumber = plan.AwsAccountID.ValueString()
	notification.ExternalID = plan.ExternalID.ValueString()
	notification.RequestType = "Create"
	notification.ResourceProperties = ResourceProperties{
		ServiceBucket: plan.BucketName.ValueString(),
		AWSAccountID:  plan.AwsAccountID.ValueString(),
		RoleArn:       plan.RoleArn.ValueString(),
		ExternalID:    plan.ExternalID.ValueString(),
	}
	_, err := r.client.NotifyNops(notification)
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
			tflog.Debug(ctx, "Upstream notification project data received for project "+strconv.Itoa(project.ID)+" name: "+project.Name)
			plan.ID = types.StringValue(strconv.Itoa(project.ID))
			plan.RoleArn = types.StringValue(project.Arn)
			plan.BucketName = types.StringValue(project.Bucket)
			plan.AwsAccountID = types.StringValue(project.AccountNumber)
			plan.ExternalID = types.StringValue(project.ExternalID)
			plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Created nOps notification resource", map[string]any{"ID": plan.ID, "AwsAccountID": plan.AwsAccountID})

}

// Read refreshes the Terraform state with the latest data.
func (r *projectNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state newProjectNotificationModel
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
			tflog.Debug(ctx, "Upstream notification project data received for project "+strconv.Itoa(project.ID)+" name: "+project.Name)
			state.ID = types.StringValue(strconv.Itoa(project.ID))
			state.RoleArn = types.StringValue(project.Arn)
			state.BucketName = types.StringValue(project.Bucket)
			state.AwsAccountID = types.StringValue(project.AccountNumber)
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
func (r *projectNotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan newProjectNotificationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Notify nOps with updated values
	var notification Notification
	notification.RoleArn = plan.RoleArn.ValueString()
	notification.BucketName = plan.BucketName.ValueString()
	notification.AccountNumber = plan.AwsAccountID.ValueString()
	notification.ExternalID = plan.ExternalID.ValueString()
	notification.RequestType = "Update"
	notification.ResourceProperties = ResourceProperties{
		ServiceBucket: plan.BucketName.ValueString(),
		AWSAccountID:  plan.AwsAccountID.ValueString(),
		RoleArn:       plan.RoleArn.ValueString(),
		ExternalID:    plan.ExternalID.ValueString(),
	}
	_, err := r.client.NotifyNops(notification)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating nOps project",
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
			tflog.Debug(ctx, "Upstream notification project data received for project "+strconv.Itoa(project.ID)+" name: "+project.Name)
			plan.ID = types.StringValue(strconv.Itoa(project.ID))
			plan.RoleArn = types.StringValue(project.Arn)
			plan.BucketName = types.StringValue(project.Bucket)
			plan.AwsAccountID = types.StringValue(project.AccountNumber)
			plan.ExternalID = types.StringValue(project.ExternalID)
			plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Updated nOps notification resource", map[string]any{"ID": plan.ID, "ExternalID": plan.ExternalID, "LastUpdated": plan.LastUpdated})
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *projectNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No current project delete API on the nOps platform, this is a manual process done in the UI
	var state newProjectNotificationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
