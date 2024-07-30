package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-trust-store/tools"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &trustStoreProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &trustStoreProvider{
			version: version,
		}
	}
}

// hashicupsProvider is the provider implementation.
type trustStoreProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type trusStoreProviderModel struct {
	TOKEN types.String `tfsdk:"token"`
	URL   types.String `tfsdk:"url"`
}

// Metadata returns the provider type name.
func (p *trustStoreProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opaas_trust_store"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *trustStoreProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Required: true,
			},
			"url": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *trustStoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config trusStoreProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustStore := tools.TrustStoreClient{
		Token: config.TOKEN.ValueString(),
		Url:   config.URL.ValueString(),
	}

	// Make the client available during DataSource and Resource
	// type Configure methods.

	resp.DataSourceData = &trustStore
	resp.ResourceData = &trustStore
}

// DataSources defines the data sources implemented in the provider.
func (p *trustStoreProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *trustStoreProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTrustStoreResource,
	}
}
