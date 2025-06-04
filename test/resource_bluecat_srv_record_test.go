package main

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceSRVRecord(t *testing.T) {
	// create with full fields and update
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSRVRecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceSRVRecordCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccSRVRecordExists(t, fmt.Sprintf("bluecat_srv_record.%s", srvResource1), srvName1, srvPriority1, srvWeight1),
				),
			},
			// // update
			resource.TestStep{
				Config: testAccresourceSRVRecordUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccSRVRecordExists(t, fmt.Sprintf("bluecat_srv_record.%s", srvResource1), srvName1, srvPriority2, srvWeight2),
				),
			},
		},
	})
	// create without some optional fields
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSRVRecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceSRVRecordCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccSRVRecordExists(t, fmt.Sprintf("bluecat_srv_record.%s", srvResource1), srvName1, srvPriority1, srvWeight1),
				),
			},
		},
	})
}

func testAccCheckSRVRecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "bluecat_srv_record" {
			_, err := objMgr.GetSRVRecord(configuration, view, rs.Primary.ID)
			if err == nil {
				msg := fmt.Sprintf("SRV record %s is not removed", rs.Primary.ID)
				log.Error(msg)
				return fmt.Errorf("SRV record %s is not removed", rs.Primary.ID)
			}
		} else {
			msg := fmt.Sprintf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
			log.Error(msg)
			// return fmt.Errorf(msg)
		}
	}
	return nil
}

func testAccSRVRecordExists(t *testing.T, resource string, name string, priority string, weight string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check SRV record on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		srvRecord, err := objMgr.GetSRVRecord(configuration, view, name)
		if err != nil {
			msg := fmt.Sprintf("Getting SRV record %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf("Getting SRV record %s failed: %s", rs.Primary.ID, err)
		}
		ttlProperty := utils.GetPropertyValue("priority", srvRecord.Properties)
		weightProperty := utils.GetPropertyValue("weight", srvRecord.Properties)
		if ttlProperty != priority || weightProperty != weight {
			msg := fmt.Sprintf("Getting SRV record %s failed: %s. Expect priority=%s weight=%s in properties, but received '%s'", rs.Primary.ID, err, priority, weight, srvRecord.Properties)
			log.Error(msg)
			return fmt.Errorf("Getting SRV record %s failed: %s. Expect priority=%s weight=%s in properties, but received '%s'", rs.Primary.ID, err, priority, weight, srvRecord.Properties)
		}
		return nil
	}
}

var testAccresourceHostSRVRecordCreate = fmt.Sprintf(
	`%s
	resource "bluecat_host_record" "srv_host_record" {
		configuration = "%s"
		view = "%s"
		absolute_name = "ns1.example.com"
		ip_address = "1.1.0.2"
		ttl = 200
		properties = ""
		depends_on = [bluecat_zone.sub_zone_test, bluecat_ipv4network.network_test]
		}`, GetTestEnvResources(), configuration, view)

var srvResource1 = "srv_record_a4"
var srvName1 = "a4.example.com"
var srvWeight1 = "20"
var srvPort1 = "8080"
var srvPriority1 = "2"
var srvTTL1 = "400"
var srvProperties1 = ""
var testAccresourceSRVRecordCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_srv_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		linked_record = "ns1.example.com"
		weight = "%s"
		port = "%s"
		priority = "%s"
		ttl = %s
		properties = "%s"
		depends_on = [bluecat_zone.sub_zone_test, bluecat_host_record.srv_host_record]
	  }`, testAccresourceHostSRVRecordCreate, srvResource1, configuration, view, zone, srvName1, srvWeight1, srvPort1, srvPriority1, srvTTL1, srvProperties1)

var testAccresourceSRVRecordCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_srv_record" "%s" {
		configuration = "%s"
		view = "%s"
		absolute_name = "%s"
		linked_record = "ns1.example.com"
		weight = "%s"
		port = "%s"
		priority = "%s"
		properties = "%s"
		depends_on = [bluecat_zone.sub_zone_test, bluecat_host_record.srv_host_record]
		}`, testAccresourceHostSRVRecordCreate, srvResource1, configuration, view, srvName1, srvWeight1, srvPort1, srvPriority1, srvProperties1)

var srvWeight2 = "21"
var srvPort2 = "8081"
var srvPriority2 = "3"
var srvTTL2 = "4000"
var srvProperties2 = ""
var testAccresourceSRVRecordUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_srv_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		linked_record = "ns1.example.com"
		weight = "%s"
		port = "%s"
		priority = "%s"
		ttl = %s
		properties = "%s"
		depends_on = [bluecat_zone.sub_zone_test, bluecat_host_record.srv_host_record]
		}`, testAccresourceHostSRVRecordCreate, srvResource1, configuration, view, zone, srvName1, srvWeight2, srvPort2, srvPriority2, srvTTL2, srvProperties2)
