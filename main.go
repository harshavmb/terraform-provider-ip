package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/kekwork/terraform-provider-ip/ip"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ip.Provider})
}
