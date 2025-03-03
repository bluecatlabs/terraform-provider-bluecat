// Copyright 2020 BlueCat Networks. All rights reserved

package main

import (
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"terraform-provider-bluecat/bluecat"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return bluecat.Provider()
		},
		ProviderAddr: "terraform-provider-bluecat/bluecat",
		Debug:        debugMode,
	}

	plugin.Serve(opts)
}
