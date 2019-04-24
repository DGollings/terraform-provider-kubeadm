package kubeadm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provisioner() terraform.ResourceProvisioner {
	return &schema.Provisioner{
		// ConnSchema: map[string]*schema.Schema{
		// 	"address": {
		// 		Type:     schema.TypeString,
		// 		Optional: true,
		// 	},
		// },

		Schema: map[string]*schema.Schema{
			"config": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "init/join configuration to use",
			},
			"join": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "seeder node to join. Or start a seeder when not provided",
			},
			// not sure really necessary: maybe we can get Changes('count'):
			"remove": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "when true, remove this node from the cluster instead of addinng it",
			},
			"nodename": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "node name in the kubernetes cluster",
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "local kubeconfig file where the remote admin config file will be copied to",
			},
			"prevent_sudo": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "prevent the use of sudo",
			},
			"ignore_checks": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "list of preflight checks to ignore by kubeadm",
			},
			"manifests": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "list of manifests to load in the API server once the master is setup",
			},
			"install": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "try to automatically install kubeadm with the built-in helper script",
						},
						"script": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "user-provided installation script",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "kubeadm version to install.",
						},
					},
				},
			},
		},

		ApplyFunc: applyFn,
	}
}
