// Copyright 2021 BlueCat Networks. All rights reserved

package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceZoneRecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceZoneRead,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_zone.%s", zoneDataSource), "zone", zoneName),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_zone.%s", zoneDataSource), "deployable", zoneDeployable),
				),
			},
		},
	})
}

var zoneDataSource = "test_zone"
var zoneName = "subzone_test.com"
var zoneDeployable = "True"

var testAccDataSourceZoneRead = fmt.Sprintf(
	`%s
	data "bluecat_zone" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
	}`, server, zoneDataSource, configuration, view, zoneName)
