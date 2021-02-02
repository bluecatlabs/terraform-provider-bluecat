package main

import (
	"fmt"
	"strings"
	"testing"
	"terraform-provider-bluecat/bluecat/utils"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceBlock(t *testing.T) {
	// create with full fields and update without some optional fields, then create sub-block
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccResourceBlockCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccBlockExists(t, fmt.Sprintf("bluecat_ipv4block.%s", blockResource1), blockName1, blockAddress1, blockCIDR1, blockAllowDuplicateHostProperty1),
				),
			},
			// update
			resource.TestStep{
				Config: testAccResourceBlockUpdateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccBlockExists(t, fmt.Sprintf("bluecat_ipv4block.%s", blockResource1), blockName2, blockAddress1, blockCIDR1, blockAllowDuplicateHostProperty2),
				),
			},
			// create sub-block
			resource.TestStep{
				Config: testAccResourceSubBlockCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccBlockExists(t, fmt.Sprintf("bluecat_ipv4block.%s", blockResource3), blockName3, blockAddress3, blockCIDR3, blockAllowDuplicateHostProperty3),
				),
			},
		},
	})
	// create without some optional fields and update with full fields
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccResourceBlockCreateNotFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccBlockExists(t, fmt.Sprintf("bluecat_ipv4block.%s", blockResource1), blockName1, blockAddress1, blockCIDR1, blockAllowDuplicateHostProperty1),
				),
			},
			resource.TestStep{
				Config: testAccResourceBlockUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccBlockExists(t, fmt.Sprintf("bluecat_ipv4block.%s", blockResource1), blockName2, blockAddress1, blockCIDR1,blockAllowDuplicateHostProperty2),
				),
			},
		},
	})
}

func testAccCheckBlockDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bluecat_ipv4block" {
			msg := fmt.Sprintf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		cidr := strings.Split(rs.Primary.ID, "/")
		_, err := objMgr.GetBlock(configuration, cidr[0], cidr[1])
		if err == nil {
			msg := fmt.Sprintf("Block %s is not removed", rs.Primary.ID)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
	}
	return nil
}

func testAccBlockExists(t *testing.T, resource string, name string, address string, cidr string, blockAllowDuplicateHostProperty string ) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check Block on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		block, err := objMgr.GetBlock(configuration, address, cidr)
		if err != nil {
			msg := fmt.Sprintf("Getting block %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		allowDuplicateHostProperty := getPropertyValue("allowDuplicateHost", block.Properties)
		if allowDuplicateHostProperty != blockAllowDuplicateHostProperty || block.Name != name {
			msg := fmt.Sprintf("Getting block %s failed: %s. Expect allowDuplicateHost=%s in properties and name=%s, but received '%s' and name=%s", rs.Primary.ID, err, blockAllowDuplicateHostProperty, name, block.Properties, block.Name)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}
}

var blockResource1 = "block_record"
var blockName1 = "block1"
var blockParent1 = ""
var blockAddress1 = "30.0.0.0"
var blockCIDR1 = "24"
var blockProperties1 = "allowDuplicateHost=disable|"
var blockAllowDuplicateHostProperty1 = "disable"
var testAccResourceBlockCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4block" "%s" {
		configuration = "%s"
		name = "%s"
		parent_block = "%s"
		address = "%s"
		cidr = "%s"
		properties = "%s"
	  }`, server, blockResource1, configuration, blockName1, blockParent1, blockAddress1, blockCIDR1, blockProperties1)

var testAccResourceBlockCreateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4block" "%s" {
		configuration = "%s"
		name = "%s"
		address = "%s"
		cidr = "%s"
		properties = "%s"
		}`, server, blockResource1, configuration, blockName1, blockAddress1, blockCIDR1, blockProperties1)

var blockName2 = "block2"
var blockProperties2 = "allowDuplicateHost=enable|"
var blockAllowDuplicateHostProperty2 = "enable"
var testAccResourceBlockUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4block" "%s" {
		configuration = "%s"
		name = "%s"
		parent_block = "%s"
		address = "%s"
		cidr = "%s"
		properties = "%s"
		}`, server, blockResource1, configuration, blockName2, blockParent1, blockAddress1, blockCIDR1, blockProperties2)

var testAccResourceBlockUpdateNotFullField = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4block" "%s" {
		configuration = "%s"
		name = "%s"
		address = "%s"
		cidr = "%s"
		properties = "%s"
		}`, server, blockResource1, configuration, blockName2, blockAddress1, blockCIDR1, blockProperties2)

var blockResource3 = "block_record1"
var blockName3 = "block3"
var blockAddress3 = "30.0.0.0"
var blockCIDR3 = "25"
var blockProperties3 = "allowDuplicateHost=enable|"
var blockAllowDuplicateHostProperty3 = "enable"
var testAccResourceSubBlockCreate = fmt.Sprintf(
	`%s
	resource "bluecat_ipv4block" "%s" {
		configuration = "%s"
		name = "%s"
		parent_block = "%s"
		address = "%s"
		cidr = "%s"
		properties = "%s"
		depends_on = [bluecat_ipv4block.%s]
		}`, testAccResourceBlockUpdateFullField, blockResource3, configuration, blockName3, blockName2, blockAddress3, blockCIDR3, blockProperties3, blockResource1)
