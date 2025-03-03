package main

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceNetwork(t *testing.T) {
	// create with full fields and update without some optional fields
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccResourceNetworkCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", netResource1), netName1, netCIDR1, netAllowDuplicateHost1, netGateway1, netReserveIPValue1, "", entities.IPV4),
				),
			},
			// update
			{
				Config: testAccResourceNetworkUpdateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", netResource1), "", netCIDR1, netAllowDuplicateHost2, netGateway2, netReserveIPValue1, "", entities.IPV4),
				),
			},

			// create ipv6 network inside created block Full Field
			{
				Config: testAccResourceIPv6NetworkCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(
						t,
						"bluecat_ipv6network.ipv6_net_record_1",
						"ipv6_net_name", "2040:B041::/64", "", "", "", "", entities.IPV6),
				),
			},
			// create ipv6 network inside Unique Local Address Space block
			{
				Config: testAccResourceIPv6NetworkInsideUniqueLocal,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(
						t,
						"bluecat_ipv6network.ipv6_net_record_2",
						"ipv6_net_name", "FC00::/64", "", "", "", "", entities.IPV6),
				),
			},
		},
	})

	// create without some optional fields and update with full fields
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			// create ipv4 network
			{
				Config: testAccResourceNetworkCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", netResource1), "", netCIDR1, netAllowDuplicateHost1, netGateway1, netReserveIPValue1, "", "ipv4"),
				),
			},
			// update ipv4 network
			{
				Config: testAccResourceNetworkUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", netResource1), netName2, netCIDR1, netAllowDuplicateHost2, netGateway2, netReserveIPValue1, "", entities.IPV4),
				),
			},

			// create ipv6 network inside created block Not Full Field
			{
				Config: testAccResourceIPv6NetworkCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, "bluecat_ipv6network.ipv6_net_record_1_not_full_field", "",
						"FC00::/64", "", "", "", "", entities.IPV6),
				),
			},
			// update ipv6 network inside created block Not Full Field
			{
				Config: testAccResourceIPv6NetworUpdateAddName,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, "bluecat_ipv6network.ipv6_net_record_1_not_full_field", "new_name",
						"FC00::/64", "", "", "", "", entities.IPV6),
				),
			},
		},
	})

	// create ipv4 network with full fields include template field
	// template is not supported for the ipv6 network
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccResourceNetworkCreateFullFieldWithTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", netResource1), netName1, netCIDR1, netAllowDuplicateHost1, netGateway3, netReserveIPValue2, templateName, entities.IPV4),
				),
			},
		},
	})

	// create next available network with full fields and update without some optional fields
	// next available network is not supported for the ipv6 network
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccResourceNextNetworkCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", nextNetResource1), nextNetName1, nextNetCIDRValue1, nextNetAllowDuplicateHost1, nextNetGatewayValue1, nextNetReserveIPValue1, "", entities.IPV4),
				),
			},
			// update
			{
				Config: testAccResourceNextNetworkUpdateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", nextNetResource1), nextNetName2, nextNetCIDRValue1, nextNetAllowDuplicateHost1, nextNetGatewayValue1, nextNetReserveIPValue1, "", entities.IPV4),
				),
			},
		},
	})
}

func testAccCheckNetworkDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	network := entities.Network{}
	network.Configuration = configuration
	network.IPVersion = entities.IPV4

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "bluecat_ipv6network" || rs.Type == "bluecat_ipv6block" {
			network.IPVersion = entities.IPV6
		}
		network.CIDR = rs.Primary.ID
		if rs.Type == "bluecat_ipv4network" || rs.Type == "bluecat_ipv6network" {
			_, err := objMgr.GetNetwork(&network)
			if err == nil {
				msg := fmt.Sprintf("Network %s is not removed", rs.Primary.ID)
				log.Error(msg)
				return fmt.Errorf(msg)
			}

		} else if rs.Type == "bluecat_ipv4block" || rs.Type == "bluecat_ipv6block" {
			cidr := strings.Split(rs.Primary.ID, "/")
			_, err := objMgr.GetBlock(configuration, cidr[0], cidr[1], network.IPVersion)
			if err == nil {
				msg := fmt.Sprintf("Block %s is not removed", rs.Primary.ID)
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

func testAccNetworkExists(t *testing.T, resource string, name string, cidr string, allowDuplicateHost string, gateway string, netReserveIPValue string, template string, ipVersion string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check Network on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector

		networkEntity := entities.Network{}
		networkEntity.Configuration = configuration
		networkEntity.CIDR = cidr
		networkEntity.IPVersion = ipVersion

		network, err := objMgr.GetNetwork(&networkEntity)
		if err != nil {
			msg := fmt.Sprintf("Getting Network %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		allowDuplicateHostProperty := getPropertyValue("allowDuplicateHost", network.Properties)
		gatewayProperty := getPropertyValue("gateway", network.Properties)
		ipAddress, err := objMgr.GetIPAddress(configuration, netReserveIPValue, entities.IPV4)
		state := getPropertyValue("state", ipAddress.Properties)

		if ipVersion == entities.IPV4 {
			if allowDuplicateHostProperty != allowDuplicateHost || gatewayProperty != gateway || network.Name != name {
				msg := fmt.Sprintf("Getting Network %s failed: %s. Expect allowDuplicateHost=%s gateway=%s in properties and name=%s, but received '%s' and name=%s", rs.Primary.ID, err, allowDuplicateHost, gateway, name, network.Properties, network.Name)
				log.Error(msg)
				return fmt.Errorf(msg)
			}

			if err != nil {
				msg := fmt.Sprintf("Getting reverse ip of Network %s failed: %s", rs.Primary.ID, err)
				log.Error(msg)
				return fmt.Errorf(msg)
			}

			if state != "RESERVED" {
				msg := fmt.Sprintf("Getting reverse ip of Network %s failed: %s. %s is not RESERVED", rs.Primary.ID, err, netReserveIPValue)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
			templateId := getPropertyValue("template", network.Properties)
			if template != "" && templateId == "" {
				msg := fmt.Sprintf("Assign %s template of Network %s failed", template, rs.Primary.ID)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
		} else if ipVersion == entities.IPV6 {
			if network.Name != name {
				msg := fmt.Sprintf("Getting Network %s failed: %s. Expect name=%s, but received name=%s", rs.Primary.ID, err, name, network.Name)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
		}
		return nil
	}
}

var blockNetResource1 = "block_record"
var testAccresourceBlockCreate = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4block" "%s" {
		configuration = "%s"
		name = "block1"
		parent_block = ""
		address = "30.0.0.0"
		cidr = "24"
		properties = ""
	  }`, server, blockNetResource1, configuration)

var netResource1 = "net_record"

// var netName1 = ""
var netName1 = "network1"
var netCIDR1 = "30.0.0.0/24"
var netGateway1 = "30.0.0.12"
var netReserveIP1 = "1"
var netReserveIPValue1 = "30.0.0.1"

// var netProperties1 = ""
var netProperties1 = "allowDuplicateHost=disable|"
var netAllowDuplicateHost1 = "disable"
var testAccResourceNetworkCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4network" "%s" {
		configuration = "%s"
		name = "%s"
		cidr = "%s"
		gateway = "%s"
		reserve_ip = %s
		properties = "%s"
		depends_on = [bluecat_ipv4block.%s]
		}`, testAccresourceBlockCreate, netResource1, configuration, netName1, netCIDR1, netGateway1, netReserveIP1, netProperties1, blockNetResource1)

var testAccResourceNetworkCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4network" "%s" {
		configuration = "%s"
		cidr = "%s"
		gateway = "%s"
		reserve_ip = %s
		properties = "%s"
		depends_on = [bluecat_ipv4block.%s]
		}`, testAccresourceBlockCreate, netResource1, configuration, netCIDR1, netGateway1, netReserveIP1, netProperties1, blockNetResource1)

var netName2 = "network2"
var netGateway2 = "30.0.0.15"

// var netProperties2 = ""
var netProperties2 = "allowDuplicateHost=enable|"
var netAllowDuplicateHost2 = "enable"
var testAccResourceNetworkUpdateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4network" "%s" {
		configuration = "%s"
		cidr = "%s"
		gateway = "%s"
		reserve_ip = %s
		properties = "%s"
		depends_on = [bluecat_ipv4block.%s]
	  }`, testAccresourceBlockCreate, netResource1, configuration, netCIDR1, netGateway2, netReserveIP1, netProperties2, blockNetResource1)

var testAccResourceNetworkUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4network" "%s" {
		configuration = "%s"
		name = "%s"
		cidr = "%s"
		gateway = "%s"
		reserve_ip = %s
		properties = "%s"
		depends_on = [bluecat_ipv4block.%s]
		}`, testAccresourceBlockCreate, netResource1, configuration, netName2, netCIDR1, netGateway2, netReserveIP1, netProperties2, blockNetResource1)

var netReserveIP2 = "1"
var netReserveIPValue2 = "30.0.0.2"
var templateName = "template1"
var netGateway3 = "30.0.0.1"
var testAccResourceNetworkCreateFullFieldWithTemplate = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4network" "%s" {
		configuration = "%s"
		name = "%s"
		cidr = "%s"
		gateway = "%s"
		reserve_ip = %s
		properties = "%s"
		template = "%s"
		depends_on = [bluecat_ipv4block.%s]
		}`, testAccresourceBlockCreate, netResource1, configuration, netName1, netCIDR1, netGateway3, netReserveIP2, netProperties1, templateName, blockNetResource1)

