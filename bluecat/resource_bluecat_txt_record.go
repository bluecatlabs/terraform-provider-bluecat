// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceTXTRecord The TXT record
func ResourceTXTRecord() *schema.Resource {
	return &schema.Resource{
		Create: createTXTRecord,
		Read:   getTXTRecord,
		Update: updateTXTRecord,
		Delete: deleteTXTRecord,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the TXT record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, record will be created under default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone in which you want to update a TXT record. If not provided, the absolute name must be FQDN ones",
			},
			"absolute_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the TXT record. Must be FQDN if the Zone is not provided",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					zone := d.Get("zone").(string)
					return checkDiffName(old, new, zone)
				},
			},
			"text": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Set the text of TXT record",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					zone := d.Get("zone").(string)
					return checkDiffName(old, new, zone)
				},
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
				Description: "TXT record's properties. Example: attribute=value|",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return checkDiffProperties(old, new)
				},
			},
		},
	}
}

// createTXTRecord Create the new TXT record
func createTXTRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create TXT record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	absoluteName := d.Get("absolute_name").(string)
	text := d.Get("text").(string)
	ttl := d.Get("ttl").(int)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	fqdnName := absoluteName

	if len(zone) > 0 {
		fqdnName = getFQDN(absoluteName, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}

	_, err := objMgr.CreateTXTRecord(configuration, view, zone, fqdnName, text, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating TXT record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.Set("absolute_name", fqdnName)
	log.Debugf("Completed to create TXT record %s", d.Get("absolute_name"))
	return getTXTRecord(d, m)
}

// getTXTRecord Get the TXT record
func getTXTRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get TXT record: %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	absoluteName := d.Get("absolute_name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	txtRecord, err := objMgr.GetTXTRecord(configuration, view, absoluteName)
	if err != nil {
		msg := fmt.Sprintf("Getting TXT record %s failed: %s", absoluteName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(txtRecord.AbsoluteName)
	d.Set("absolute_name", txtRecord.AbsoluteName)
	d.Set("properties", txtRecord.Properties)
	log.Debugf("Completed reading TXT record %s", d.Get("absolute_name"))
	return nil
}

// updateTXTRecord Update the existing TXT record
func updateTXTRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update TXT record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	absoluteName := d.Get("absolute_name").(string)
	text := d.Get("text").(string)
	ttl := d.Get("ttl").(int)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	fqdnName := absoluteName

	if len(zone) > 0 {
		fqdnName = getFQDN(absoluteName, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}

	_, err := objMgr.UpdateTXTRecord(configuration, view, zone, fqdnName, text, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating TXT record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.Set("absolute_name", fqdnName)
	log.Debugf("Completed to update TXT record %s", d.Get("absolute_name"))
	return getTXTRecord(d, m)
}

// deleteTXTRecord Delete the TXT record
func deleteTXTRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete TXT record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	absoluteName := d.Get("absolute_name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteTXTRecord(configuration, view, absoluteName)
	if err != nil {
		msg := fmt.Sprintf("Getting TXT record %s failed: %s", absoluteName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Completed to delete TXT record %s", d.Get("absolute_name"))
	return nil
}
