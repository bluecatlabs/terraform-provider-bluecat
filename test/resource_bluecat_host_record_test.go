package main

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceHostRecord(t *testing.T) {
	// create with full fields and update
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckHostRecordDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccresourceHostRecordCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccHostRecordExists(t, fmt.Sprintf("bluecat_host_record.%s", hostResource1), hostName1, hostTTL1, hostIP1, hostReverseProperty1),
				),
			},
			// update
			{
				Config: testAccresourceHostRecordUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccHostRecordExists(t, fmt.Sprintf("bluecat_host_record.%s", hostResource1), hostName1, hostTTL2, hostIP2, hostReverseProperty2),
				),
			},
		},
	})
	// create without some optional fields
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckHostRecordDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceHostRecordCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccHostRecordExists(t, fmt.Sprintf("bluecat_host_record.%s", hostResource1), hostName1, hostTTL1, hostIP1, hostReverseProperty1),
				),
			},
		},
	})
}

func testAccCheckHostRecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bluecat_host_record" {
			msg := fmt.Sprintf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
			log.Error(msg)
			// return fmt.Errorf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
		}
		_, err := objMgr.GetHostRecord(configuration, view, rs.Primary.ID)
		if err == nil {
			msg := fmt.Sprintf("Host record %s is not removed", rs.Primary.ID)
			log.Error(msg)
			return fmt.Errorf("Host record %s is not removed", rs.Primary.ID)
		}
	}
	return nil
}

func testAccHostRecordExists(t *testing.T, resource string, name string, ttl string, ip string, hostReverseProperty string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check Host record on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		hostRecord, err := objMgr.GetHostRecord(configuration, view, name)
		if err != nil {
			msg := fmt.Sprintf("Getting Host record %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf("Getting Host record %s failed: %s", rs.Primary.ID, err)
		}
		if checkValidHostRecord(*hostRecord, ttl, ip, hostReverseProperty) == false {
			msg := fmt.Sprintf("Getting Host record %s failed: %s. Expect ttl=%s addresses=%s reverseRecord=%s in properties, but received '%s'", rs.Primary.ID, err, ttl, ip, hostReverseProperty, hostRecord.Properties)
			log.Error(msg)
			return fmt.Errorf("Getting Host record %s failed: %s. Expect ttl=%s addresses=%s reverseRecord=%s in properties, but received '%s'", rs.Primary.ID, err, ttl, ip, hostReverseProperty, hostRecord.Properties)
		}
		return nil
	}
}

func checkValidHostRecord(hostRecord entities.HostRecord, ttl string, ip string, hostReverseProperty string) bool {
	ttlProperty := utils.GetPropertyValue("ttl", hostRecord.Properties)
	ipProperty := utils.GetPropertyValue("addresses", hostRecord.Properties)
	reverseProperty := utils.GetPropertyValue("reverseRecord", hostRecord.Properties)
	if ttlProperty != ttl || ipProperty != ip || reverseProperty != hostReverseProperty {
		return false
	}
	return true
}

var hostResource1 = "host_record_a2"
var hostName1 = "a2.example.com"
var hostIP1 = "1.1.0.2"
var hostTTL1 = "200"
var hostProperties1 = "reverseRecord=true|"
var hostReverseProperty1 = "true"
var testAccresourceHostRecordCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		ip_address = "%s"
		ttl = %s
		properties = "%s"
		depends_on = [bluecat_zone.sub_zone_test, bluecat_ipv4network.network_test]
	  }`, GetTestEnvResources(), hostResource1, configuration, view, zone, hostName1, hostIP1, hostTTL1, hostProperties1)

var testAccresourceHostRecordCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		absolute_name = "%s"
		ip_address = "%s"
		ttl = %s
		properties = "%s"
		depends_on = [bluecat_zone.sub_zone_test, bluecat_ipv4network.network_test]
		}`, GetTestEnvResources(), hostResource1, configuration, view, hostName1, hostIP1, hostTTL1, hostProperties1)

var hostIP2 = "1.1.0.3"
var hostTTL2 = "300"
var hostProperties2 = "reverseRecord=false|"
var hostReverseProperty2 = "false"
var testAccresourceHostRecordUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_host_record" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		absolute_name = "%s"
		ip_address = "%s"
		ttl = %s
		properties = "%s"
		depends_on = [bluecat_zone.sub_zone_test, bluecat_ipv4network.network_test]
		}`, GetTestEnvResources(), hostResource1, configuration, view, zone, hostName1, hostIP2, hostTTL2, hostProperties2)
