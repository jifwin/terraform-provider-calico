// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"flag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/projectcalico/api/pkg/client/clientset_generated/clientset"
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

func (p *CalicoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CalicoProviderModel

	//TODO: should verify connection to the cluster, etc
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new config based on kubeconfig file.
	var kubeconfig *string //TODO: read from context
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse() //TODO: read from context
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	//TODO: read from content, not from file
	// Build a clientset based on the provided kubeconfig file.
	cs, err := clientset.NewForConfig(config) //TODO: pack this into resp?
	if err != nil {
		panic(err)
	}
	resp.ResourceData = cs
	resp.DataSourceData = cs

	//TODO: implement the first resource
	//TODO: move out
	_, err = cs.ProjectcalicoV3().GlobalNetworkPolicies().List(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	//TODo: verify that patch
	//TODO: implement enable the patch resource
	//TODO: EnableGlobalNetworkPoliciesPatch resource?

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