var nextNetResource1 = "next_net_record"

var nextNetName1 = "next network1"
var nextNetReserveIP1 = "1"
var nextNetParentBlock1 = "30.0.0.0/24"
var nextNetSize1 = "256"

var nextNetGatewayValue1 = "30.0.0.1"
var nextNetReserveIPValue1 = "30.0.0.2"
var nextNetCIDRValue1 = "30.0.0.0/24"

// var nextNetProperties1 = ""
var nextNetProperties1 = "allowDuplicateHost=disable|"
var nextNetAllowDuplicateHost1 = "disable"
var testAccResourceNextNetworkCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4network" "%s" {
		configuration = "%s"
		name = "%s"
		reserve_ip = %s
		properties = "%s"
		parent_block = "%s"
		size = %s
		depends_on = [bluecat_ipv4block.%s]
		}`, testAccresourceBlockCreate, nextNetResource1, configuration, nextNetName1, nextNetReserveIP1, nextNetProperties1, nextNetParentBlock1, nextNetSize1, blockNetResource1)

var nextNetName2 = "next network2"

var testAccResourceNextNetworkUpdateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4network" "%s" {
		configuration = "%s"
		name = "%s"
		reserve_ip = %s
		properties = "%s"
		parent_block = "%s"
		size = %s
		depends_on = [bluecat_ipv4block.%s]
		}`, testAccresourceBlockCreate, nextNetResource1, configuration, nextNetName2, nextNetReserveIP1, nextNetProperties1, nextNetParentBlock1, nextNetSize1, blockNetResource1)

var testAccResourceIPv6BlockCreate = fmt.Sprintf(
	`%s
	resource "bluecat_ipv6block" "ipv6_block_record_1" {
		configuration = "%s"
		name = "ipv6_block_name"
		parent_block = ""
		address = "2040:B041::"
		cidr = "64"
		properties = ""
		ip_version = "ipv6"
	  }`, server, configuration)

var testAccResourceIPv6NetworkCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv6network" "ipv6_net_record_1" {
		configuration = "%s"
		name = "ipv6_net_name"
		cidr = "2040:B041::/64"
		properties = ""
		ip_version = "ipv6"
		depends_on = [bluecat_ipv6block.ipv6_block_record_1]
		}`, testAccResourceIPv6BlockCreate, configuration)

var testAccResourceIPv6NetworkInsideUniqueLocal = fmt.Sprintf(
	`%s
	resource "bluecat_ipv6network" "ipv6_net_record_2" {
		configuration = "%s"
		name = "ipv6_net_name"
		cidr = "FC00::/64"
		properties = ""
		ip_version = "ipv6"
		}`, server, configuration)

// do not send name and properties option
var testAccResourceIPv6NetworkCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv6network" "ipv6_net_record_1_not_full_field" {
		configuration = "%s"
		cidr = "FC00::/64"
		ip_version = "ipv6"
		}`, server, configuration)

// do not send name and properties option
var testAccResourceIPv6NetworUpdateAddName = fmt.Sprintf(
	`%s
	resource "bluecat_ipv6network" "ipv6_net_record_1_not_full_field" {
		configuration = "%s"
		cidr = "FC00::/64"
		name = "new_name"
		ip_version = "ipv6"
		}`, server, configuration)
