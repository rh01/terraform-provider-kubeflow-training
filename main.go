package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/rh01/terraform-provider-kubeflow-training/kubeflowtraining"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kubeflowtraining.Provider})
}
