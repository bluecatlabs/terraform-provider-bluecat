// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Description: "The Zone in which you want to update a host record",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the IP address/Host record",
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
			"ip4_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address. If no value is given, a next available IP address in the network will be allocated",
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
		},
	}
}

// createIPAllocation Allocate the IP address/Host record
// Create the host record if the zone name is provided
// In case of allocating the IP address, the network must be specified
func createIPAllocation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to allocate IP address in the network %s", d.Get("network"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	network := d.Get("network").(string)
	ip4Address := d.Get("ip4_address").(string)
	macAddress := d.Get("mac_address").(string)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	addrCIDR := strings.Split(network, "/")
	fqdnName := getFQDN(name, zone)

	if len(ip4Address) == 0 || len(zone) == 0 {
		createIP := true
		if len(ip4Address) != 0 {
			_, err := objMgr.GetIPAddress(configuration, ip4Address)
			if err != nil {
				log.Debugf("The linked IP address doesn't exist")
			} else {
				createIP = false
				err = updateAllocatedResource(d, m)
				if err != nil {
					return err
				}
			}
		}
		if createIP {
			log.Debugf("Allocating the IP address under network %s", network)
			ipProperties := properties
			if len(zone)>0 {
				ipProperties = ""
			}
			newIPAddress, err := objMgr.CreateStaticIP(configuration, addrCIDR[0], ip4Address, macAddress, fqdnName, ipProperties)
			if err != nil {
				msg := fmt.Sprintf("Error allocating IP from network %s: %s", network, err)
				log.Debug(msg)
				return fmt.Errorf(msg)
			}
			if len(ip4Address) == 0 {
				//No IP address, so need to get the IP after got the new ones in the above step
				ip4Address = getPropertyValue("address", newIPAddress.Properties)
				log.Debugf("Got the IP address %s", ip4Address)
			}
		}
		d.Set("ip4_address", ip4Address)
		d.Set("name", fqdnName)
	}

	if len(zone) > 0 {
		log.Debugf("Creating the Host record %s", fqdnName)
		_, err := objMgr.CreateHostRecord(configuration, view, zone, fqdnName, ip4Address, -1, properties)
		if err != nil {
			msg := fmt.Sprintf("Error creating the Host record %s: %s", fqdnName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		if len(macAddress) > 0 {
			log.Debugf("Updating the MAC address for the IP address %s", ip4Address)
			_, err = objMgr.SetMACAddress(configuration, ip4Address, macAddress)
			if err != nil {
				msg := fmt.Sprintf("Updating IP address %s failed: %s", ip4Address, err)
				log.Debug(msg)
				return fmt.Errorf(msg)
			}
		}
		d.Set("name", fqdnName)
		log.Debugf("Finished to create the Host record %s", fqdnName)
	}

	log.Debugf("Completed to allocate IP address %s", ip4Address)
	return getIPAllocation(d, m)
}

// getIPAllocation Get the allocated IP address/Host info
func getIPAllocation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get IP address: %s", d.Get("ip4_address").(string))
	configuration := d.Get("configuration").(string)
	zone := d.Get("zone").(string)
	ip4Address := d.Get("ip4_address").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	fqdnName := getFQDN(d.Get("name").(string), zone)

	properties := ""
	if len(zone) > 0 {
		view := d.Get("view").(string)
		log.Debugf("Getting Host record info %s", fqdnName)
		hostRecord, err := objMgr.GetHostRecord(configuration, view, fqdnName)
		if err != nil {
			msg := fmt.Sprintf("Getting Host record %s failed: %s", fqdnName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		properties = hostRecord.Properties
	} else {
		log.Debugf("Getting IP address info %s", ip4Address)
		ipAddress, err := objMgr.GetIPAddress(configuration, ip4Address)
		if err != nil {
			msg := fmt.Sprintf("Getting IP address %s failed: %s", ip4Address, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		properties = ipAddress.Properties
	}
	d.Set("name", fqdnName)
	d.Set("properties", properties)
	d.SetId(fqdnName)
	log.Debugf("Completed reading IP address %s", ip4Address)
	return nil
}

// updateIPAllocation Update the allocation
func updateIPAllocation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update the allocation for the IP address %s", d.Get("ip4_address"))
	err := updateAllocatedResource(d, m)
	if err != nil {
		return err
	}
	log.Debugf("Completed to update allocation %s", d.Get("ip4_address"))
	return getIPAllocation(d, m)
}

// deleteIPAllocation Delete the allocated IP address/Host record
func deleteIPAllocation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to release an IP allocated in the network %s", d.Get("network"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	ip4Address := d.Get("ip4_address").(string)
	fqdnName := getFQDN(d.Get("name").(string), zone)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	if len(zone) > 0 {
		log.Debugf("Deleting the host record %s", fqdnName)
		_, err := objMgr.DeleteHostRecord(configuration, view, fqdnName)
		if err != nil {
			msg := fmt.Sprintf("Delete Host record %s failed: %s", fqdnName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	} else {
		log.Debugf("Deleting the IP address %s", ip4Address)
		_, err := objMgr.DeleteIPAddress(configuration, ip4Address)
		if err != nil {
			msg := fmt.Sprintf("Delete IP address %s failed: %s", ip4Address, err)
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
	log.Debugf("Updating allocated resource in network %s", d.Get("network"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	ip4Address := d.Get("ip4_address").(string)
	macAddress := d.Get("mac_address").(string)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	fqdnName := getFQDN(name, zone)

	if len(zone) > 0 {
		log.Debugf("Updating host record %s", fqdnName)
		hostRecord, err := objMgr.GetHostRecord(configuration, view, fqdnName)
		if err != nil {
			msg := fmt.Sprintf("Getting Host record %s failed: %s", name, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}

		// Keeps values as in the server
		log.Debugf(hostRecord.Properties)
		TTL := getPropertyValue("ttl", hostRecord.Properties)
		rrTTL, err := strconv.Atoi(TTL)
		if err != nil {
			msg := fmt.Sprintf("Convert Host record TTL %s failed: %s", TTL, err)
			log.Debug(msg)
			rrTTL = -1
		}

		associtaeIPs := ip4Address
		if len(associtaeIPs) == 0 {
			associtaeIPs := getPropertyValue("addresses", hostRecord.Properties)
			ip4Address = strings.Split(associtaeIPs, ",")[0]
		}
		_, err = objMgr.UpdateHostRecord(configuration, view, zone, fqdnName, associtaeIPs, rrTTL, properties)
		if err != nil {
			msg := fmt.Sprintf("Error updating Host record %s: %s", fqdnName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	}
	log.Debugf("Updating IP address %s", ip4Address)
	ipAddress, err := objMgr.GetIPAddress(configuration, ip4Address)
	if err != nil {
		msg := fmt.Sprintf("Getting IP address %s failed: %s", ip4Address, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	// The properties field belongs to Host record if the zone field is not none
	if len(zone)>0 {
		properties = ""
	} else {
		// 'address' and 'state' can not be changed.
		// macAddress is updated by the mac_address field
		properties = removeAttributeFromProperties("address", properties)
		properties = removeAttributeFromProperties("state", properties)
		properties = removeAttributeFromProperties("macAddress", properties)
	}
	_, err = objMgr.UpdateIPAddress(configuration, ip4Address, macAddress, fqdnName, ipAddress.Action, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating IP address %s: %s", ip4Address, err)
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
