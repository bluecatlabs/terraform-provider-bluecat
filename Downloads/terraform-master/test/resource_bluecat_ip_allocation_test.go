package main

import (
	"testing"
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP1, "", ipAllocateName1, ipAllocateMac1),
				),
			},
			resource.TestStep{
				Config: testAccresourceIPAllocationUpdateNotZoneMacProperties,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP1, "", ipAllocateName1, ipAllocateMac2),
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
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP1, "", ipAllocateName1, ipAllocateMac1),
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
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP1, zone, ipAllocateName1, ipAllocateMac1),
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
var ipAllocateIP1 = "1.1.0.2"
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
		}`, server, ipAllocateResource1, configuration, view, ipAllocateName1, ipAllocateNet1, ipAllocateIP1, ipAllocateMac1, ipAllocateProperties1)

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

var ipAllocateIP2 = "1.1.0.3"
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
		}`, server, ipAllocateResource1, configuration, view, ipAllocateName1, ipAllocateNet1, ipAllocateIP1, ipAllocateMac2, ipAllocateProperties2)