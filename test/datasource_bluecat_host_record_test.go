// Copyright 2021 BlueCat Networks. All rights reserved

package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHostRecord(t *testing.T) {
	// getting with full field
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceHostRecordsReadFullField,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_host_record.%s", hostDataSource1), "zone", zone),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_host_record.%s", hostDataSource1), "fqdn", fmt.Sprintf("%s.%s", fqdnName1, zone)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_host_record.%s", hostDataSource1), "ip_address", ipAddress1),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_host_record.%s", hostDataSource1), "ttl", "200"),
				),
			},
		},
	})
	// getting without some optional fields
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceHostRecordsReadNotFullField,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_host_record.%s", hostDataSource2), "zone", zone),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_host_record.%s", hostDataSource2), "fqdn", fqdnName2),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_host_record.%s", hostDataSource2), "ip_address", ipAddress2),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_host_record.%s", hostDataSource2), "ttl", "200"),
				),
			},
		},
	})
}

var dataHostResource1 = "host_record_1"
var hostDataSource1 = "test_host_record_1"
var fqdnName1 = "host1"
var ipAddress1 = "1.1.0.30"
var testAccDataSourceHostRecordsReadFullField = fmt.Sprintf(
	`%s
	resource "bluecat_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		ip_address = "%s"
		ttl = 200
		properties = ""
		depends_on = [bluecat_zone.sub_zone_test, bluecat_ipv4network.network_test]
		}

	data "bluecat_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		fqdn = bluecat_host_record.%s.absolute_name
		ip_address = bluecat_host_record.%s.ip_address
		}`, GetTestEnvResources(), dataHostResource1, configuration, view, zone, fqdnName1, ipAddress1, hostDataSource1, configuration,
	view, zone, dataHostResource1, dataHostResource1)

var dataHostResource2 = "host_record_2"
var hostDataSource2 = "test_host_record_2"
var fqdnName2 = "host2.example.com"
var ipAddress2 = "1.1.0.31"
var testAccDataSourceHostRecordsReadNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		ip_address = "%s"
		ttl = 200
		properties = ""
		depends_on = [bluecat_zone.sub_zone_test, bluecat_ipv4network.network_test]
		}

	data "bluecat_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		fqdn = bluecat_host_record.%s.absolute_name
		ip_address = bluecat_host_record.%s.ip_address
		}`, GetTestEnvResources(), dataHostResource2, configuration, view, zone, fqdnName2, ipAddress2, hostDataSource2, configuration, view, dataHostResource2, dataHostResource2)
