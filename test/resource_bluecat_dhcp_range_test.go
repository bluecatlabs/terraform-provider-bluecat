package main

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceDHCPRange(t *testing.T) {
	// create with full fields and update without some optional fields
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDHCPRangeDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccResourceDHCPRangeCreateNotTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccDHCPRangeExists(t, fmt.Sprintf("bluecat_dhcp_range.%s", dhcpRangeResource), dhcpRangeNetwork, dhcpRangeStart, dhcpRangeEnd, dhcpRangeTemplateName),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDHCPRangeDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccResourceDHCPRangeCreateWithTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccDHCPRangeExists(t, fmt.Sprintf("bluecat_dhcp_range.%s", dhcpRangeResource2), dhcpRangeNetwork, dhcpRangeStart2, dhcpRangeEnd2, dhcpRangeTemplateName2),
				),
			},
			// update
			resource.TestStep{
				Config: testAccResourceDHCPRangeUpdateWithTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccDHCPRangeExists(t, fmt.Sprintf("bluecat_dhcp_range.%s", dhcpRangeResource2), dhcpRangeNetwork, dhcpRangeStart2, dhcpRangeEnd2, dhcpRangeTemplateName3),
				),
			},
		},
	})
}

func testAccCheckDHCPRangeDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "bluecat_dhcp_range" {
			listParam := strings.Split(rs.Primary.ID, "-")
			start, end := listParam[0], listParam[1]
			_, err := objMgr.GetDHCPRange(configuration, dhcpRangeNetwork, start, end)
			if err == nil {
				msg := fmt.Sprintf("DHCP Range %s is not removed", rs.Primary.ID)
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

func testAccDHCPRangeExists(t *testing.T, resource string, network string, start string, end string, template string) resource.TestCheckFunc {
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
		dhcpRange, err := objMgr.GetDHCPRange(configuration, network, start, end)
		if err != nil {
			msg := fmt.Sprintf("Getting DHCP Range %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		templateId := getPropertyValue("template", dhcpRange.Properties)
		if template != "" && templateId == "" {
			msg := fmt.Sprintf("Assign %s template of DHCP Range %s failed %s", template, rs.Primary.ID, templateId)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		return nil
	}
}

var dhcpRangeResource = "dhcp_range"

var dhcpRangeStart = "1.1.0.5"
var dhcpRangeEnd = "1.1.0.10"
var dhcpRangeNetwork = "1.1.0.0/16"
var dhcpRangeProperties = ""
var dhcpRangeTemplateName = ""
var testAccResourceDHCPRangeCreateNotTemplate = fmt.Sprintf(
	`%s
	resource "bluecat_dhcp_range" %s {
		configuration = "%s"
		start = "%s"
		end = "%s"
		network = "%s"
		properties = "%s"
		}`, server, dhcpRangeResource, configuration, dhcpRangeStart, dhcpRangeEnd, dhcpRangeNetwork, dhcpRangeProperties)

var dhcpRangeResource2 = "dhcp_range_2"

var dhcpRangeStart2 = "1.1.0.12"
var dhcpRangeEnd2 = "1.1.0.15"
var dhcpRangeTemplateName2 = "template1"
var testAccResourceDHCPRangeCreateWithTemplate = fmt.Sprintf(
	`%s
	resource "bluecat_dhcp_range" %s {
		configuration = "%s"
		start = "%s"
		end = "%s"
		network = "%s"
		properties = "%s"
		template = "%s"
		}`, server, dhcpRangeResource2, configuration, dhcpRangeStart2, dhcpRangeEnd2, dhcpRangeNetwork, dhcpRangeProperties, dhcpRangeTemplateName2)

var dhcpRangeTemplateName3 = ""
var testAccResourceDHCPRangeUpdateWithTemplate = fmt.Sprintf(
	`%s
	resource "bluecat_dhcp_range" %s {
		configuration = "%s"
		start = "%s"
		end = "%s"
		network = "%s"
		properties = "%s"
		template = "%s"
		}`, server, dhcpRangeResource2, configuration, dhcpRangeStart2, dhcpRangeEnd2, dhcpRangeNetwork, dhcpRangeProperties, dhcpRangeTemplateName3)
				