package main
import (
	"fmt"
	"terraform-provider-bluecat/bluecat/logging"
	"terraform-provider-bluecat/bluecat"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"testing"
	"github.com/sirupsen/logrus"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var log logrus.Logger

func init() {
	log = *logging.GetLogger()
}

func init() {
	testAccProvider = bluecat.Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"bluecat": testAccProvider,
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
		server = "127.0.0.1"
		api_version = "1"
		transport = "http"
		port = "8000"
		username = "terraform"
		password = "123456"
	  }`)
