package main

import (
	"fmt"
	"net"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceIPAllocation(t *testing.T) {
	// create with full fields and update IP, Mac, properties
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			// allocate IPv4 address
			{
				Config: testAccresourceIPAllocationCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP1, zone, ipAllocateName1, ipAllocateMac1, entities.IPV4),
				),
			},
			//// update IPv4 address
			//// change MAC address from 223344556699 to 887788888888
			//// change state from MAKE_STATIC to MAKE_DHCP_RESERVED
			{
				Config: testAccresourceIPAllocationUpdateIPMacProperties,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP1, zone, ipAllocateName1, ipAllocateMac2, entities.IPV4),
				),
			},

			// allocate IPv6 address
			{
				Config: testAccresourceIPv6AllocationCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(
						t,
						"bluecat_ip_allocation.ipv6_address_allocation",
						"2003:1000::15", "example", "ip6allocation.example.com", "A1:B1:C1:D1:E1:F1", entities.IPV6),
				),
			},
			// update IPv6 address (changing MAC address from A1:B1:C1:D1:E1:F1 to A2:B2:C2:D2:E2:F2)
			{
				Config: testAccResourceIPv6AllocationUpdateIPMacProperties,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(
						t,
						"bluecat_ip_allocation.ipv6_address_allocation",
						"2003:1000::15", "example", "ip6allocation.example.com", "A2:B2:C2:D2:E2:F2", entities.IPV6),
				),
			},
		},
	})

	// create without zone and update Mac, properties
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			// allocate IPv4 address 1.1.0.12 without Zone
			{
				Config: testAccresourceIPAllocationCreateNotZone,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(
						t,
						fmt.Sprintf("bluecat_ip_allocation.%s",
							ipAllocateResource1), ipAllocateIP5, "", ipAllocateName1, ipAllocateMac5, entities.IPV4,
					),
				),
			},
			// change MAC from 223344556699 to 887788888888 for the IPv4
			{
				Config: testAccresourceIPAllocationUpdateNotZoneMacProperties,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(
						t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1),
						ipAllocateIP5, "", ipAllocateName1, ipAllocateMac6, entities.IPV4,
					),
				),
			},

			// allocate IPv6 address 2003:1000::100 without Zone
			{
				Config: testAccResourceIPv6AllocationCreateNotZone,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(
						t, "bluecat_ip_allocation.ipv6_address_allocation_wo_zone",
						"2003:1000::100", "", "ip6allocation_wo_zone.example.com", "11:11:11:11:11:11",
						entities.IPV6,
					),
				),
			},
			// change MAC from 11:11:11:11:11:11 to 22:22:22:22:22:22 for the IPv6
			{
				Config: testAccresourceIPv6AllocationUpdateNotZoneMacProperties,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(
						t, "bluecat_ip_allocation.ipv6_address_allocation_wo_zone",
						"2003:1000::100", "", "ip6allocation_wo_zone.example.com", "22:22:22:22:22:22",
						entities.IPV6,
					),
				),
			},
		},
	})

	// allocate IP address - no Zone and IP passed
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			// allocate IPv4 address without passing Zone and IP address
			{
				Config: testAccresourceIPAllocationCreateNotZoneNotIP,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipCheckExists1, "", ipAllocateName1, ipAllocateMac1, entities.IPV4),
				),
			},
			// allocate IPv6 address without passing Zone and IP address
			{
				Config: testAccresourceIPv6AllocationCreateNotZoneNotIP,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(
						t, "bluecat_ip_allocation.ipv6_address_allocation_wo_zone_and_ip",
						"2003:1000::1", "", "ip6allocation_wo_zone_and_ip.example.com", "AA:AA:AA:11:11:11",
						entities.IPV6,
					),
				),
			},
		},
	})

	// create without IP
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			// allocate IPv4 address without passing IP address
			{
				Config: testAccresourceIPAllocationCreateNotIP,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipCheckExists1, zone, ipAllocateName1, ipAllocateMac1, entities.IPV4),
				),
			},
			// allocate IPv6 address without passing IP address
			{
				Config: testAccResourceIPv6AllocationCreateNotIP,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(
						t, "bluecat_ip_allocation.ipv6_address_allocation_wo_ip",
						"2003:1000::1", zone, "ip6allocation_wo_ip.example.com", "AA:BB:CC:11:22:33",
						entities.IPV6,
					),
				),
			},
		},
	})

	// create with full fields include template and action field
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccresourceIPAllocationCreateWithActionTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipAllocateIP3, zone, ipAllocateName3, ipAllocateMac3, entities.IPV4),
				),
			},
		},
	})

	// create without IP and assign template
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPAllocationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccresourceIPAllocationCreateNotIPWithTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccIPAllocationExists(t, fmt.Sprintf("bluecat_ip_allocation.%s", ipAllocateResource1), ipCheckExists1, zone, ipAllocateName4, ipAllocateMac4, entities.IPV4),
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
			// return fmt.Errorf(msg)
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

