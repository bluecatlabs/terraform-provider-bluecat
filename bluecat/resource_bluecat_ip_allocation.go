// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceIPAllocation The IP Allocation
func ResourceIPAllocation() *schema.Resource {
	return &schema.Resource{
		Create: createIPAllocation,
		Read:   getIPAllocation,
		Update: updateIPAllocation,
		Delete: deleteIPAllocation,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Allocating the IP address/Host record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, record will be created under default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone in which you want to update a host record. If not provided, the absolute name must be FQDN ones",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Host record. Must be FQDN if the Zone is not provided",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					zone := d.Get("zone").(string)
					return checkDiffName(old, new, zone)
				},
			},
			"network": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Network address in CIDR format",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IPv4/IPv6 address. If no value is given, a next available IP address in the network will be allocated",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new != "" {
						if new == old {
							return true
						}
						return false
					}
					return true
				},
			},
			"mac_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The MAC address",
			},
			"properties": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP address/Host record's properties. Example: attribute=value|",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return checkDiffProperties(old, new)
				},
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Desired IP4 address state: MAKE_STATIC / MAKE_RESERVED / MAKE_DHCP_RESERVED",
				Default:     entities.AllocateStatic,
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 Template which you want to assign",
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP Address version: ipv4 or ipv6",
			},
		},
	}
}

