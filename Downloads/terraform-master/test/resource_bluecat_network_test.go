package main

import (
	"testing"
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccResourceNetwork(t *testing.T) {
	// create with full fields and update without some optional fields
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccResourceNetworkCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", netResource1), netName1, netCIDR1, netAllowDuplicateHost1, netGateway1, netReserveIPValue1),
				),
			},
			// // update
			resource.TestStep{
				Config: testAccResourceNetworkUpdateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", netResource1), "", netCIDR1, netAllowDuplicateHost2, netGateway2, netReserveIPValue1),
				),
			},
		},
	})
	// create without some optional fields and update with full fields
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccResourceNetworkCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", netResource1), "", netCIDR1, netAllowDuplicateHost1, netGateway1, netReserveIPValue1),
				),
			},
			// update
			resource.TestStep{
				Config: testAccResourceNetworkUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccNetworkExists(t, fmt.Sprintf("bluecat_ipv4network.%s", netResource1), netName2, netCIDR1, netAllowDuplicateHost2, netGateway2, netReserveIPValue1),
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
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "bluecat_ipv4network" {
			_, err := objMgr.GetNetwork(configuration, rs.Primary.ID)
			if err == nil {
				msg := fmt.Sprintf("Network %s is not removed", rs.Primary.ID)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
			
		} else if rs.Type == "bluecat_ipv4block" {
			cidr := strings.Split(rs.Primary.ID, "/")
			_, err := objMgr.GetBlock(configuration, cidr[0], cidr[1])
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

func testAccNetworkExists(t *testing.T, resource string, name string, cidr string, allowDuplicateHost string, gateway string, netReserveIPValue string) resource.TestCheckFunc {
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
		network, err := objMgr.GetNetwork(configuration, cidr)
		if err != nil {
			msg := fmt.Sprintf("Getting Network %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		allowDuplicateHostProperty := getPropertyValue("allowDuplicateHost", network.Properties)
		gatewayPropertiy := getPropertyValue("gateway", network.Properties)
		if allowDuplicateHostProperty != allowDuplicateHost || gatewayPropertiy != gateway || network.Name != name {
			msg := fmt.Sprintf("Getting Network %s failed: %s. Expect allowDuplicateHost=%s gateway=%s in properties and name=%s, but received '%s' and name=%s", rs.Primary.ID, err, allowDuplicateHost, gateway, name, network.Properties, network.Name)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		ipAddress, err := objMgr.GetIPAddress(configuration, netReserveIPValue)
		state := getPropertyValue("state", ipAddress.Properties)
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
