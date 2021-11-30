// Copyright 2020 BlueCat Networks. All rights reserved

package main

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourcePTRRecord(t *testing.T) {
	// create with full fields and update
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPTRRecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourcePTRRecordCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccPTRRecordExists(t, fmt.Sprintf("bluecat_ptr_record.%s", ptrResource1), ptrName1, ptrTTL1, ptrIP1),
				),
			},
			// update
			resource.TestStep{
				Config: testAccresourcePTRRecordUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccPTRRecordExists(t, fmt.Sprintf("bluecat_ptr_record.%s", ptrResource1), ptrName1, ptrTTL1, ptrIP2),
				),
			},
		},
	})
	// create without some optional fields
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPTRRecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourcePTRRecordCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccPTRRecordExists(t, fmt.Sprintf("bluecat_ptr_record.%s", ptrResource1), ptrName1, ptrTTL1, ptrIP1),
				),
			},
		},
	})
}

func testAccCheckPTRRecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bluecat_ptr_record" {
			msg := fmt.Sprintf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		_, err := objMgr.GetHostRecord(configuration, view, rs.Primary.ID)
		if err == nil {
			msg := fmt.Sprintf("Host record %s is not removed", rs.Primary.ID)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
	}
	return nil
}

func testAccPTRRecordExists(t *testing.T, resource string, name string, ttl string, ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check ptr on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		hostRecord, err := objMgr.GetHostRecord(configuration, view, name)
		if err != nil {
			msg := fmt.Sprintf("Getting Host record %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		if checkValidHostRecord(*hostRecord, ttl, ip, "true") == false {
			msg := fmt.Sprintf("Getting Host record %s failed: %s. Expect ttl=%s addresses=%s reverseRecord=%s in properties, but received '%s'", rs.Primary.ID, err, ttl, ip, "true", hostRecord.Properties)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}
}

var ptrResource1 = "ptr_record"
var ptrName1 = "host.example.com"
var ptrNet1 = "1.1.0.0/16"
var ptrIP1 = "1.1.0.5"
var ptrTTL1 = "1000"
var ptrProperties1 = ""
var testAccresourcePTRRecordCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ptr_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "%s"
		network = "%s"
		ip4_address = "%s"
		ttl = %s
		properties = "%s"
		reverse_record = "True"
	  }`, server, ptrResource1, configuration, view, zone, ptrName1, ptrNet1, ptrIP1, ptrTTL1, ptrProperties1)

var testAccresourcePTRRecordCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ptr_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "%s"
		network = "%s"
		ip4_address = "%s"
		ttl = %s
		reverse_record = "True"
		}`, server, ptrResource1, configuration, view, zone, ptrName1, ptrNet1, ptrIP1, ptrTTL1)

var ptrIP2 = "1.1.0.6"
var testAccresourcePTRRecordUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ptr_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "%s"
		network = "%s"
		ip4_address = "%s"
		ttl = %s
		properties = "%s"
		reverse_record = "True"
		}`, server, ptrResource1, configuration, view, zone, ptrName1, ptrNet1, ptrIP2, ptrTTL1, ptrProperties1)
