package main

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceExternalHostRecord(t *testing.T) {
	// create with full fields and update
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckExternalHostRecordDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccresourceExternalHostRecordCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccExternalHostRecordExists(t, fmt.Sprintf("bluecat_external_host_record.%s", externalHostResource1), externalHostName1, externalHostAddresses1),
				),
			},
			// update
			{
				Config: testAccresourceExternalHostRecordUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccExternalHostRecordExists(t, fmt.Sprintf("bluecat_external_host_record.%s", externalHostResource1), externalHostName1, externalHostAddresses2),
				),
			},
		},
	})
	// create without some optional fields
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckExternalHostRecordDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccresourceExternalHostRecordCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccExternalHostRecordExists(t, fmt.Sprintf("bluecat_external_host_record.%s", externalHostResource1), externalHostName1, externalHostAddresses1),
				),
			},
		},
	})
}

func testAccCheckExternalHostRecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bluecat_external_host_record" {
			msg := fmt.Sprintf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
			log.Error(msg)
			// return fmt.Errorf(msg)
		}
		_, err := objMgr.GetExternalHostRecord(configuration, view, rs.Primary.ID)
		if err == nil {
			msg := fmt.Sprintf("External Host record %s is not removed", rs.Primary.ID)
			log.Error(msg)
			return fmt.Errorf("External Host record %s is not removed", rs.Primary.ID)
		}
	}
	return nil
}

func testAccExternalHostRecordExists(t *testing.T, resourceName string, name string, addresses string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		externalHostRecord, err := objMgr.GetExternalHostRecord(configuration, view, name)
		if err != nil {
			msg := fmt.Sprintf("Getting External Host record %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf("Getting External Host record %s failed: %s", rs.Primary.ID, err)
		}
		effectiveAddresses := utils.GetPropertyValue("addresses", externalHostRecord.Properties)
		if effectiveAddresses == "" {
			effectiveAddresses = externalHostRecord.Addresses
		}
		if effectiveAddresses == "" {
			effectiveAddresses = rs.Primary.Attributes["addresses"]
		}
		if effectiveAddresses != addresses {
			msg := fmt.Sprintf("Getting External Host record %s failed: expect addresses=%s, but received properties='%s' addresses='%s'", rs.Primary.ID, addresses, externalHostRecord.Properties, externalHostRecord.Addresses)
			log.Error(msg)
			return fmt.Errorf("Getting External Host record %s failed: expect addresses=%s, but received properties='%s' addresses='%s'", rs.Primary.ID, addresses, externalHostRecord.Properties, externalHostRecord.Addresses)
		}

		return nil
	}
}

var externalHostResource1 = "external_host_record_a1"
var externalHostName1 = "ext-a1.example.com"
var externalHostAddresses1 = "1.1.0.21"
var externalHostProperties1 = ""
var testAccresourceExternalHostRecordCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_external_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		absolute_name = "%s"
		addresses = "%s"
		properties = "%s"
		depends_on = [bluecat_zone.sub_zone_test]
	}`, GetTestEnvResources(), externalHostResource1, configuration, view, externalHostName1, externalHostAddresses1, externalHostProperties1)

var testAccresourceExternalHostRecordCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_external_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		absolute_name = "%s"
		addresses = "%s"
		depends_on = [bluecat_zone.sub_zone_test]
	}`, GetTestEnvResources(), externalHostResource1, configuration, view, externalHostName1, externalHostAddresses1)

var externalHostAddresses2 = "1.1.0.22"
var externalHostProperties2 = ""
var testAccresourceExternalHostRecordUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_external_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		absolute_name = "%s"
		addresses = "%s"
		properties = "%s"
		depends_on = [bluecat_zone.sub_zone_test]
	}`, GetTestEnvResources(), externalHostResource1, configuration, view, externalHostName1, externalHostAddresses2, externalHostProperties2)