func testAccIPAllocationExists(t *testing.T, resource string, ip string, zone string, name string, mac string, ipVersion string) resource.TestCheckFunc {
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
		ipAddress, err := objMgr.GetIPAddress(configuration, ip, ipVersion)
		if err != nil {
			msg := fmt.Sprintf("Getting ip %s failed: %s", ip, err)
			log.Error(msg)
			return fmt.Errorf("Getting ip %s failed: %s", ip, err)
		}
		macProperty := reFormatMac(utils.GetPropertyValue("macAddress", ipAddress.Properties))
		mac = reFormatMac(mac)
		if macProperty != mac {
			msg := fmt.Sprintf("Getting IP %s failed: %s. Expect %s in macAddress property, but received %s.", ip, err, mac, macProperty)
			log.Error(msg)
			return fmt.Errorf("Getting IP %s failed: %s. Expect %s in macAddress property, but received %s.", ip, err, mac, macProperty)
		}

		hostRecord, err := objMgr.GetHostRecord(configuration, view, name)
		if err != nil {
			msg := fmt.Sprintf("Getting Host record %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf("Getting Host record %s failed: %s", rs.Primary.ID, err)
		}
		ipProperty := utils.GetPropertyValue("addresses", hostRecord.Properties)
		ipAddressLong := net.ParseIP(ipProperty)
		fmt.Sprintf("%s", ipAddressLong)
		if ipAddressLong.String() != ip {
			msg := fmt.Sprintf("Getting Host record %s failed: %s. Expect addresses=%s in properties, but received %s.", rs.Primary.ID, err, ip, ipProperty)
			log.Error(msg)
			return fmt.Errorf("Getting Host record %s failed: %s. Expect addresses=%s in properties, but received %s.", rs.Primary.ID, err, ip, ipProperty)
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
		ip_address = "%s"
		mac_address = "%s"
		properties = "%s"
		action = "MAKE_STATIC"
		depends_on = [bluecat_ipv4network.network_test, bluecat_zone.sub_zone_test]
	  }`, GetTestEnvResources(), ipAllocateResource1, configuration, view, zone, ipAllocateName1, ipAllocateNet1, ipAllocateIP1, ipAllocateMac1, ipAllocateProperties1)

var ipAllocateIP5 = "1.1.0.12"
var ipAllocateMac5 = "223344556699"

var testAccresourceIPAllocationCreateNotZone = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		name = "%s"
		network = "%s"
		ip_address = "%s"
		mac_address = "%s"
		properties = "%s"
		depends_on = [bluecat_ipv4network.network_test, bluecat_zone.sub_zone_test]
		}`, GetTestEnvResources(), ipAllocateResource1, configuration, view, ipAllocateName1, ipAllocateNet1, ipAllocateIP5, ipAllocateMac5, ipAllocateProperties1)

var ipCheckExists1 = "1.1.0.2"
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
		depends_on = [bluecat_ipv4network.network_test, bluecat_zone.sub_zone_test]
		}`, GetTestEnvResources(), ipAllocateResource1, configuration, view, zone, ipAllocateName1, ipAllocateNet1, ipAllocateMac1, ipAllocateProperties1)

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
		ip_address = "%s"
		mac_address = "%s"
		properties = "%s"
		action = "MAKE_DHCP_RESERVED"
		depends_on = [bluecat_ipv4network.network_test, bluecat_zone.sub_zone_test]
	}`, GetTestEnvResources(), ipAllocateResource1, configuration, view, zone, ipAllocateName1, ipAllocateNet1, ipAllocateIP1, ipAllocateMac2, ipAllocateProperties2)

var ipAllocateMac6 = "887788888888"

