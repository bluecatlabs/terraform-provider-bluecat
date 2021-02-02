// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceNetwork The IPv4 Network
func ResourceNetwork() *schema.Resource {

	return &schema.Resource{
		Create: createNetwork,
		Read:   getNetwork,
		Update: updateNetwork,
		Delete: deleteNetwork,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the Network in the default Configuration if doesn't specify",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Network name",
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network address in CIDR format",
			},
			"reserve_ip": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Reserves the number of IP's for later use",
			},
			"gateway": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Give the IP you want to reserve for gateway, by default the first IP gets reserved for gateway",
			},
			"properties": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 Network's properties. Example: attribute=value|",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return checkDiffProperties(old, new)
				  },
			},
		},
	}
}

// createNetwork Create the new IPv4 Network
func createNetwork(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create Network %s", d.Get("cidr"))
	configuration := d.Get("configuration").(string)
	name := d.Get("name").(string)
	cidr := d.Get("cidr").(string)
	numReserved := d.Get("reserve_ip").(int)
	gateway := d.Get("gateway").(string)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	addrCIDR := strings.Split(cidr, "/")

	// Get the block
	block, err := objMgr.GetBlock(configuration, addrCIDR[0], "0")
	if err != nil {
		msg := fmt.Sprintf("Failed to getting the IPv4 Block for (%s): %s", cidr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	_, err = objMgr.CreateNetwork(configuration, block.AddressCIDR(), name, cidr, gateway, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating Network (%s): %s", cidr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	log.Debugf("Successful to create Network %s", d.Get("cidr"))
	if numReserved > 0 {
		log.Debugf("Reserving %d IP Addresses on the Network %s", numReserved, d.Get("cidr"))
		for i := 0; i < numReserved; i++ {
			_, err = objMgr.ReserveIPAddress(configuration, addrCIDR[0])
			if err != nil {
				msg := fmt.Sprintf("Reservation IP Address failed in network %s:%s", cidr, err)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
		}
	}

	return getNetwork(d, m)
}

// getNetwork Get the IPv4 Network
func getNetwork(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get Network %s", d.Get("cidr"))
	configuration := d.Get("configuration").(string)
	cidr := d.Get("cidr").(string)
	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	network, err := objMgr.GetNetwork(configuration, cidr)
	if err != nil {
		msg := fmt.Sprintf("Getting Network %s failed: %s", cidr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(network.CIDR)
	d.Set("name", network.Name)
	d.Set("properties", network.Properties)
	log.Debugf("Completed getting Network %s", d.Get("cidr"))
	return nil
}

// updateNetwork Update the existing IPv4 Network
func updateNetwork(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update Network %s", d.Get("cidr"))
	configuration := d.Get("configuration").(string)
	name := d.Get("name").(string)
	cidr := d.Get("cidr").(string)
	gateway := d.Get("gateway").(string)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.UpdateNetwork(configuration, name, cidr, gateway, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating Network (%s): %s", cidr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	log.Debugf("Completed to update Network %s", d.Get("cidr"))
	return getNetwork(d, m)
}

// deleteNetwork Delete the IPv4 Network
func deleteNetwork(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete Network %s", d.Get("cidr"))
	configuration := d.Get("configuration").(string)
	cidr := d.Get("cidr").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteNetwork(configuration, cidr)
	if err != nil {
		msg := fmt.Sprintf("Delete Network %s failed: %s", cidr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Deletion of Network complete ")
	return nil
}
