package main

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceIPAllocation(t *testing.T) {
	// create with full fields and update IP, Mac, properties
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceIPAllocationCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP1, zone, ipAllocateName1, ipAllocateMac1),
				),
			},
			resource.TestStep{
				Config: testAccresourceIPAllocationUpdateIPMacProperties,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP2, zone, ipAllocateName1, ipAllocateMac2),
				),
			},
		},
	})
	// create without zone and update Mac, properties
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceIPAllocationCreateNotZone,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP5, "", ipAllocateName1, ipAllocateMac5),
				),
			},
			resource.TestStep{
				Config: testAccresourceIPAllocationUpdateNotZoneMacProperties,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP5, "", ipAllocateName1, ipAllocateMac6),
				),
			},
		},
	})
	// create without zone and IP
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceIPAllocationCreateNotZoneNotIP,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipCheckExists1, "", ipAllocateName1, ipAllocateMac1),
				),
			},
		},
	})
	// create without IP
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceIPAllocationCreateNotIP,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipCheckExists1, zone, ipAllocateName1, ipAllocateMac1),
				),
			},
		},
	})
	// create with full fields include template and action field
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceIPAllocationCreateWithActionTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP3, zone, ipAllocateName3, ipAllocateMac3),
				),
			},
		},
	})
	// create without IP and assign template
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceIPAllocationCreateNotIPWithTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipCheckExists2, zone, ipAllocateName4, ipAllocateMac4),
				),
			},
		},
	})
}

func testAccCheckIPAllocationDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bluecat_ip_allocation" {
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

func testAccIPAllocationExists(t *testing.T, resource string, ip string, zone string, name string, mac string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check IP allocation on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		ipAddress, err := objMgr.GetIPAddress(configuration, ip)
		if err != nil {
			msg := fmt.Sprintf("Getting ip %s failed: %s", ip, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		macProperty := reFormatMac(getPropertyValue("macAddress", ipAddress.Properties))
		mac = reFormatMac(mac)
		if macProperty != mac {
			msg := fmt.Sprintf("Getting IP %s failed: %s. Expect %s in macAddress property, but received %s.", ip, err, mac, macProperty)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		hostRecord, err := objMgr.GetHostRecord(configuration, view, name)
		if zone == "" {
			if err == nil {
				msg := fmt.Sprintf("There is an unexpected Host record %s: %s", name, err)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
		} else {
			if err != nil {
				msg := fmt.Sprintf("Getting Host record %s failed: %s", rs.Primary.ID, err)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
			ipProperty := getPropertyValue("addresses", hostRecord.Properties)
			if ipProperty != ip {
				msg := fmt.Sprintf("Getting Host record %s failed: %s. Expect addresses=%s in properties, but received %s.", rs.Primary.ID, err, ip, ipProperty)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
		}
		return nil
	}
}

func reFormatMac(mac string) string {
	mac = strings.ReplaceAll(mac, "-", "")
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ToLower(mac)
	return mac
}

var ipAllocateResource1 = "address_allocate"

// var ipAllocateName1 = ""
var ipAllocateName1 = "allocation.example.com"
var ipAllocateNet1 = "1.1.0.0/16"
var ipAllocateIP1 = "1.1.0.15"
var ipAllocateMac1 = "223344556677"
var ipAllocateProperties1 = ""
var testAccresourceIPAllocationCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "%s"
		network = "%s"
		ip4_address = "%s"
		mac_address = "%s"
		properties = "%s"
	  }`, server, ipAllocateResource1, configuration, view, zone, ipAllocateName1, ipAllocateNet1, ipAllocateIP1, ipAllocateMac1, ipAllocateProperties1)

var ipAllocateIP5 = "1.1.0.12"
var ipAllocateMac5 = "223344556699"

var testAccresourceIPAllocationCreateNotZone = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		name = "%s"
		network = "%s"
		ip4_address = "%s"
		mac_address = "%s"
		properties = "%s"
		}`, server, ipAllocateResource1, configuration, view, ipAllocateName1, ipAllocateNet1, ipAllocateIP5, ipAllocateMac5, ipAllocateProperties1)

var ipCheckExists1 = "1.1.0.4"
var testAccresourceIPAllocationCreateNotIP = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "%s"
		network = "%s"
		mac_address = "%s"
		properties = "%s"
		}`, server, ipAllocateResource1, configuration, view, zone, ipAllocateName1, ipAllocateNet1, ipAllocateMac1, ipAllocateProperties1)

var testAccresourceIPAllocationCreateNotZoneNotIP = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		name = "%s"
		network = "%s"
		mac_address = "%s"
		properties = "%s"
		}`, server, ipAllocateResource1, configuration, view, ipAllocateName1, ipAllocateNet1, ipAllocateMac1, ipAllocateProperties1)

var ipAllocateIP2 = "1.1.0.17"
var ipAllocateMac2 = "888888888888"
var ipAllocateProperties2 = ""
var testAccresourceIPAllocationUpdateIPMacProperties = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "%s"
		network = "%s"
		ip4_address = "%s"
		mac_address = "%s"
		properties = "%s"
	}`, server, ipAllocateResource1, configuration, view, zone, ipAllocateName1, ipAllocateNet1, ipAllocateIP2, ipAllocateMac2, ipAllocateProperties2)

var ipAllocateMac6 = "887788888888"

var testAccresourceIPAllocationUpdateNotZoneMacProperties = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		name = "%s"
		network = "%s"
		ip4_address = "%s"
		mac_address = "%s"
		properties = "%s"
		}`, server, ipAllocateResource1, configuration, view, ipAllocateName1, ipAllocateNet1, ipAllocateIP5, ipAllocateMac6, ipAllocateProperties2)

var actionDhcpReserved = "MAKE_DHCP_RESERVED"
var template = "template1"
var ipAllocateIP3 = "1.1.0.8"
var ipAllocateMac3 = "888888888877"
var ipAllocateName3 = "allocation3.example.com"
var testAccresourceIPAllocationCreateWithActionTemplate = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "%s"
		network = "%s"
		ip4_address = "%s"
		mac_address = "%s"
		properties = "%s"
		action = "%s"
		template = "%s"
		}`, server, ipAllocateResource1, configuration, view, zone, ipAllocateName3, ipAllocateNet1, ipAllocateIP3, ipAllocateMac3, ipAllocateProperties1, actionDhcpReserved, template)

var ipCheckExists2 = "1.1.0.5"
var ipAllocateMac4 = "778888888877"
var ipAllocateName4 = "allocation4.example.com"

var testAccresourceIPAllocationCreateNotIPWithTemplate = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "%s"
		network = "%s"
		mac_address = "%s"
		properties = "%s"
		action = "%s"
		template = "%s"
		}`, server, ipAllocateResource1, configuration, view, zone, ipAllocateName4, ipAllocateNet1, ipAllocateMac4, ipAllocateProperties1, actionDhcpReserved, template)
