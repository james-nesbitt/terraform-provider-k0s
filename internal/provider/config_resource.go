package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/k0sproject/k0sctl/phase"
	"google.golang.org/grpc/balancer/grpclb/state"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ConfigResource{}
var _ resource.ResourceWithImportState = &ConfigResource{}

func NewConfigResource() resource.Resource {
	return &ConfigResource{}
}

// ConfigResource defines the resource implementation.
type ConfigResource struct {
	disable_telemetry bool
}

// ConfigResourceModel
type ConfigResourceModel struct {
	Concurrency int `tfsdk:"concurrency"`
	ConcurrentUploads int `tfsdk:"concurrent_uploads"`
	DisableDowngradeCheck bool `tfsdk:"disable_downgrade_check"`
	RestoreFrom bool `tfsdk:"restore_from"`
	NoDrain bool `tfsdk:"no_drain"`
	NoWait bool `tfsdk:"no_wait"`

	KubeConfigAPIAddress str `tfsdk:"kubeconfig_ali_address"`

	KubeConfig str
	
	K0SConfig_V1Beta1
}

func (r *ConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (r *ConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = V1Beta1Schema()
}

func (r *ConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	provData, ok := req.ProviderData.(*K0sProviderModel)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected K0sProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.disable_telemetry = provData.DisableTelemetry.ValueBool()
}

func (r *ConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ConfigResourceModel
	ctx := context.Background()
	
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cfg := data.Cluster()
	tflog.Trace(context.Background(), "Cluster configuration built from resource block", map[string]interface{}{"config": cfg})
	
	phase.NoWait = data.NoWait

	manager := phase.Manager{Config: cfg, Concurrency: data.Concurrency, ConcurrentUploads: data.ConcurrentUploads}
	lockPhase := &phase.Lock{}

	manager.AddPhase(
		&phase.Connect{},
		&phase.DetectOS{},
		lockPhase,
		&phase.PrepareHosts{},
		&phase.GatherFacts{},
		&phase.DownloadBinaries{},
		&phase.UploadFiles{},
		&phase.ValidateHosts{},
		&phase.GatherK0sFacts{},
		&phase.ValidateFacts{
			SkipDowngradeCheck: data.DisableDowngradeCheck,
		},
		&phase.UploadBinaries{},
		&phase.DownloadK0s{},
		&phase.InstallBinaries{},
		&phase.RunHooks{
			Stage: "before",
			Action: "apply",
		},
		&phase.PrepareArm{},
		&phase.ConfigureK0s{},
		&phase.Restore{
			RestoreFrom: data.RestoreFrom,
		},
		&phase.InitializeK0s{},
		&phase.InstallControllers{},
		&phase.InstallWorkers{},
		&phase.UpgradeControllers{},
		&phase.UpgradeWorkers{
			NoDrain: data.NoDrain,
		},
		&phase.ResetWorkers{
			NoDrain: data.NoDrain,
		},
		&phase.ResetControllers{
			NoDrain:data.NoDrain,
		},
		&phase.RunHooks{Stage: "after", Action: "apply"},
		&phase.GetKubeconfig{
			APIAddress: data.KubeConfigAPIAddress,
		},
		&phase.Unlock{
			Cancel: lockPhase.Cancel,
		},
		&phase.Disconnect{},		
	)

	if err := manager.Run(); err != nil {
		tflog.Error(context.Background(), "K0S install failed", map[string]interface{}{"log": "still trying to figure out how to get this."})
	} else {
		data.Id = "GET this from the cluster config"
		data.KubeConfig = manager.Config.Metadata.Kubeconfigs

		for _, host := range manager.Config.Spec.Hosts {
			if host.Reset {
				tflog.Warn(context.Background(), "There were nodes that got uninstalled during the apply phase. Please remove them from your k0sctl config file")
				break
			}
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)	
}

func (r *ConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ConfigResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ConfigResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ConfigResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *ConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