var testAccresourceIPAllocationUpdateNotZoneMacProperties = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		name = "%s"
		network = "%s"
		ip_address = "%s"
		mac_address = "%s"
		properties = "%s"
		depends_on = [bluecat_ipv4network.network_test, bluecat_zone.sub_zone_test]
		}`, GetTestEnvResources(), ipAllocateResource1, configuration, view, ipAllocateName1, ipAllocateNet1, ipAllocateIP5, ipAllocateMac6, ipAllocateProperties2)

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
		ip_address = "%s"
		mac_address = "%s"
		properties = "%s"
		action = "%s"
		template = "%s"
		depends_on = [bluecat_ipv4network.network_test, bluecat_zone.sub_zone_test]
		}`, GetTestEnvResources(), ipAllocateResource1, configuration, view, zone, ipAllocateName3, ipAllocateNet1, ipAllocateIP3, ipAllocateMac3, ipAllocateProperties1, actionDhcpReserved, template)

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
		depends_on = [bluecat_ipv4network.network_test, bluecat_zone.sub_zone_test]
		}`, GetTestEnvResources(), ipAllocateResource1, configuration, view, zone, ipAllocateName4, ipAllocateNet1, ipAllocateMac4, ipAllocateProperties1, actionDhcpReserved, template)

var testAccresourceIPAllocationCreateNotZoneNotIP = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "%s" {
		configuration = "%s"
		view = "%s"
		name = "%s"
		network = "%s"
		mac_address = "%s"
		properties = "%s"
		depends_on = [bluecat_ipv4network.network_test, bluecat_zone.sub_zone_test]
		}`, GetTestEnvResources(), ipAllocateResource1, configuration, view, ipAllocateName1, ipAllocateNet1, ipAllocateMac1, ipAllocateProperties1)

var testAccresourceIPv6AllocationCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "ipv6_address_allocation" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "ip6allocation.example.com"
		network = "2003:1000::/64"
		ip_address = "2003:1000::15"
		mac_address = "A1:B1:C1:D1:E1:F1"
		properties = ""
		action = "MAKE_STATIC"
		ip_version = "ipv6"
		depends_on = [bluecat_ipv6network.ipv6_network_test, bluecat_zone.sub_zone_test]
	  }`, GetTestEnvResources(), configuration, view, zone)

var testAccResourceIPv6AllocationUpdateIPMacProperties = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "ipv6_address_allocation" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "ip6allocation.example.com"
		network = "2003:1000::/64"
		ip_address = "2003:1000::15"
		mac_address = "A2:B2:C2:D2:E2:F2"
		properties = ""
		ip_version = "ipv6"
		depends_on = [bluecat_ipv6network.ipv6_network_test, bluecat_zone.sub_zone_test]
	}`, GetTestEnvResources(), configuration, view, zone)

var testAccResourceIPv6AllocationCreateNotZone = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "ipv6_address_allocation_wo_zone" {
		configuration = "%s"
		view = "%s"
		name = "ip6allocation_wo_zone.example.com"
		network = "2003:1000::"
		ip_address = "2003:1000::100"
		mac_address = "11:11:11:11:11:11"
		properties = ""
		ip_version = "ipv6"
		depends_on = [bluecat_ipv6network.ipv6_network_test, bluecat_zone.sub_zone_test]
		}`, GetTestEnvResources(), configuration, view)

var testAccresourceIPv6AllocationUpdateNotZoneMacProperties = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "ipv6_address_allocation_wo_zone" {
		configuration = "%s"
		view = "%s"
		name = "%s"
		network = "2003:1000::"
		ip_address = "2003:1000::100"
		mac_address = "22:22:22:22:22:22"
		properties = ""
		ip_version = "ipv6"
		depends_on = [bluecat_ipv6network.ipv6_network_test, bluecat_zone.sub_zone_test]
		}`, GetTestEnvResources(), configuration, view, ipAllocateName1)

var testAccresourceIPv6AllocationCreateNotZoneNotIP = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "ipv6_address_allocation_wo_zone_and_ip" {
		configuration = "%s"
		view = "%s"
		name = "ip6allocation_wo_zone_and_ip.example.com"
		network = "2003:1000::"
		mac_address = "AA:AA:AA:11:11:11"
		properties = ""
		ip_version = "ipv6"
		depends_on = [bluecat_ipv6network.ipv6_network_test, bluecat_zone.sub_zone_test]
		}`, GetTestEnvResources(), configuration, view)

var testAccResourceIPv6AllocationCreateNotIP = fmt.Sprintf(
	`%s
	resource "bluecat_ip_allocation" "ipv6_address_allocation_wo_ip" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		name = "ip6allocation_wo_ip.example.com"
		network = "2003:1000::"
		mac_address = "AA:BB:CC:11:22:33"
		properties = ""
		ip_version = "ipv6"
		depends_on = [bluecat_ipv6network.ipv6_network_test, bluecat_zone.sub_zone_test]
		}`, GetTestEnvResources(), configuration, view, zone)
