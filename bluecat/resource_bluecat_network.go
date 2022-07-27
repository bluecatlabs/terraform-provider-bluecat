// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
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
				Optional:    true,
				Description: "The network address in CIDR format",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new != old && new == "" {
						return true
					}
					return new == old
				},
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
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 Template",
			},
			"parent_block": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parent block of the network in CIDR format",
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The size of the network expressed in the power of 2",
			},
			"allocated_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The allocated id of the next available network",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != ""
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
	template := d.Get("template").(string)
	parentBlock := d.Get("parent_block").(string)
	size := d.Get("size").(string)
	allocatedId := d.Get("allocated_id").(string)

	var networkAddress string

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	if cidr != "" {
		// Create specified network

		networkAddress = strings.Split(cidr, "/")[0]

		block, err := objMgr.GetBlock(configuration, networkAddress, "0")
		if err != nil {
			msg := fmt.Sprintf("Failed to getting the IPv4 Block for (%s): %s", cidr, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		_, err = objMgr.CreateNetwork(configuration, block.AddressCIDR(), name, cidr, gateway, properties, template)
		if err != nil {
			msg := fmt.Sprintf("Error creating Network (%s): %s", cidr, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		log.Debugf("Successful to create Network %s", d.Get("cidr"))

	} else {
		// Create next available network

		if parentBlock == "" {
			msg := "'parent_block' is a required property to get next available network"
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		if allocatedId == "" {
			msg := "'allocated_id' is a required property to get next available network"
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		sizeNumber, err := strconv.Atoi(size)
		if err != nil && !isPowerOfTwo(sizeNumber) {
			msg := "'size' is a required property and must be power of 2 to get next available network"
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		block, err := objMgr.GetBlock(configuration, parentBlock, "0")
		if err != nil {
			msg := fmt.Sprintf("Failed to getting the IPv4 Block for (%s): %s", cidr, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		_, ref, err := objMgr.CreateNextAvailableNetwork(configuration, block.AddressCIDR(), name, gateway, properties, template, size, allocatedId)
		if err != nil {
			msg := fmt.Sprintf("Error creating next available Network of Block(%s): %s", parentBlock, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		log.Debugf("Successful to create next available Network of Block %s", d.Get("parent_block"))

		cidr = getObjectFieldValue("CIDR", ref)
		d.Set("cidr", cidr)

		networkAddress = strings.Split(cidr, "/")[0]
	}

	if numReserved > 0 {
		log.Debugf("Reserving %d IP Addresses on the Network %s", numReserved, cidr)
		for i := 0; i < numReserved; i++ {
			_, err := objMgr.ReserveIPAddress(configuration, networkAddress)
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

	parentBlock := d.Get("parent_block").(string)
	allocatedId := d.Get("allocated_id").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	var network *entities.Network
	var err error

	if cidr != "" {
		network, err = objMgr.GetNetwork(configuration, cidr)
		if err != nil {
			msg := fmt.Sprintf("Getting Network %s failed: %s", cidr, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
	} else if allocatedId != "" && parentBlock != "" {
		network, err = objMgr.GetNetworkByAllocatedId(configuration, parentBlock, allocatedId)
		if err != nil {
			msg := fmt.Sprintf("Getting Network in block %s failed: %s", parentBlock, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
	}

	d.SetId(network.CIDR)
	d.Set("cidr", network.CIDR)
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

func isPowerOfTwo(number int) bool {
	return (number != 0) && (number&(number-1) == 0)
}

func getObjectFieldValue(fieldName, ref string) (val string) {

	object := entities.Network{}
	json.Unmarshal([]byte(ref), &object)

	r := reflect.ValueOf(&object)
	f := reflect.Indirect(r).FieldByName(fieldName)
	if (f != reflect.Value{} && f.String() != "") {
		return f.String()
	}

	return getPropertyValue(fieldName, object.Properties)
}
