package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-starrocks/internal/client"
)

var _ provider.Provider = &starrocksProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &starrocksProvider{version: version}
	}
}

type starrocksProvider struct {
	version string
}

type starrocksProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *starrocksProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "starrocks"
	resp.Version = p.version
}

func (p *starrocksProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *starrocksProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config starrocksProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("STARROCKS_HOST")
	username := os.Getenv("STARROCKS_USERNAME")
	password := os.Getenv("STARROCKS_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(path.Root("host"), "Missing StarRocks Host", "Set host in configuration or STARROCKS_HOST environment variable")
	}
	if username == "" {
		resp.Diagnostics.AddAttributeError(path.Root("username"), "Missing StarRocks Username", "Set username in configuration or STARROCKS_USERNAME environment variable")
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(path.Root("password"), "Missing StarRocks Password", "Set password in configuration or STARROCKS_PASSWORD environment variable")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := client.NewClient(host, username, password)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create StarRocks Client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *starrocksProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *starrocksProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewResourceGroupResource,
	}
}
