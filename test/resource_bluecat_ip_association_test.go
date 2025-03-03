package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"
)

func TestAccResourceIPAssociation(t *testing.T) {
	// create with full fields and update
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccresourceHostRecordCreate2,
			},
			// create
			{
				Config: testAccresourceIPAssocitionCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAssociationExists(t, fmt.Sprintf("bluecat_ip_association.%s", IPAssociateResource1), IPAssociateName1, IPAssociateIP1, IPAssociateMac1),
				),
			},
			// update
			{
				Config: testAccresourceIPAssocitionUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAssociationExists(t, fmt.Sprintf("bluecat_ip_association.%s", IPAssociateResource1), IPAssociateName1, IPAssociateIP1, IPAssociateMac2),
				),
			},
		},
	})
	// create without some optional fields and update
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPAssociationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceHostRecordCreate2,
			},
			// create
			resource.TestStep{
				Config: testAccresourceIPAssocitionCreateNotZone,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAssociationExists(t, fmt.Sprintf("bluecat_ip_association.%s", IPAssociateResource1), IPAssociateName1, IPAssociateIP1, IPAssociateMac1),
				),
			},
			// update
			resource.TestStep{
				Config: testAccresourceIPAssocitionUpdateNotZone,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAssociationExists(t, fmt.Sprintf("bluecat_ip_association.%s", IPAssociateResource1), IPAssociateName1, IPAssociateIP1, IPAssociateMac2),
				),
			},
		},
	})
}

func testAccCheckIPAssociationDestroy(s *terraform.State) error {
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
		} else if rs.Type != "bluecat_ip_association" {
			msg := fmt.Sprintf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

	}
	return nil
}

func testAccIPAssociationExists(t *testing.T, resource string, name string, ip string, mac string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check IP Association on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		ipAddress, err := objMgr.GetIPAddress(configuration, ip, "ipv4")
		if err != nil {
			msg := fmt.Sprintf("Getting IP %s failed: %s", ip, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		macProperty := reFormatMac(getPropertyValue("macAddress", ipAddress.Properties))
		mac = reFormatMac(mac)
		if macProperty != mac {
			msg := fmt.Sprintf("Getting IP %s failed: %s. Expect macAddress=%s in properties, but received %s.", ip, err, mac, macProperty)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}
}

var testAccresourceHostRecordCreate2 = fmt.Sprintf(
	`%s
	resource "bluecat_host_record" "host_record_a2" {
		configuration = "%s"
		view = "%s"
		absolute_name = "a1.example.com"
		ip_address = "1.1.0.10"
		ttl = 200
		properties = ""
		}`, server, configuration, view)

var IPAssociateResource1 = "address_allocate"
var IPAssociateName1 = "a1.example.com"
var IPAssociateNet1 = "1.1.0.0/16"
var IPAssociateIP1 = "1.1.0.10"
var IPAssociateMac1 = "AB3344556677"
var IPAssociateProperties1 = ""
var IPAssociateDescriptionProperty1 = "terraform testing config"
var testAccresourceIPAssocitionCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ip_association" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "%s"
		network = "%s"
		ip_address = "%s"
		mac_address = "%s"
		properties = "%s"
	  }`, testAccresourceHostRecordCreate2, IPAssociateResource1, configuration, view, zone, IPAssociateName1, IPAssociateNet1, IPAssociateIP1, IPAssociateMac1, IPAssociateProperties1)

var testAccresourceIPAssocitionCreateNotZone = fmt.Sprintf(
	`%s
	resource "bluecat_ip_association" "%s" {
		configuration = "%s"
		view = "%s"
		name = "%s"
		network = "%s"
		ip_address = "%s"
		mac_address = "%s"
		properties = "%s"
		}`, testAccresourceHostRecordCreate2, IPAssociateResource1, configuration, view, IPAssociateName1, IPAssociateNet1, IPAssociateIP1, IPAssociateMac1, IPAssociateProperties1)

var IPAssociateMac2 = "ABABABABABAB"
var testAccresourceIPAssocitionUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ip_association" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "%s"
		network = "%s"
		ip_address = "%s"
		mac_address = "%s"
		properties = "%s"
	  }`, testAccresourceHostRecordCreate2, IPAssociateResource1, configuration, view, zone, IPAssociateName1, IPAssociateNet1, IPAssociateIP1, IPAssociateMac2, IPAssociateProperties1)

var testAccresourceIPAssocitionUpdateNotZone = fmt.Sprintf(
	`%s
	resource "bluecat_ip_association" "%s" {
		configuration = "%s"
		view = "%s"
		name = "%s"
		network = "%s"
		ip_address = "%s"
		mac_address = "%s"
		properties = "%s"
		}`, testAccresourceHostRecordCreate2, IPAssociateResource1, configuration, view, IPAssociateName1, IPAssociateNet1, IPAssociateIP1, IPAssociateMac2, IPAssociateProperties1)
