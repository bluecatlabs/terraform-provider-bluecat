// Copyright 2021 BlueCat Networks. All rights reserved

package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceZoneRecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
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
var zoneName = "test_zone_1.com"
var zoneDeployable = "True"

var testAccDataSourceZoneRead = fmt.Sprintf(
	`%s
	resource "bluecat_zone" "test_zone_1" {
		configuration = "%s"
		view = "test"
		zone = "%s"
		deployable = "%s"
		properties = ""
		depends_on = [bluecat_view.view_test]
		}
	
	data "bluecat_zone" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		deployable = "%s"
		depends_on = [bluecat_view.view_test, bluecat_zone.test_zone_1]
		}`, GetTestEnvResources(), configuration, zoneName, zoneDeployable, zoneDataSource, configuration, view, zoneName, zoneDeployable)