// createIPAllocation Allocate the IPv4/IPv6 address/Host record
// Create the host record if the zone name is provided
// In case of allocating the IP address, the network must be specified
func createIPAllocation(d *schema.ResourceData, m interface{}) error {
	objMgr := GetObjManager(m)

	address := entities.IPAddress{}
	if !address.InitIPAddress(d) {
		log.Error(address.InitError)
		return fmt.Errorf(address.InitError)
	}

	// these props are not directly related to address
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	network := d.Get("network").(string)
	log.Debugf("Beginning to allocate IP address in the network %s", network)

	fqdnName := address.Name
	if len(zone) > 0 {
		fqdnName = getFQDN(address.Name, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}

	createIP := true
	if len(address.Address) != 0 {
		_, err := objMgr.GetIPAddress(address.Configuration, address.Address, address.IPVersion)
		if err != nil {
			log.Debugf("The linked IP address doesn't exist")
		} else {
			createIP = false
			err = updateAllocatedResource(d, m)
			if err != nil {
				return err
			}
		}
	} else {
		address.CIDR = strings.Split(network, "/")[0]
	}
	if createIP {
		log.Debugf("Allocating the IP address under network %s", network)
		ipProperties := address.Properties
		if len(zone) > 0 {
			ipProperties = ""
		}
		if address.Action == entities.AllocateStatic && address.IPVersion == entities.IPV4 {
			ipProperties = fmt.Sprintf("%s|excludeDHCPRange=true", ipProperties)
		} else if address.Action == entities.AllocateReserved {
			address.Name = strings.Split(fqdnName, fmt.Sprintf(".%s", zone))[0][0:]
		}
		newIPAddress, err := objMgr.CreateIPAddress(address)
		if err != nil {
			msg := fmt.Sprintf("Error allocating IP from network %s: %s", network, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		if len(address.Address) == 0 {
			//No IP address, so need to get the IP after got the new ones in the above step
			address.Address = newIPAddress.Address
			log.Debugf("Got the IP address %s", address.Address)
		}
	}

	if len(address.Mac) > 0 {
		log.Debugf("Updating the MAC address for the IP address %s", address.Address)
		address.Action = ""
		_, err := objMgr.SetMACAddress(address)
		if err != nil {
			msg := fmt.Sprintf("Updating IP address %s failed: %s", address.Address, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	}

	d.Set("ip_address", address.Address)
	d.Set("name", fqdnName)

	if len(zone) > 0 {
		if address.Action != entities.AllocateReserved {
			log.Debugf("Creating the Host record %s", fqdnName)
			_, err := objMgr.CreateHostRecord(address.Configuration, view, zone, fqdnName, address.Address, -1, address.Properties)
			if err != nil {
				msg := fmt.Sprintf("Error creating the Host record %s: %s", fqdnName, err)
				log.Debug(msg)
				return fmt.Errorf(msg)
			}
		}

		d.Set("name", fqdnName)
		log.Debugf("Finished to create the Host record %s", fqdnName)
	}

	log.Debugf("Completed to allocate IP address %s", address.Address)
	return getIPAllocation(d, m)
}

// getIPAllocation Get the allocated IP address/Host info
func getIPAllocation(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	address := entities.IPAddress{}
	if !address.InitIPAddress(d) {
		log.Error(address.InitError)
		return fmt.Errorf(address.InitError)
	}
	log.Debugf("Beginning to get IP address: %s", address.Address)

	zone := d.Get("zone").(string)

	fqdnName := address.Name
	if len(zone) > 0 {
		fqdnName = getFQDN(address.Name, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}

	properties := ""
	if len(zone) > 0 {
		view := d.Get("view").(string)
		log.Debugf("Getting Host record info %s", fqdnName)
		hostRecord, err := objMgr.GetHostRecord(address.Configuration, view, fqdnName)
		if err == nil {
			properties = hostRecord.Properties
		}
	}
	if properties == "" {
		log.Debugf("Getting IP address info %s", address.Address)
		ipAddress, err := objMgr.GetIPAddress(address.Configuration, address.Address, address.IPVersion)
		if err != nil {
			msg := fmt.Sprintf("Getting IP address %s failed: %s", address.Address, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		properties = ipAddress.Properties
	}
	d.Set("name", fqdnName)
	d.Set("properties", properties)
	d.SetId(fqdnName)
	log.Debugf("Completed reading IP address %s", address.Address)
	return nil
}

// updateIPAllocation Update the allocation
func updateIPAllocation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update the allocation for the IP address %s", d.Get("ip_address"))
	err := updateAllocatedResource(d, m)
	if err != nil {
		return err
	}
	log.Debugf("Completed to update allocation %s", d.Get("ip4_address"))
	return getIPAllocation(d, m)
}

// deleteIPAllocation Delete the allocated IP address/Host record
func deleteIPAllocation(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	address := entities.IPAddress{}
	if !address.InitIPAddress(d) {
		log.Error(address.InitError)
		return fmt.Errorf(address.InitError)
	}
	log.Debugf("Beginning to release an IP allocated in the network %s", d.Get("network"))

	log.Debugf("Checking the IP address %s for deletion", address.Address)
	_, err := objMgr.GetIPAddress(address.Configuration, address.Address, address.IPVersion)
	if err != nil {
		msg := fmt.Sprintf("The IP address %s not found: %s", address.Address, err)
		log.Debug(msg)
	} else {
		log.Debugf("Deleting the IP address %s", address.Address)
		_, err := objMgr.DeleteIPAddress(address.Configuration, address.Address, address.IPVersion)
		if err != nil {
			msg := fmt.Sprintf("Delete IP address %s failed: %s", address.Address, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	}

	d.SetId("")
	log.Debugf("Completed to release an IP allocated in the network %s", d.Get("network"))
	return nil
}

// updateAllocatedResource Update the allocated IP address/Host record
func updateAllocatedResource(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	address := entities.IPAddress{}
	if !address.InitIPAddress(d) {
		log.Error(address.InitError)
		return fmt.Errorf(address.InitError)
	}

	log.Debugf("Updating allocated resource in network %s", d.Get("network"))

	view := d.Get("view").(string)
	zone := d.Get("zone").(string)

	fqdnName := address.Name
	if len(zone) > 0 {
		fqdnName = getFQDN(address.Name, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}

	if len(zone) > 0 {
		log.Debugf("Updating host record %s", fqdnName)
		hostRecord, err := objMgr.GetHostRecord(address.Configuration, view, fqdnName)
		if err == nil {
			// Keeps values as in the server
			log.Debugf(hostRecord.Properties)
			TTL := getPropertyValue("ttl", hostRecord.Properties)
			rrTTL, err := strconv.Atoi(TTL)
			if err != nil {
				msg := fmt.Sprintf("Convert Host record TTL %s failed: %s", TTL, err)
				log.Debug(msg)
				rrTTL = -1
			}

			associateIPs := address.Address
			currentAssociateIPs := getPropertyValue("addresses", hostRecord.Properties)
			if len(currentAssociateIPs) > 0 {
				associateIPs = fmt.Sprintf("%s,%s", currentAssociateIPs, address.Address)
			}

			var immutableProperties = []string{"parentId", "parentType"} // these properties will raise error on the rest-api
			address.Properties = utils.RemoveImmutableProperties(address.Properties, immutableProperties)

			_, err = objMgr.UpdateHostRecord(address.Configuration, view, zone, fqdnName, associateIPs, rrTTL, address.Properties)
			if err != nil {
				msg := fmt.Sprintf("Error updating Host record %s: %s", fqdnName, err)
				log.Debug(msg)
				return fmt.Errorf(msg)
			}
		} else {
			msg := fmt.Sprintf("Getting Host record %s failed: %s", address.Name, err)
			log.Debug(msg)
		}
	}
	log.Debugf("Updating IP address %s", address.Address)
	ipAddress, err := objMgr.GetIPAddress(address.Configuration, address.Address, address.IPVersion)
	if err != nil {
		msg := fmt.Sprintf("Getting IP address %s failed: %s", address.Address, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	// The properties field belongs to Host record if the zone field is not none
	if len(zone) > 0 {
		address.Properties = ""
	} else {
		// 'address' and 'state' can not be changed.
		// macAddress is updated by the mac_address field
		address.Properties = removeAttributeFromProperties("address", address.Properties)
		address.Properties = removeAttributeFromProperties("state", address.Properties)
		address.Properties = removeAttributeFromProperties("macAddress", address.Properties)
	}

	ipAddress.Action = getAttributeFromProperties("state", ipAddress.Properties)
	ipAddress.SetAction()

	// if old ip address state is the same as the new then do not try to change to the same state
	// because that will raise a bug on the REST-API side
	if address.Action == ipAddress.Action {
		address.Action = ""
	}
	if address.Action == "RESERVED" || address.Action == entities.AllocateReserved {
		address.Name = strings.Split(fqdnName, fmt.Sprintf(".%s", zone))[0][0:]
	}
	_, err = objMgr.UpdateIPAddress(address)
	if err != nil {
		msg := fmt.Sprintf("Error updating IP address %s: %s", address.Address, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	log.Debugf("Completed to update the allocated resource in network %s", d.Get("network"))
	return nil
}

// getPropertyValue Get the property value by key from the properties string
func getPropertyValue(key, props string) (val string) {
	properties := strings.Split(props, "|")
	for i := 0; i < len(properties); i++ {
		prop := strings.Split(properties[i], "=")
		if prop[0] == key {
			val = prop[1]
			return
		}
	}
	return
}
