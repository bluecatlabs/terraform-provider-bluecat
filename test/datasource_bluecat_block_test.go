// Copyright 2021 BlueCat Networks. All rights reserved

package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBlockRecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlockRead,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_ipv4block.%s", ipBlockDataSource), "cidr", cidrBlock),
				),
			},
		},
	})
}

var ipBlockDataSource = "test_ip4block"
var cidrBlock = "1.1.0.0/16"
var testAccDataSourceBlockRead = fmt.Sprintf(
	`%s
	data "bluecat_ipv4block" "%s" {
		configuration = "%s"
		cidr = "%s"
	}`, server, ipBlockDataSource, configuration, cidrBlock)
