package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"
)

func TestAccResourceTXTRecord(t *testing.T) {
	// create with full fields and update
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckTXTRecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceTXTRecordCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccTXTRecordExists(t, fmt.Sprintf("bluecat_txt_record.%s", txtResource1), txtName1, txtTTL1, txtText1),
				),
			},
			// // update
			resource.TestStep{
				Config: testAccresourceTXTRecordUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccTXTRecordExists(t, fmt.Sprintf("bluecat_txt_record.%s", txtResource1), txtName1, txtTTL2, txtText2),
				),
			},
		},
	})
	// create without some optional fields
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckTXTRecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceTXTRecordCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccTXTRecordExists(t, fmt.Sprintf("bluecat_txt_record.%s", txtResource1), txtName1, txtTTL1, txtText1),
				),
			},
		},
	})
}

func testAccCheckTXTRecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "bluecat_txt_record" {
			_, err := objMgr.GetTXTRecord(configuration, view, rs.Primary.ID)
			if err == nil {
				msg := fmt.Sprintf("TXT record %s is not removed", rs.Primary.ID)
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

func testAccTXTRecordExists(t *testing.T, resource string, name string, ttl string, text string) resource.TestCheckFunc {
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
		txtRecord, err := objMgr.GetTXTRecord(configuration, view, name)
		if err != nil {
			msg := fmt.Sprintf("Getting TXT record %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		ttlProperty := getPropertyValue("ttl", txtRecord.Properties)
		textProperty := getPropertyValue("txt", txtRecord.Properties)
		if ttlProperty != ttl || textProperty != text {
			msg := fmt.Sprintf("Getting TXT record %s failed: %s. Expect ttl=%s text=%s in properties, but received '%s'", rs.Primary.ID, err, ttl, text, txtRecord.Properties)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}
}

var txtResource1 = "txt_record_a4"
var txtName1 = "a4.example.com"
var txtText1 = "test"
var txtTTL1 = "400"
var txtProperties1 = ""
var testAccresourceTXTRecordCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_txt_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		text = "%s"
		ttl = %s
		properties = "%s"
	  }`, server, txtResource1, configuration, view, zone, txtName1, txtText1, txtTTL1, txtProperties1)

var testAccresourceTXTRecordCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_txt_record" "%s" {
		configuration = "%s"
		view = "%s"
		absolute_name = "%s"
		text = "%s"
		ttl = %s
		properties = "%s"
		}`, server, txtResource1, configuration, view, txtName1, txtText1, txtTTL1, txtProperties1)

var txtText2 = "test2"
var txtTTL2 = "4000"
var txtProperties2 = ""
var testAccresourceTXTRecordUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_txt_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		text = "%s"
		ttl = %s
		properties = "%s"
		}`, server, txtResource1, configuration, view, zone, txtName1, txtText2, txtTTL2, txtProperties2)
