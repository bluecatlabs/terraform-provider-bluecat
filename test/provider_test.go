package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sirupsen/logrus"
	"terraform-provider-bluecat/bluecat"
	"terraform-provider-bluecat/bluecat/logging"
	"testing"
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

var configuration = "terraform_test"
var view = "test"
var zone = "example.com"
var server = fmt.Sprintf(
	`provider "bluecat" {
		server = "10.244.85.185"
		api_version = "1"
		transport = "http"
		port = "80"
		username = "admin"
		password = "admin"
	  }`)
