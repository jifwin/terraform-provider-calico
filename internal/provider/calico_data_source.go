package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/projectcalico/api/pkg/client/clientset_generated/clientset"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCoffeesDataSource() datasource.DataSource {
	return &coffeesDataSource{}
}

type coffeesDataSource struct {
	client *clientset.Clientset
}

// TODO: rename
type GlobalNetworkPoliciesDataSourceModel struct {
	Policies []globalNetworkPolicyModel `tfsdk:"policies"`
}

type globalNetworkPolicyModel struct {
	Name types.String `tfsdk:"name"`
}

// TODO: rename the datasource
// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &coffeesDataSource{}
	_ datasource.DataSourceWithConfigure = &coffeesDataSource{}
)

func (d *coffeesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_coffees"
}

// Configure adds the provider configured client to the data source.
func (d *coffeesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clientset.Clientset)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			//TODO: update
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *coffeesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"policies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *coffeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GlobalNetworkPoliciesDataSourceModel
	//TODO: define the state model and set the state on it

	//TODO: what is this contxt?
	//TODO: replace the context

	//TODO: finish implementing and create a kubeconfig to test

	//TODO: replace
	//result, err := d.client.ProjectcalicoV3().GlobalNetworkPolicies().List(context.Background(), v1.ListOptions{})
	result, err := d.client.ProjectcalicoV3().KubeControllersConfigurations().List(context.Background(), v1.ListOptions{})

	//TODO:, error, not panic?
	if err != nil {
		panic(err)
	}

	for _, item := range result.Items {
		data := globalNetworkPolicyModel{Name: types.StringValue(item.ObjectMeta.Name)}
		state.Policies = append(state.Policies, data)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//TODO: implement enable the patch resource
	//TODO: EnableGlobalNetworkPoliciesPatch resource?
}
