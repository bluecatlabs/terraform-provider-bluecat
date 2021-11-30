package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"
)

func TestAccResourceGenericRecord(t *testing.T) {
	// create with full fields and update
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGenericRecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceGenericRecordCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccGenericRecordExists(t, fmt.Sprintf("bluecat_generic_record.%s", genericResource1), genericType1, genericName1, genericTTL1, genericData1),
				),
			},
			// // update
			resource.TestStep{
				Config: testAccresourceGenericRecordUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccGenericRecordExists(t, fmt.Sprintf("bluecat_generic_record.%s", genericResource1), genericType1, genericName1, genericTTL2, genericData2),
				),
			},
		},
	})
	// create without some optional fields
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGenericRecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceGenericRecordCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccGenericRecordExists(t, fmt.Sprintf("bluecat_generic_record.%s", genericResource1), genericType1, genericName1, genericTTL1, genericData1),
				),
			},
		},
	})
}

func testAccCheckGenericRecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "bluecat_generic_record" {
			_, err := objMgr.GetGenericRecord(configuration, view, rs.Primary.ID)
			if err == nil {
				msg := fmt.Sprintf("Generic record %s is not removed", rs.Primary.ID)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
		} else {
			msg := fmt.Sprintf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
	}
	return nil
}

func testAccGenericRecordExists(t *testing.T, resource string, typerr string, name string, ttl string, data string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check TXT record on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		genericRecord, err := objMgr.GetGenericRecord(configuration, view, name)
		if err != nil {
			msg := fmt.Sprintf("Getting Generic record %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		ttlProperty := getPropertyValue("ttl", genericRecord.Properties)
		dataProperty := getPropertyValue("rdata", genericRecord.Properties)
		typeProperty := getPropertyValue("type", genericRecord.Properties)

		if ttlProperty != ttl || dataProperty != data || typeProperty != typerr {
			msg := fmt.Sprintf("Getting Generic record %s failed: %s. Expect ttl=%s data=%s type=%s in properties, but received '%s'", rs.Primary.ID, err, ttl, data, typerr, genericRecord.Properties)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}
}

var genericResource1 = "generic_record_a4"
var genericName1 = "a4.example.com"
var genericData1 = "test"
var genericType1 = "NS"
var genericTTL1 = "400"
var genericProperties1 = ""
var testAccresourceGenericRecordCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_generic_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
        type = "%s"
		absolute_name = "%s"
		data = "%s"
		ttl = %s
		properties = "%s"
	  }`, server, genericResource1, configuration, view, zone, genericType1, genericName1, genericData1, genericTTL1, genericProperties1)

var testAccresourceGenericRecordCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_generic_record" "%s" {
		configuration = "%s"
		view = "%s"
		type = "%s"
		absolute_name = "%s"
		data = "%s"
		ttl = %s
		properties = "%s"
		}`, server, genericResource1, configuration, view, genericType1, genericName1, genericData1, genericTTL1, genericProperties1)

var genericData2 = "test2"
var genericTTL2 = "4000"
var genericProperties2 = ""
var testAccresourceGenericRecordUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_generic_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		type = "%s"
		absolute_name = "%s"
		data = "%s"
		ttl = %s
		properties = "%s"
		}`, server, genericResource1, configuration, view, zone, genericType1, genericName1, genericData2, genericTTL2, genericProperties2)

func TestAccResourceA4Record(t *testing.T) {
	// create with full fields and update
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckA4RecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceA4RecordCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccA4RecordExists(t, fmt.Sprintf("bluecat_generic_record.%s", genericResource3), genericType3, genericName3, genericTTL3, genericData3),
				),
			},
			// // update
			resource.TestStep{
				Config: testAccresourceA4RecordUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccA4RecordExists(t, fmt.Sprintf("bluecat_generic_record.%s", genericResource3), genericType3, genericName3, genericTTL4, genericData4),
				),
			},
		},
	})
	// create without some optional fields
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckA4RecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceA4RecordCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccA4RecordExists(t, fmt.Sprintf("bluecat_generic_record.%s", genericResource3), genericType3, genericName3, genericTTL3, genericData3),
				),
			},
		},
	})
}

func testAccCheckA4RecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "bluecat_generic_record" {
			_, err := objMgr.GetGenericRecord(configuration, view, rs.Primary.ID)
			if err == nil {
				msg := fmt.Sprintf("Generic record %s is not removed", rs.Primary.ID)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
		} else {
			msg := fmt.Sprintf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
	}
	return nil
}

func testAccA4RecordExists(t *testing.T, resource string, typerr string, name string, ttl string, data string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check TXT record on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		genericRecord, err := objMgr.GetGenericRecord(configuration, view, name)
		if err != nil {
			msg := fmt.Sprintf("Getting Generic record %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		ttlProperty := getPropertyValue("ttl", genericRecord.Properties)
		dataProperty := getPropertyValue("rdata", genericRecord.Properties)
		typeProperty := getPropertyValue("type", genericRecord.Properties)

		if ttlProperty != ttl || dataProperty != data || typeProperty != typerr {
			msg := fmt.Sprintf("Getting Generic record %s failed: %s. Expect ttl=%s data=%s type=%s in properties, but received '%s'", rs.Primary.ID, err, ttl, data, typerr, genericRecord.Properties)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}
}

var genericResource3 = "generic_record_a4"
var genericName3 = "a4.example.com"
var genericData3 = "ab::123"
var genericType3 = "AAAA"
var genericTTL3 = "400"
var genericProperties3 = ""
var testAccresourceA4RecordCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_generic_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
        type = "%s"
		absolute_name = "%s"
		data = "%s"
		ttl = %s
		properties = "%s"
	  }`, server, genericResource3, configuration, view, zone, genericType3, genericName3, genericData3, genericTTL3, genericProperties3)

var testAccresourceA4RecordCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_generic_record" "%s" {
		configuration = "%s"
		view = "%s"
		type = "%s"
		absolute_name = "%s"
		data = "%s"
		ttl = %s
		properties = "%s"
		}`, server, genericResource3, configuration, view, genericType3, genericName3, genericData3, genericTTL3, genericProperties3)

var genericData4 = "ab::124"
var genericTTL4 = "4000"
var genericProperties4 = ""
var testAccresourceA4RecordUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_generic_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		type = "%s"
		absolute_name = "%s"
		data = "%s"
		ttl = %s
		properties = "%s"
		}`, server, genericResource3, configuration, view, zone, genericType3, genericName3, genericData4, genericTTL4, genericProperties4)
