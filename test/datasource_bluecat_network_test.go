// Copyright 2021 BlueCat Networks. All rights reserved

package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceNetworkRecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceNetworkRead,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_ipv4network.%s", ipNetworkDataSource), "cidr", cidrNetwork),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_ipv4network.%s", ipNetworkDataSource), "name", name),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_ipv4network.%s", ipNetworkDataSource), "gateway", gateway),
				),
			},
		},
	})
}

var ipNetworkDataSource = "test_ip4network"
var name = "network"
var cidrNetwork = "1.1.0.0/16"
var gateway = "1.1.0.1"
var testAccDataSourceNetworkRead = fmt.Sprintf(
	`%s
	data "bluecat_ipv4network" "%s" {
		configuration = "%s"
		cidr = "%s"
		}`, server, ipNetworkDataSource, configuration, cidrNetwork)
