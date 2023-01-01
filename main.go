package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hugolesta/terraform-provider-gitguardian/gitguardian"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gitguardian.Provider,
	})
}
