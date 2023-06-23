// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	DefaultConcurrency = 30
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &K0sProvider{}

// K0sProvider defines the provider implementation.
type K0sProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// K0sProviderModel describes the provider data model.
type K0sProviderModel struct {
	DisableTelemetry types.Bool `tfsdk:"disable_telemetry"`
	Concurrency types.Int64 `tfsdk:"concurrency"`
}

func (p *K0sProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "scaffolding"
	resp.Version = p.version
}

func (p *K0sProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"disable_telemetry": schema.BoolAttribute{
				MarkdownDescription: "Disable the anonymous telemetry.",
				Optional:            true,
			},
			"concurrency": schema.Int64Attribute{
				MarkdownDescription: "How many machines to process in parrallel in apply and reset operations",
				Optional: true,
		},
	}
}

func (p *K0sProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data K0sProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = data
	resp.ResourceData = data
}

func (p *K0sProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *K0sProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &K0sProvider{
			version: version,
		}
	}
}
