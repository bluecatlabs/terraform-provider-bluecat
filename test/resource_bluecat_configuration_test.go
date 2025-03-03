package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"terraform-provider-bluecat/bluecat/utils"
	"testing"
)

func TestAccResourceConfiguration(t *testing.T) {
	// create with full fields and update
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckConfigurationDestroy,
		Steps: []resource.TestStep{
			// create
			resource.TestStep{
				Config: testAccresourceConfigurationCreateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccConfigurationExists(t, fmt.Sprintf("bluecat_configuration.%s", confResource1), confName1, confDescriptionProperty1),
				),
			},
			// update
			resource.TestStep{
				Config: testAccresourceConfigurationUpdateFullField,
				Check: resource.ComposeTestCheckFunc(
					testAccConfigurationExists(t, fmt.Sprintf("bluecat_configuration.%s", confResource1), confName1, confDescriptionProperty2),
				),
			},
		},
	})
}

func testAccCheckConfigurationDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()
	connector := meta.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bluecat_configuration" {
			msg := fmt.Sprintf("There is an unexpected resource %s %s", rs.Primary.ID, rs.Type)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		_, err := objMgr.GetConfiguration("terraform_test_configuration_1")
		if err == nil {
			msg := fmt.Sprintf("Configuration %s is not removed", rs.Primary.ID)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
	}
	return nil
}

func testAccConfigurationExists(t *testing.T, resource string, name string, confDescriptionProperty string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		// check Configuration on BAM
		meta := testAccProvider.Meta()
		connector := meta.(*utils.Connector)
		objMgr := new(utils.ObjectManager)
		objMgr.Connector = connector
		conf, err := objMgr.GetConfiguration(name)
		if err != nil {
			msg := fmt.Sprintf("Getting configuration %s failed: %s", rs.Primary.ID, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		descriptionProperty := getPropertyValue("description", conf.Properties)
		if descriptionProperty != confDescriptionProperty {
			msg := fmt.Sprintf("Getting configuration %s failed: %s. Expect description='%s' in properties, but received '%s'", rs.Primary.ID, err, confDescriptionProperty, descriptionProperty)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		return nil
	}
}

var confResource1 = "conf_record"
var confName1 = "terraform_test_configuration_1"
var confProperties1 = "description=terraform testing config|"
var confDescriptionProperty1 = "terraform testing config"
var testAccresourceConfigurationCreateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_configuration" "%s" {
		name = "%s"
		properties = "%s"
	  }`, server, confResource1, confName1, confProperties1)

var confProperties2 = "description=updated config|"
var confDescriptionProperty2 = "updated config"
var testAccresourceConfigurationUpdateFullField = fmt.Sprintf(
	`%s
	resource "bluecat_configuration" "%s" {
		name = "%s"
		properties = "%s"
		}`, server, confResource1, confName1, confProperties2)
