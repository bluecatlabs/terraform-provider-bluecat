// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceParseId(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected ip_address/cidr", id)
	}
	address := parts[0]
	cidr := parts[1]

	return address, cidr, nil
}

func getIpVersion(d *schema.ResourceData, address string) string {
	ip := net.ParseIP(address)
	ip_from_id := net.ParseIP(d.Id())
	ipVersion := d.Get("ip_version").(string)
	if ipVersion == "" {
		if ip.To4() != nil || ip_from_id.To4() != nil {
			ipVersion = "ipv4"
		} else {
			ipVersion = "ipv6"
		}
	}
	return ipVersion
}

// ResourceBlock The IPv4 Block
func ResourceBlock() *schema.Resource {
	return &schema.Resource{
		Create: createBlock,
		Read:   getBlock,
		Update: updateBlock,
		Delete: deleteBlock,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the Block in the default Configuration if doesn't specify",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Block name",
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IPv4 Block's address",
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Block prefix length",
			},
			"parent_block": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parent Block. Specified to creating the child Block. THe Block in CIDR format",
			},
			"properties": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 Block's properties. Example: attribute=value|",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return checkDiffProperties(old, new)
				},
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Block IP version: ipv4 or ipv6",
			},
		},
		Importer: &schema.ResourceImporter{
			State: resourceImporter,
		},
	}
}

// createIP4Block Create the new IPv4 Block
func createBlock(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create Block %s", d.Get("address"))

	block := entities.Block{}
	block.InitBlock(d)

	_, err := strconv.Atoi(block.CIDR)
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", block.CIDR, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err = objMgr.CreateBlock(block)
	if err != nil {
		msg := fmt.Sprintf("Error creating Block (%s): %s", block.Address, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	log.Debugf("Completed to create Block %s", d.Get("address"))
	return getBlock(d, m)
}

// getBlock Get the IPv4 Block
func getBlock(d *schema.ResourceData, m interface{}) error {
	var address, cidrStr string
	var err error
	if d.Id() != "" {
		address, cidrStr, err = resourceServiceParseId(d.Id())
	}

	if err != nil {
		return err
	}
	log.Debugf("Beginning to get Block %s", d.Get("address"))
	configuration := d.Get("configuration").(string)

	if address == "" {
		address = d.Get("address").(string)
	}
	ipVersion := getIpVersion(d, address)
	if cidrStr == "" {
		cidrStr = d.Get("cidr").(string)
	}
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	block, err := objMgr.GetBlock(configuration, address, cidrStr, ipVersion)
	if err != nil {
		// check to see if some resource exist. If it does not exist, and it is in the plan - create that resource
		if d.Id() != "" {
			//d.Set("address", address)
			d.Set("ip_version", ipVersion)
			err := createBlock(d, m)
			if err != nil {
				msg := fmt.Sprintf("Something gone wrong %v", err)
				return fmt.Errorf(msg)
			}
		}
	}

	d.Set("name", block.Name)
	d.Set("properties", block.Properties)
	d.SetId(block.AddressCIDR())
	log.Debugf("Completed getting Block %s", d.Get("address"))
	return nil
}

// updateBlock Update the existing IPv4 Block
func updateBlock(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update Block %s", d.Get("address"))

	block := entities.Block{}
	address := d.Get("address").(string)
	d.Set("ip_version", getIpVersion(d, address))
	block.InitBlock(d)
	d.Set("ip_version", nil)

	_, err := strconv.Atoi(block.CIDR)
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", block.CIDR, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err = objMgr.UpdateBlock(block)
	if err != nil {
		msg := fmt.Sprintf("Error updating Block (%s): %s", block.Address, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	log.Debugf("Completed to update Block %s", d.Get("address"))
	return getBlock(d, m)
}

// deleteBlock Delete the IPv4 Block
func deleteBlock(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to Delete Block %s", d.Get("address"))
	configuration := d.Get("configuration").(string)
	address := d.Get("address").(string)
	cidrStr := d.Get("cidr").(string)
	ipVersion := d.Get("ip_version").(string)
	_, err := strconv.Atoi(cidrStr)
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	ipVersion = getIpVersion(d, address)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err = objMgr.DeleteBlock(configuration, address, cidrStr, ipVersion)
	if err != nil {
		msg := fmt.Sprintf("Delete Block %s/%s failed: %s", address, cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Deletion of Block complete ")
	return nil
}
