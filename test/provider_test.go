package main

import (
	"fmt"
	"terraform-provider-bluecat/bluecat"
	"terraform-provider-bluecat/bluecat/logging"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sirupsen/logrus"
)

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider
var log logrus.Logger

func init() {
	log = *logging.GetLogger()
	testAccProvider = bluecat.Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"bluecat": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := bluecat.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func GetTestEnvResources() string {
	return fmt.Sprintf(
		`%s
		%s
		%s
		%s
		%s
		%s
		%s`, server, blockTestResource, networkTestResource, viewTestResource, zoneTestResource, zoneTestResource2, ipv6NetworkTestResource)
}

var configuration = "terraform_test"
var view = "test"
var zone = "example.com"
var server = fmt.Sprintf(
	`provider "bluecat" {
		server = "10.244.12.9"
		api_version = "1"
		transport = "http"
		port = "80"
		username = "admin"
		password = "admin"
	  }`)
var viewTestResource = fmt.Sprintf(
	`resource "bluecat_view" "view_test" {
		name = "%s"
		configuration = "%s"
	}`, view, configuration)
var zoneTestResource = fmt.Sprintf(
	`resource "bluecat_zone" "sub_zone_test" {
		configuration = "%s"
		view          = "%s"
		zone          = "%s"
		deployable    = false
		depends_on = [bluecat_view.view_test]
	}`, configuration, view, zone)
var zoneTestResource2 = fmt.Sprintf(
	`resource "bluecat_zone" "zone_org" {
		configuration = "%s"
		view          = "%s"
		zone          = "org"
		deployable    = false
		depends_on = [bluecat_view.view_test]
	}`, configuration, view)
var blockTestResource = fmt.Sprintf(
	`resource "bluecat_ipv4block" "block_test" {
			configuration = "%s"
			name = "test_block"
			address = "1.1.0.0"
			cidr = "16"
		}`, configuration)
var networkTestResource = fmt.Sprintf(
	`resource "bluecat_ipv4network" "network_test" {
			configuration = "%s"
			name = "network_test"
			cidr = "1.1.0.0/16"
			depends_on = [bluecat_ipv4block.block_test]
		}`, configuration)
var ipv6BlockTestResource = fmt.Sprintf(
	`resource "bluecat_ipv6block" "ipv6_block_test" {
      configuration = "%s"
      name = "test_ipv6"
      parent_block = ""
      address = "2003::"
      cidr = "3"
      properties = ""
      ip_version = "ipv6"
    }`, configuration)
var ipv6NetworkTestResource = fmt.Sprintf(
	`resource "bluecat_ipv6network" "ipv6_network_test" {
      configuration = "%s"
      name = "test_ipv6"
      cidr = "2003:1000::/64"
      ip_version = "ipv6"
      properties = ""
    }`, configuration)
