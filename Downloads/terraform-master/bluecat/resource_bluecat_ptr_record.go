// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourcePTRRecord The PTR record
func ResourcePTRRecord() *schema.Resource {
	return &schema.Resource{
		Create: createPTRRecord,
		Read:   getPTRRecord,
		Update: updatePTRRecord,
		Delete: deletePTRRecord,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the PTR record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, record will be created under default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Zone in which you want to update a host record",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the host record",
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
				Required:    true,
				Description: "The IP address that will be created the PTR record for",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The TTL value",
				Default:     -1,
			},
			"properties": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Record's properties. Example: attribute=value|",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return checkDiffProperties(old, new)
				  },
			},
		},
	}
}

// createPTRRecord Create the new PTR record
// Create the Host record, then server will create the PTR
func createPTRRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create PTR record %s", d.Get("name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	network := d.Get("network").(string)
	ip4Address := d.Get("ip4_address").(string)
	ttl := d.Get("ttl").(int)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	fqdnName := name

	if len(zone) > 0 {
		fqdnName = getFQDN(fqdnName, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}
	// Check if the network is already exist?
	_, err := objMgr.GetNetwork(configuration, network)
	if err != nil {
		msg := fmt.Sprintf("Getting Network %s failed: %s", network, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	properties = fmt.Sprintf("%s|reverseRecord=true", properties)
	_, err = objMgr.CreateHostRecord(configuration, view, zone, fqdnName, ip4Address, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating PTR record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.Set("name", fqdnName)
	log.Debugf("Completed to create PTR record %s", d.Get("name"))
	return getPTRRecord(d, m)
}

// getPTRRecord Get the PTR record
func getPTRRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get PTR record: %s", d.Get("name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	name := d.Get("name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	hostRecord, err := objMgr.GetHostRecord(configuration, view, name)
	if err != nil {
		msg := fmt.Sprintf("Getting PTR record %s failed: %s", name, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(hostRecord.AbsoluteName)
	d.Set("name", hostRecord.AbsoluteName)
	d.Set("properties", hostRecord.Properties)
	log.Debugf("Completed reading PTR record %s", d.Get("name"))
	return nil
}

// updatePTRRecord Update the existing PTR record
func updatePTRRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update PTR record %s", d.Get("name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	network := d.Get("network").(string)
	ip4Address := d.Get("ip4_address").(string)
	ttl := d.Get("ttl").(int)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	fqdnName := name

	if len(zone) > 0 {
		fqdnName = getFQDN(fqdnName, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}
	// Check network
	_, err := objMgr.GetNetwork(configuration, network)
	if err != nil {
		msg := fmt.Sprintf("Getting Network %s failed: %s", network, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	// Check IP address, create the new one if doesn't exist
	_, err = objMgr.GetIPAddress(configuration, ip4Address)
	if err != nil {
		log.Debugf("The linked IP address doesn't exist, allocating the IP address %s", ip4Address)
		_, err = objMgr.CreateStaticIP(configuration, strings.Split(network, "/")[0], ip4Address, "", fqdnName, "")
		if err != nil {
			msg := fmt.Sprintf("Error allocating IP addess %s from network %s: %s", ip4Address, network, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	}
	// Update the host record
	_, err = objMgr.UpdateHostRecord(configuration, view, zone, fqdnName, ip4Address, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating PTR record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.Set("name", fqdnName)
	log.Debugf("Completed to update PTR record %s", d.Get("name"))
	return getPTRRecord(d, m)
}

// deletePTRRecord Delete the PTR record
// To delete the PTR, just delete an IP address from the server
func deletePTRRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete PTR record %s", d.Get("name"))
	configuration := d.Get("configuration").(string)
	name := d.Get("name").(string)
	ip4Address := d.Get("ip4_address").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteIPAddress(configuration, ip4Address)
	if err != nil {
		msg := fmt.Sprintf("Getting PTR record %s failed: %s", name, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Completed to delete PTR record %s", d.Get("name"))
	return nil
}
