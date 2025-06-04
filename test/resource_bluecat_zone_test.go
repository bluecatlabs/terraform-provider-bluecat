package main

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func skipTest() (bool, error) {
	return true, nil
}

func TestAccResourceZone(t *testing.T) {
	configuration = "terraform_test"
	// create top zone with full fields and update
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceZoneCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneExists(t, fmt.Sprintf("bluecat_zone.%s", zoneResource1), zoneName1, zoneDeployable1, zoneProperties1),
				),
			},
			// update
			resource.TestStep{
				Config: testAccresourceZoneUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneExists(t, fmt.Sprintf("bluecat_zone.%s", zoneResource1), zoneName1, zoneDeployable2, zoneProperties1),
				),
			},
		},
	})
	// create sub zone with full fields and update
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceSubZoneCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneExists(t, fmt.Sprintf("bluecat_zone.%s", zoneResource2), zoneName2, zoneDeployable1, zoneProperties1),
				),
			},
			// update
			resource.TestStep{
				Config: testAccresourceSubZoneUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccZoneExists(t, fmt.Sprintf("bluecat_zone.%s", zoneResource2), zoneName2, zoneDeployable2, zoneProperties1),
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
			// return fmt.Errorf(msg)
		}
		_, err := objMgr.GetZone(configuration, view, rs.Primary.ID)
		if err == nil {
			msg := fmt.Sprintf("Zone %s is not removed", rs.Primary.ID)
			log.Error(msg)
			return fmt.Errorf("Zone %s is not removed", rs.Primary.ID)
		}
	}
	return nil
}

func testAccZoneExists(t *testing.T, resource string, zoneName string, deployable string, properties string) resource.TestCheckFunc {
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
			return fmt.Errorf("Getting Zone %s failed: %s", rs.Primary.ID, err)
		}
		if checkValidZone(objMgr, *zone, zoneName, deployable) == false {
			msg := fmt.Sprintf("Getting Zone %s failed: deployable property or list server_roles does not match", rs.Primary.ID)
			log.Error(msg)
			return fmt.Errorf("Getting Zone %s failed: deployable property or list server_roles does not match", rs.Primary.ID)
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

func checkValidZone(objMgr *utils.ObjectManager, zone entities.Zone, zoneName string, deployable string) bool {
	deployableProperty := utils.GetPropertyValue("deployable", zone.Properties)
	deployableValues := []string{"yes", "true", "1"}
	if contains(deployableValues, deployable) != contains(deployableValues, deployableProperty) {
		return false
	}

	return true
}

var zoneResource1 = "top_zone_org"
var zoneName1 = "top_zone.org"
var zoneDeployable1 = "true"
var zoneProperties1 = ""
var testAccresourceZoneCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_zone" "%s" {
		configuration = "terraform_test"
		view = "test"
		zone = "%s"
		deployable = "%s"
		properties = "%s"
		depends_on = [bluecat_view.view_test, bluecat_zone.zone_org]
		}`, GetTestEnvResources(), zoneResource1, zoneName1, zoneDeployable1, zoneProperties1)

var zoneDeployable2 = "false"
var testAccresourceZoneUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_zone" "%s" {
		configuration = "terraform_test"
		view = "test"
		zone = "%s"
		deployable = "%s"
		properties = "%s"
		depends_on = [bluecat_view.view_test, bluecat_zone.zone_org]
		}`, GetTestEnvResources(), zoneResource1, zoneName1, zoneDeployable2, zoneProperties1)

var zoneResource2 = "sub_zone"
var zoneName2 = "subzone.com"
var testAccresourceSubZoneCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_zone" "%s" {
		configuration = "terraform_test"
		view = "%s"
		zone = "%s"
		deployable = "%s"
		properties = "%s"
		depends_on = [bluecat_view.view_test]
		}`, GetTestEnvResources(), zoneResource2, view, zoneName2, zoneDeployable1, zoneProperties1)

var testAccresourceSubZoneUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_zone" "%s" {
		configuration = "terraform_test"
		view = "%s"
		zone = "%s"
		deployable = "%s"
		properties = "%s"
		depends_on = [bluecat_view.view_test]
		}`, GetTestEnvResources(), zoneResource2, view, zoneName2, zoneDeployable2, zoneProperties1)
