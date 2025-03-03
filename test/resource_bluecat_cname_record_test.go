package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"
)

func TestAccResourceCNAMERecord(t *testing.T) {
	// create with full fields and update
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckCNAMERecordDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccresourceCNAMERecordCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccCNAMERecordExists(t, fmt.Sprintf("bluecat_cname_record.%s", cnameResource1), cnameName1, cnameTTL1, cnameLink1),
				),
			},
			// // update
			{
				Config: testAccresourceCNAMERecordUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccCNAMERecordExists(t, fmt.Sprintf("bluecat_cname_record.%s", cnameResource1), cnameName1, cnameTTL2, cnameLink2),
				),
			},
		},
	})
	// create without some optional fields
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckCNAMERecordDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccresourceCNAMERecordCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccCNAMERecordExists(t, fmt.Sprintf("bluecat_cname_record.%s", cnameResource1), cnameName1, cnameTTL1, cnameLink1),
				),
			},
		},
	})
}

func testAccCheckCNAMERecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "bluecat_host_record" {
			_, err := objMgr.GetHostRecord(configuration, view, rs.Primary.ID)
			if err == nil {
				msg := fmt.Sprintf("Host record %s is not removed", rs.Primary.ID)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
		} else if rs.Type == "bluecat_cname_record" {
			_, err := objMgr.GetCNAMERecord(configuration, view, rs.Primary.ID)
			if err == nil {
				msg := fmt.Sprintf("CNAME record %s is not removed", rs.Primary.ID)
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

func testAccCNAMERecordExists(t *testing.T, resource string, name string, ttl string, link string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check CNAME record on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		cnameRecord, err := objMgr.GetCNAMERecord(configuration, view, name)
		if err != nil {
			msg := fmt.Sprintf("Getting CNAME record %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		ttlProperty := getPropertyValue("ttl", cnameRecord.Properties)
		linkProperty := getPropertyValue("linkedRecordName", cnameRecord.Properties)
		if ttlProperty != ttl || linkProperty != link {
			msg := fmt.Sprintf("Getting CNAME record %s failed: %s. Expect ttl=%s linkedRecordName=%s in properties, but received '%s'", rs.Primary.ID, err, ttl, link, cnameRecord.Properties)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}
}

var hostCnameResource1 = "host_record_a2"
var hostCnameResource2 = "host_record_a3"
var testAccresourceHostRecordCreate1 = fmt.Sprintf(
	`%s
	resource "bluecat_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		absolute_name = "a2.example.com"
		ip_address = "1.1.0.2"
		ttl = 200
		properties = ""
		}

	resource "bluecat_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		absolute_name = "a3.example.com"
		ip_address = "1.1.0.3"
		ttl = 300
		properties = ""
		}`, server, hostCnameResource1, configuration, view, hostCnameResource2, configuration, view)

var cnameResource1 = "cname_record_a4"
var cnameName1 = "a4.example.com"
var cnameLink1 = "a2.example.com"
var cnameTTL1 = "400"
var cnameProperties1 = ""
var testAccresourceCNAMERecordCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_cname_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		linked_record = "%s"
		ttl = %s
		properties = "%s"
		depends_on = [bluecat_host_record.%s, bluecat_host_record.%s]
	  }`, testAccresourceHostRecordCreate1, cnameResource1, configuration, view, zone, cnameName1, cnameLink1, cnameTTL1, cnameProperties1, hostCnameResource1, hostCnameResource2)

var testAccresourceCNAMERecordCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_cname_record" "%s" {
		configuration = "%s"
		view = "%s"
		absolute_name = "%s"
		linked_record = "%s"
		ttl = %s
		properties = "%s"
		depends_on = [bluecat_host_record.%s, bluecat_host_record.%s]
		}`, testAccresourceHostRecordCreate1, cnameResource1, configuration, view, cnameName1, cnameLink1, cnameTTL1, cnameProperties1, hostCnameResource1, hostCnameResource2)

var cnameLink2 = "a3.example.com"
var cnameTTL2 = "4000"
var cnameProperties2 = ""
var testAccresourceCNAMERecordUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_cname_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		linked_record = "%s"
		ttl = %s
		properties = "%s"
		depends_on = [bluecat_host_record.%s, bluecat_host_record.%s]
		}`, testAccresourceHostRecordCreate1, cnameResource1, configuration, view, zone, cnameName1, cnameLink2, cnameTTL2, cnameProperties2, hostCnameResource1, hostCnameResource2)
