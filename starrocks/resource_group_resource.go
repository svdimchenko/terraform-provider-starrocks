package starrocks

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	client *Client
}

type resourceGroupResourceModel struct {
	Name                     types.String `tfsdk:"name"`
	CPUWeight                types.Int64  `tfsdk:"cpu_weight"`
	ExclusiveCPUCores        types.Int64  `tfsdk:"exclusive_cpu_cores"`
	CPUCoreLimit             types.Int64  `tfsdk:"cpu_core_limit"`
	MaxCPUCores              types.Int64  `tfsdk:"max_cpu_cores"`
	MemLimit                 types.String `tfsdk:"mem_limit"`
	ConcurrencyLimit         types.Int64  `tfsdk:"concurrency_limit"`
	BigQueryMemLimit         types.String `tfsdk:"big_query_mem_limit"`
	BigQueryScanRowsLimit    types.Int64  `tfsdk:"big_query_scan_rows_limit"`
	BigQueryCPUSecondLimit   types.Int64  `tfsdk:"big_query_cpu_second_limit"`
	Classifiers              types.List   `tfsdk:"classifiers"`
}

func (m *resourceGroupResourceModel) GetName() types.String { return m.Name }
func (m *resourceGroupResourceModel) GetCPUWeight() types.Int64 { return m.CPUWeight }
func (m *resourceGroupResourceModel) GetExclusiveCPUCores() types.Int64 { return m.ExclusiveCPUCores }
func (m *resourceGroupResourceModel) GetCPUCoreLimit() types.Int64 { return m.CPUCoreLimit }
func (m *resourceGroupResourceModel) GetMaxCPUCores() types.Int64 { return m.MaxCPUCores }
func (m *resourceGroupResourceModel) GetMemLimit() types.String { return m.MemLimit }
func (m *resourceGroupResourceModel) GetConcurrencyLimit() types.Int64 { return m.ConcurrencyLimit }
func (m *resourceGroupResourceModel) GetBigQueryMemLimit() types.String { return m.BigQueryMemLimit }
func (m *resourceGroupResourceModel) GetBigQueryScanRowsLimit() types.Int64 { return m.BigQueryScanRowsLimit }
func (m *resourceGroupResourceModel) GetBigQueryCPUSecondLimit() types.Int64 { return m.BigQueryCPUSecondLimit }
func (m *resourceGroupResourceModel) GetClassifiers() types.List { return m.Classifiers }

type classifierModel struct {
	User   types.String  `tfsdk:"user"`
	Role   types.String  `tfsdk:"role"`
	QueryType types.String `tfsdk:"query_type"`
	SourceIP  types.String `tfsdk:"source_ip"`
	DB        types.String `tfsdk:"db"`
}

func formatMemLimit(s string) string {
	if len(s) > 0 && s[len(s)-1] == '%' {
		val := s[:len(s)-1]
		if !strings.Contains(val, ".") {
			return val + ".0%"
		}
	}
	return s
}

type memLimitModifier struct{}

func (m memLimitModifier) Description(ctx context.Context) string {
	return "Formats mem_limit to match StarRocks format"
}

func (m memLimitModifier) MarkdownDescription(ctx context.Context) string {
	return "Formats mem_limit to match StarRocks format"
}

func (m memLimitModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	resp.PlanValue = types.StringValue(formatMemLimit(req.ConfigValue.ValueString()))
}

func (r *resourceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_group"
}

func (r *resourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":                          schema.StringAttribute{Required: true},
			"cpu_weight":                    schema.Int64Attribute{Optional: true},
			"exclusive_cpu_cores":           schema.Int64Attribute{Optional: true},
			"cpu_core_limit":                schema.Int64Attribute{Optional: true},
			"max_cpu_cores":                 schema.Int64Attribute{Optional: true},
			"mem_limit":                     schema.StringAttribute{Optional: true, PlanModifiers: []planmodifier.String{memLimitModifier{}, stringplanmodifier.UseStateForUnknown()}},
			"concurrency_limit":             schema.Int64Attribute{Optional: true},
			"big_query_mem_limit":           schema.StringAttribute{Optional: true},
			"big_query_scan_rows_limit":     schema.Int64Attribute{Optional: true},
			"big_query_cpu_second_limit":    schema.Int64Attribute{Optional: true},
			"classifiers": schema.ListNestedAttribute{
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

	// Update only the fields that GetResourceGroup returns, keep classifiers from state
	state.Name = rg.Name
	if !rg.CPUWeight.IsNull() {
		state.CPUWeight = rg.CPUWeight
	}
	if !rg.ExclusiveCPUCores.IsNull() {
		state.ExclusiveCPUCores = rg.ExclusiveCPUCores
	}
	if !rg.CPUCoreLimit.IsNull() {
		state.CPUCoreLimit = rg.CPUCoreLimit
	}
	if !rg.MaxCPUCores.IsNull() {
		state.MaxCPUCores = rg.MaxCPUCores
	}
	// Keep mem_limit from state to avoid drift from "80%" vs "80.0%"
	if !rg.ConcurrencyLimit.IsNull() {
		state.ConcurrencyLimit = rg.ConcurrencyLimit
	}
	if !rg.BigQueryMemLimit.IsNull() {
		state.BigQueryMemLimit = rg.BigQueryMemLimit
	}
	if !rg.BigQueryScanRowsLimit.IsNull() {
		state.BigQueryScanRowsLimit = rg.BigQueryScanRowsLimit
	}
	if !rg.BigQueryCPUSecondLimit.IsNull() {
		state.BigQueryCPUSecondLimit = rg.BigQueryCPUSecondLimit
	}
	// Keep classifiers from existing state since GetResourceGroup doesn't return them properly
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete and recreate
	if err := r.client.DeleteResourceGroup(plan.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Unable to Delete Resource Group", err.Error())
		return
	}

	if err := r.client.CreateResourceGroup(&plan); err != nil {
		resp.Diagnostics.AddError("Unable to Create Resource Group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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

	c, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *Client, got: %T", req.ProviderData))
		return
	}

	r.client = c
}
