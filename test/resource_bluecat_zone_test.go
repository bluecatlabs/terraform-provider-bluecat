package main

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceZone(t *testing.T) {
	// create top zone with full fields and update
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceZoneCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneExists(t, fmt.Sprintf("bluecat_zone.%s", zoneResource1), zoneName1, zoneDeployable1, zoneServerRoles1, zoneProperties1),
				),
			},
			// update
			resource.TestStep{
				Config: testAccresourceZoneUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneExists(t, fmt.Sprintf("bluecat_zone.%s", zoneResource1), zoneName1, zoneDeployable2, zoneServerRoles2, zoneProperties1),
				),
			},
		},
	})
	// create sub zone with full fields and update
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceSubZoneCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneExists(t, fmt.Sprintf("bluecat_zone.%s", zoneResource2), zoneName2, zoneDeployable1, zoneServerRoles1, zoneProperties1),
				),
			},
			// update
			resource.TestStep{
				Config: testAccresourceSubZoneUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneExists(t, fmt.Sprintf("bluecat_zone.%s", zoneResource2), zoneName2, zoneDeployable2, zoneServerRoles2, zoneProperties1),
				),
			},
		},
	})
}

func testAccCheckZoneDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bluecat_zone" {
			msg := fmt.Sprintf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		_, err := objMgr.GetZone(configuration, view, rs.Primary.ID)
		if err == nil {
			msg := fmt.Sprintf("Zone %s is not removed", rs.Primary.ID)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
	}
	return nil
}

func testAccZoneExists(t *testing.T, resource string, zoneName string, deployable string, serverRoles []string, properties string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check Zone on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		zone, err := objMgr.GetZone(configuration, view, zoneName)
		if err != nil {
			msg := fmt.Sprintf("Getting Zone %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		if checkValidZone(objMgr, *zone, zoneName, deployable, serverRoles) == false {
			msg := fmt.Sprintf("Getting Zone %s failed: deployable property or list server_roles does not match", rs.Primary.ID)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func checkValidZone(objMgr *utils.ObjectManager, zone entities.Zone, zoneName string, deployable string, serverRoles []string) bool {
	deployableProperty := getPropertyValue("deployable", zone.Properties)
	deployableValues := []string{"yes", "true", "1"}
	if contains(deployableValues, deployable) != contains(deployableValues, deployableProperty) {
		return false
	}

	for _, serverRole := range serverRoles {
		prop := strings.Split(serverRole, ",")
		serverFQDN := strings.Trim(prop[1], " ")
		_, err := objMgr.GetDeploymentRole(configuration, view, zoneName, serverFQDN)
		if err != nil {
			return false
		}
	}

	return true
}

var zoneResource1 = "top_zone"
var zoneName1 = "org"
var zoneDeployable1 = "true"
var zoneServerRolesRaw1 = "[\"primary, server1\", \"secondary, server2\"]"
var zoneServerRoles1 = []string{"primary, server1", "secondary, server2"}
var zoneProperties1 = ""
var testAccresourceZoneCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_zone" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		deployable = "%s"
		server_roles = %s
		properties = "%s"
		}`, server, zoneResource1, configuration, view, zoneName1, zoneDeployable1, zoneServerRolesRaw1, zoneProperties1)

var zoneDeployable2 = "false"
var zoneServerRolesRaw2 = "[\"secondary, server1\", \"primary, server3\"]"
var zoneServerRoles2 = []string{"secondary, server1", "primary, server3"}
var testAccresourceZoneUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_zone" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		deployable = "%s"
		server_roles = %s
		properties = "%s"
		}`, server, zoneResource1, configuration, view, zoneName1, zoneDeployable2, zoneServerRolesRaw2, zoneProperties1)

var zoneResource2 = "sub_zone"
var zoneName2 = "subzone.com"
var testAccresourceSubZoneCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_zone" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		deployable = "%s"
		server_roles = %s
		properties = "%s"
		}`, server, zoneResource2, configuration, view, zoneName2, zoneDeployable1, zoneServerRolesRaw1, zoneProperties1)

var testAccresourceSubZoneUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_zone" "%s" {
		configuration = "%s"
		view = "%s"
		zone = "%s"
		deployable = "%s"
		server_roles = %s
		properties = "%s"
		}`, server, zoneResource2, configuration, view, zoneName2, zoneDeployable2, zoneServerRolesRaw2, zoneProperties1)
