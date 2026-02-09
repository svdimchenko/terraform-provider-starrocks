package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-starrocks/internal/client"
)

var (
	_ resource.Resource                = &resourceGroupResource{}
	_ resource.ResourceWithConfigure   = &resourceGroupResource{}
	_ resource.ResourceWithImportState = &resourceGroupResource{}
)

func NewResourceGroupResource() resource.Resource {
	return &resourceGroupResource{}
}

type resourceGroupResource struct {
	client *client.Client
}

type resourceGroupResourceModel struct {
	Name                     types.String `tfsdk:"name"`
	CPUCoreLimit             types.Int64  `tfsdk:"cpu_core_limit"`
	MemLimit                 types.String `tfsdk:"mem_limit"`
	ConcurrencyLimit         types.Int64  `tfsdk:"concurrency_limit"`
	BigQueryMemLimit         types.String `tfsdk:"big_query_mem_limit"`
	BigQueryScanRowsLimit    types.Int64  `tfsdk:"big_query_scan_rows_limit"`
	BigQueryCPUSecondLimit   types.Int64  `tfsdk:"big_query_cpu_second_limit"`
	Classifiers              types.Set    `tfsdk:"classifiers"`
}

func (m *resourceGroupResourceModel) GetName() types.String { return m.Name }
func (m *resourceGroupResourceModel) GetCPUCoreLimit() types.Int64 { return m.CPUCoreLimit }
func (m *resourceGroupResourceModel) GetMemLimit() types.String { return m.MemLimit }
func (m *resourceGroupResourceModel) GetConcurrencyLimit() types.Int64 { return m.ConcurrencyLimit }
func (m *resourceGroupResourceModel) GetBigQueryMemLimit() types.String { return m.BigQueryMemLimit }
func (m *resourceGroupResourceModel) GetBigQueryScanRowsLimit() types.Int64 { return m.BigQueryScanRowsLimit }
func (m *resourceGroupResourceModel) GetBigQueryCPUSecondLimit() types.Int64 { return m.BigQueryCPUSecondLimit }
func (m *resourceGroupResourceModel) GetClassifiers() types.Set { return m.Classifiers }

type classifierModel struct {
	User   types.String  `tfsdk:"user"`
	Role   types.String  `tfsdk:"role"`
	QueryType types.String `tfsdk:"query_type"`
	SourceIP  types.String `tfsdk:"source_ip"`
	DB        types.String `tfsdk:"db"`
}

func (r *resourceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_group"
}

func (r *resourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":                          schema.StringAttribute{Required: true},
			"cpu_core_limit":                schema.Int64Attribute{Optional: true},
			"mem_limit":                     schema.StringAttribute{Optional: true},
			"concurrency_limit":             schema.Int64Attribute{Optional: true},
			"big_query_mem_limit":           schema.StringAttribute{Optional: true},
			"big_query_scan_rows_limit":     schema.Int64Attribute{Optional: true},
			"big_query_cpu_second_limit":    schema.Int64Attribute{Optional: true},
			"classifiers": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"user":       schema.StringAttribute{Optional: true},
						"role":       schema.StringAttribute{Optional: true},
						"query_type": schema.StringAttribute{Optional: true},
						"source_ip":  schema.StringAttribute{Optional: true},
						"db":         schema.StringAttribute{Optional: true},
					},
				},
			},
		},
	}
}

func (r *resourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.CreateResourceGroup(&plan); err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource Group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rg, err := r.client.GetResourceGroup(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource group", err.Error())
		return
	}

	state.Name = rg.Name
	state.CPUCoreLimit = rg.CPUCoreLimit
	state.MemLimit = rg.MemLimit
	state.ConcurrencyLimit = rg.ConcurrencyLimit
	state.BigQueryMemLimit = rg.BigQueryMemLimit
	state.BigQueryScanRowsLimit = rg.BigQueryScanRowsLimit
	state.BigQueryCPUSecondLimit = rg.BigQueryCPUSecondLimit
	state.Classifiers = rg.Classifiers
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateResourceGroup(&plan, &state); err != nil {
		resp.Diagnostics.AddError("Failed to update resource group", err.Error())
		return
	}

	resp.State.Set(ctx, &plan)
}

func (r *resourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteResourceGroup(state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Unable to Delete Resource Group", err.Error())
	}
}

func (r *resourceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.State.SetAttribute(ctx, path.Root("name"), req.ID)
}

func (r *resourceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}

	r.client = client
}
