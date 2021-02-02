// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strconv"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceBlock The IPv4 Block
func ResourceBlock() *schema.Resource {
	return &schema.Resource{
		Create: createIP4Block,
		Read:   getIP4Block,
		Update: updateIP4Block,
		Delete: deleteIP4Block,

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
		},
	}
}

// createIP4Block Create the new IPv4 Block
func createIP4Block(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create Block %s", d.Get("address"))
	configuration := d.Get("configuration").(string)
	name := d.Get("name").(string)
	parentBlock := d.Get("parent_block").(string)
	address := d.Get("address").(string)
	cidrStr := d.Get("cidr").(string)
	_, err := strconv.Atoi(cidrStr)
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err = objMgr.CreateBlock(configuration, name, address, cidrStr, parentBlock, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating Block (%s): %s", address, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	log.Debugf("Completed to create Block %s", d.Get("address"))
	return getIP4Block(d, m)
}

// getIP4Block Get the IPv4 Block
func getIP4Block(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get Block %s", d.Get("address"))
	configuration := d.Get("configuration").(string)
	address := d.Get("address").(string)
	cidrStr := d.Get("cidr").(string)
	_, err := strconv.Atoi(cidrStr)
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	block, err := objMgr.GetBlock(configuration, address, cidrStr)
	if err != nil {
		msg := fmt.Sprintf("Getting Block %s/%s failed: %s", address, cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(block.AddressCIDR())
	d.Set("name", block.Name)
	d.Set("properties", block.Properties)
	log.Debugf("Completed getting Block %s", d.Get("address"))
	return nil
}

// updateIP4Block Update the existing IPv4 Block
func updateIP4Block(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update Block %s", d.Get("address"))
	configuration := d.Get("configuration").(string)
	name := d.Get("name").(string)
	parentBlock := d.Get("parent_block").(string)
	address := d.Get("address").(string)
	cidrStr := d.Get("cidr").(string)
	_, err := strconv.Atoi(cidrStr)
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err = objMgr.UpdateBlock(configuration, name, address, cidrStr, parentBlock, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating Block (%s): %s", address, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	log.Debugf("Completed to update Block %s", d.Get("address"))
	return getIP4Block(d, m)
}

// deleteIP4Block Delete the IPv4 Block
func deleteIP4Block(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to Delete Block %s", d.Get("address"))
	configuration := d.Get("configuration").(string)
	address := d.Get("address").(string)
	cidrStr := d.Get("cidr").(string)
	_, err := strconv.Atoi(cidrStr)
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err = objMgr.DeleteBlock(configuration, address, cidrStr)
	if err != nil {
		msg := fmt.Sprintf("Delete Block %s/%s failed: %s", address, cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Deletion of Block complete ")
	return nil
}
