// Copyright 2021 BlueCat Networks. All rights reserved

package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceCNAMERecord(t *testing.T) {
	// getting with full field
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDataSourceCNAMERecordsReadFullField,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_cname_record.%s", cnameDataSource1), "zone", zone),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_cname_record.%s", cnameDataSource1), "linked_record", linkedRecord1),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_cname_record.%s", cnameDataSource1), "canonical", fmt.Sprintf("%s.%s", canonicalName1, zone)),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_cname_record.%s", cnameDataSource1), "ttl", "200"),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// getting without some optional fields
			resource.TestStep{
				Config: testAccDataSourceCNAMERecordsReadNotFullField,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_cname_record.%s", cnameDataSource2), "zone", zone),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_cname_record.%s", cnameDataSource2), "linked_record", linkedRecord2),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_cname_record.%s", cnameDataSource2), "canonical", canonicalName2),
					resource.TestCheckResourceAttr(fmt.Sprintf("data.bluecat_cname_record.%s", cnameDataSource2), "ttl", "200"),
				),
			},
		},
	})
}

var resourceHostCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_host_record" "host_record_1" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "host1.example.com"
		ip4_address = "1.1.0.9"
		ttl = 200
		properties = ""
	}`, server, configuration, view, zone)

var dataCNAMEResource1 = "cname_record_1"
var cnameDataSource1 = "test_cname_record_1"
var canonicalName1 = "cname1"
var linkedRecord1 = "host1.example.com"
var testAccDataSourceCNAMERecordsReadFullField = fmt.Sprintf(
	`%s
	resource "bluecat_cname_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		linked_record = "%s"
		ttl = 200
		properties = ""
  		depends_on = [bluecat_host_record.host_record_1]		
	}

	data "bluecat_cname_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		linked_record = bluecat_cname_record.%s.linked_record
		canonical = bluecat_cname_record.%s.absolute_name
		}`, resourceHostCreateFullField, dataCNAMEResource1, configuration, view, zone, canonicalName1, linkedRecord1,
	cnameDataSource1, configuration, view, zone, dataCNAMEResource1, dataCNAMEResource1)

var resourceHostCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_host_record" "host_record_2" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "host2.example.com"
		ip4_address = "1.1.0.9"
		ttl = 200
		properties = ""
	}`, server, configuration, view, zone)

var dataCNAMEResource2 = "cname_record_2"
var cnameDataSource2 = "test_cname_record_2"
var canonicalName2 = "cname2.example.com"
var linkedRecord2 = "host2.example.com"
var testAccDataSourceCNAMERecordsReadNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_cname_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		linked_record = "%s"
		ttl = 200
		properties = ""
		depends_on = [bluecat_host_record.host_record_2]
		}

	data "bluecat_cname_record" "%s" {
		configuration = "%s"
		view = "%s"
		linked_record = bluecat_cname_record.%s.linked_record
		canonical = bluecat_cname_record.%s.absolute_name
		}`, resourceHostCreateNotFullField, dataCNAMEResource2, configuration, view, zone, canonicalName2, linkedRecord2,
	cnameDataSource2, configuration, view, dataCNAMEResource2, dataCNAMEResource2)
