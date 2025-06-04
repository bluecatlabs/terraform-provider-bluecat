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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceNetwork The IPv4/IPv6 Network
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
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network IP version: ipv4 or ipv6",
			},
		},
		Importer: &schema.ResourceImporter{
			State: resourceImporter,
		},
	}
}

// createNetwork Create the new IPv4/IPv6 Network
func createNetwork(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	network := entities.Network{}
	if !network.InitNetwork(d) {
		log.Error(network.InitError)
		return fmt.Errorf(network.InitError)
	}
	log.Debugf("Beginning to create Network %s", network.CIDR)

	numReserved := d.Get("reserve_ip").(int)

	var networkAddress string
	if network.CIDR != "" {
		// Create specified network

		networkAddress = strings.Split(network.CIDR, "/")[0]

		var parentBlockCidrNotation string
		block, err := objMgr.GetBlock(network.Configuration, networkAddress, "0", network.IPVersion)
		if block.IPVersion == entities.IPV6 {
			parentBlockCidrNotation = block.GetIPv6BlockFromPropsPrefix()
		}
		if err != nil {
			msg := fmt.Sprintf("Failed to getting the IPv4 Block for (%s): %s", network.CIDR, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		if network.IPVersion == "" || network.IPVersion == entities.IPV4 {
			network.IPVersion = entities.IPV4
			network.BlockAddr = block.AddressCIDR()
		} else if network.IPVersion == entities.IPV6 {
			// we use this since 2000::/3 or FC00::/6 is a default block
			// if there is no any specific block under that block
			network.BlockAddr = parentBlockCidrNotation
		}

		_, err = objMgr.CreateNetwork(network)
		if err != nil {
			msg := fmt.Sprintf("Error creating Network (%s): %s", network.CIDR, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		log.Debugf("Successful to create Network %s", network.CIDR)

	} else if network.IPVersion != entities.IPV6 {
		// Create next available network

		if network.ParentBlock == "" {
			msg := "'parent_block' is a required property to get next available network"
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		sizeNumber, err := strconv.Atoi(network.Size)
		if err != nil && !isPowerOfTwo(sizeNumber) {
			msg := "'size' is a required property and must be power of 2 to get next available network"
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		if strings.Contains(network.ParentBlock, "/") {
			network.ParentBlock = strings.Split(network.ParentBlock, "/")[0]
		}
		block, err := objMgr.GetBlock(network.Configuration, network.ParentBlock, "0", network.IPVersion)
		if err != nil {
			msg := fmt.Sprintf("Failed to getting the IPv4 Block for (%s): %s", network.CIDR, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		network.BlockAddr = block.AddressCIDR()

		_, ref, err := objMgr.CreateNextAvailableNetwork(network)
		if err != nil {
			msg := fmt.Sprintf("Error creating next available Network of Block(%s): %s", network.ParentBlock, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		log.Debugf("Successful to create next available Network of Block %s", d.Get("parent_block"))

		network.CIDR = getObjectFieldValue("CIDR", ref)
		d.Set("cidr", network.CIDR)

		networkAddress = strings.Split(network.CIDR, "/")[0]
	}

	if network.IPVersion == entities.IPV4 {
		// numReserved is now only used for ipv4
		// TODO: Develop this feature in some of the next releases
		if numReserved > 0 {
			log.Debugf("Reserving %d IP Addresses on the Network %s", numReserved, network.CIDR)
			for i := 0; i < numReserved; i++ {
				_, err := objMgr.ReserveIPAddress(network.Configuration, networkAddress, network.IPVersion)
				if err != nil {
					msg := fmt.Sprintf("Reservation IP Address failed in network %s:%s", network.CIDR, err)
					log.Error(msg)
					return fmt.Errorf(msg)
				}
			}
		}
	}
	return getNetwork(d, m)
}

// getNetwork Get the IPv4/IPv6 Network
func getNetwork(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	var address, cidr, cidrStr, ipVersion string
	var err error
	if d.Id() != "" {
		address, cidrStr, err = resourceServiceParseId(d.Id())
		cidr = fmt.Sprintf("%s/%s", address, cidrStr)
	}

	if err != nil {
		return err
	}
	log.Debugf("Beginning to get Network %s", d.Get("cidr"))
	configuration := d.Get("configuration").(string)
	if cidr == "" {
		cidr = d.Get("cidr").(string)
	}
	d.Set("cidr", cidr)
	//cidrStr, err = strconv.Atoi(cidrStr)

	parentBlock := d.Get("parent_block").(string)
	allocatedId := d.Get("allocated_id").(string)
	if parentBlock != "" {
		ipVersion = getIpVersion(d, strings.Split(parentBlock, "/")[0])
	} else {
		if cidr != "" {
			address = strings.Split(cidr, "/")[0]
		}
		ipVersion = getIpVersion(d, address)
	}
	if ipVersion == entities.IPV6 {
		d.Set("ip_version", ipVersion)
	}

	networkEntity := entities.Network{}
	if !networkEntity.InitNetwork(d) {
		log.Error(networkEntity.InitError)
		return fmt.Errorf(networkEntity.InitError)
	}

	var network *entities.Network
	var getNetworkError error
	if networkEntity.CIDR != "" {
		network, getNetworkError = objMgr.GetNetwork(&networkEntity)
		if getNetworkError != nil {
			// check to see if some resource exist. If it does not exist, and it is in the plan - create that resource
			if d.Id() != "" {
				e := createNetwork(d, m)
				if e != nil {
					msg := fmt.Sprintf("Something gone wrong %v", e)
					return fmt.Errorf(msg)
				}
			} else {
				msg := fmt.Sprintf("Getting Network %s failed: %s", cidr, err)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
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
	d.Set("properties", network.Properties)
	log.Debugf("Completed getting Network %s", d.Get("cidr"))
	return nil
}

// updateNetwork Update the existing IPv4/IPv6 Network
func updateNetwork(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	network := entities.Network{}
	if !network.InitNetwork(d) {
		log.Error(network.InitError)
		return fmt.Errorf(network.InitError)
	}
	log.Debugf("Beginning to update Network %s", network.CIDR)

	_, err := objMgr.UpdateNetwork(network)
	if err != nil {
		msg := fmt.Sprintf("Error updating Network (%s): %s", network.CIDR, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	log.Debugf("Completed to update Network %s", network.CIDR)
	return getNetwork(d, m)
}

// deleteNetwork Delete the IPv4/IPv6 Network
func deleteNetwork(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	network := entities.Network{}
	if !network.InitNetwork(d) {
		log.Error(network.InitError)
		return fmt.Errorf(network.InitError)
	}
	log.Debugf("Beginning to delete Network %s", network.CIDR)

	_, err := objMgr.DeleteNetwork(network)
	if err != nil {
		msg := fmt.Sprintf("Delete Network %s failed: %s", network.CIDR, err)
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

	return utils.GetPropertyValue(fieldName, object.Properties)
}
