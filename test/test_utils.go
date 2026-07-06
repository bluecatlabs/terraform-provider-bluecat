package main

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCheckCNAMERecordDestroy(rs *terraform.ResourceState, objMgr *utils.ObjectManager, configuration string) (string, error) {
	var msg string
	var err error
	if rs.Type == "bluecat_ipv4block" {
		fmt.Println("Checking for block, ", rs.Primary.ID)
		if rs.Primary.ID != "" {
			cidr := strings.Split(rs.Primary.ID, "/")
			_, err := objMgr.GetBlock(configuration, cidr[0], cidr[1], entities.IPV4)
			if err == nil {
				msg = fmt.Sprintf("Block %s is not removed", rs.Primary.ID)
			}
		}
	} else if rs.Type == "bluecat_ipv4network" {
		fmt.Println("Checking for network, ", rs.Primary.ID)
		// cidr := strings.Split(rs.Primary.ID, "/")
		// Initne
		// _, err := objMgr.GetNetwork(configuration, cidr[0], cidr[1], entities.IPV6)
		// if err == nil {
		// 	msg := fmt.Sprintf("Block %s is not removed", rs.Primary.ID)
		// 	log.Error(msg)
		// 	return fmt.Errorf("Block %s is not removed", rs.Primary.ID)
		// }
	} else if rs.Type == "bluecat_zone" {
		fmt.Println("Checking for zone, ", rs.Primary.ID)
	} else if rs.Type == "bluecat_view" {
		fmt.Println("Checking for view, ", rs.Primary.ID)

	}
	return msg, err
}
