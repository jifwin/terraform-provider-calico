// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/projectcalico/api/pkg/client/clientset_generated/clientset"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &CalicoProvider{}
var _ provider.ProviderWithFunctions = &CalicoProvider{}

// CalicoProvider defines the provider implementation.
type CalicoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CalicoProviderModel describes the provider data model.
type CalicoProviderModel struct {
	Kubeconfig types.String `tfsdk:"kubeconfig"`
}

func (p *CalicoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "calico"
	resp.Version = p.version
}

func (p *CalicoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"kubeconfig": schema.StringAttribute{ //TODO: consider other options to authenticate to the cluster
				//TODO: maybe split kubeconfig to token, host, etc
				MarkdownDescription: "kubeconfig",
				Required:            true,
			},
		},
	}
}

// TODO: should verify connection to the cluster, etc
func (p *CalicoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config CalicoProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//TODO: verify and refactor
	if config.Kubeconfig.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("kubeconfig"),
			"Unknown kubeconfig",
			"Kubeconfig not configured",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	kubeconfigGetter := func() (*clientcmdapi.Config, error) {
		var kubeconfig = config.Kubeconfig.ValueString() //TODO: ValueString vs String
		return clientcmd.Load([]byte(kubeconfig))
	}

	//TODO: check that empty string arg
	client, err := clientcmd.BuildConfigFromKubeconfigGetter("", kubeconfigGetter)

	//TODO: proper err handling via diagnostics?
	if err != nil {
		panic([]byte(config.Kubeconfig.String()))
		//panic(err.Error())
	}

	clientset, err := clientset.NewForConfig(client)

	if err != nil {
		panic(err.Error())
	}

	resp.ResourceData = clientset
	resp.DataSourceData = clientset
}

// TODO: how to wait for calico resources to be available
func (p *CalicoProvider) Resources(ctx context.Context) []func() resource.Resource {
	//TODO: move to another file
	//TODO: clientset in configuration, use here
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *CalicoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
		NewCoffeesDataSource,
	}
}

func (p *CalicoProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CalicoProvider{
			version: version,
		}
	}
}
