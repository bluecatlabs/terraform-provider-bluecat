// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceIPAssociation The IP Association
func ResourceIPAssociation() *schema.Resource {
	return &schema.Resource{
		Create: createIPAssociation,
		Read:   getIPAssociation,
		Update: updateIPAssociation,
		Delete: deleteIPAssociation,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Associate the IP address/Host record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, uses the default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone in which you want to update a host record",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Host record",
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
				Description: "The IP address",
			},
			"mac_address": {
				Type:        schema.TypeString,
				Required:    true,
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

// createIPAssociation Associate the IP address/Host record
func createIPAssociation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to associate IP address %s", d.Get("ip4_address"))
	err := updateAllocatedResource(d, m)
	if err != nil {
		return err
	}
	log.Debugf("Completed to associate IP address %s", d.Get("ip4_address"))
	return getIPAssociation(d, m)
}

// getIPAssociation Get the allocated IP address/Host info
func getIPAssociation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get IP address: %s", d.Get("ip4_address").(string))
	err := getIPAllocation(d, m)
	if err != nil {
		return err
	}
	log.Debugf("Completed reading IP address %s", d.Get("ip4_address"))
	return nil
}

// updateIPAssociation Update the association
func updateIPAssociation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update the association for the IP address %s", d.Get("ip4_address"))
	err := updateAllocatedResource(d, m)
	if err != nil {
		return err
	}
	log.Debugf("Completed to update association %s", d.Get("ip4_address"))
	return getIPAssociation(d, m)
}

// deleteIPAssociation Delete the association IP address/Host record
func deleteIPAssociation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to release an association for the IP address %s", d.Get("ip4_address"))
	configuration := d.Get("configuration").(string)
	ip4Address := d.Get("ip4_address").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.SetMACAddress(configuration, ip4Address, "00:00:00:00:00:00")
	if err != nil {
		msg := fmt.Sprintf("Releasing the IP address %s failed: %s", ip4Address, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Completed to release an association for the IP address %s", ip4Address)
	return nil
}
