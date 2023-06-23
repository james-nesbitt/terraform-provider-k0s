package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"	
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1"
)

// V1Beta1Schema Terraform schema for k0sctl
func V1Beta1Schema() schema.Schema {
	return schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "K0S installation using k0sctl, parametrized",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"skip_destroy": schema.BoolAttribute{
				MarkdownDescription: "Do not bother uninstalling on destroy",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},

			"kubeconfig": schema.StringAttribute{
				MarkdownDescription: "Output kubeconfig file contents for the cluster.",
				Computed:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"metadata": schema.SingleNestedBlock{
				MarkdownDescription: "Metadata for the launchpad cluster",

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "Cluster name",
						Required:            true,
					},
				},
			},

			"spec": schema.SingleNestedBlock{
				MarkdownDescription: "Launchpad install specifications",

				Blocks: map[string]schema.Block{

					"k0s": schema.SingleNestedBlock{
						MarkdownDescription: "K0S installation configuration",

						Attributes: map[string]schema.Attribute{
							"version": schema.StringAttribute{
								MarkdownDescription: "MKE version to install",
								Required:            true,
							},
							"image_repo": schema.StringAttribute{
								MarkdownDescription: "Image repo for MKE images",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("docker.io/mirantis"),
							},
							"admin_username": schema.StringAttribute{
								MarkdownDescription: "MKE admin user name",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("admin"),
							},
							"admin_password": schema.StringAttribute{
								MarkdownDescription: "MKE admin user password",
								Required:            true,
								Sensitive:           true,
							},

							"install_flags": schema.ListAttribute{
								MarkdownDescription: "Optional MKE bootstrapper install flags",
								ElementType:         types.StringType,
								Optional:            true,
								Computed:            true,
							},
							"upgrade_flags": schema.ListAttribute{
								MarkdownDescription: "Optional MKE bootstrapper update flags",
								ElementType:         types.StringType,
								Optional:            true,
								Computed:            true,
							},
						},
					},
				},
			},
		},
	}
}

// Data structure to mimic the K0S config structure
type K0SConfig_V1Beta1 struct {
	Id types.String `tfsdk:"id"`
}

func (kc *K0SConfig_V1Beta1) Cluster() v1beta1.Cluster {
	return v1beta1.Cluster{}
}
