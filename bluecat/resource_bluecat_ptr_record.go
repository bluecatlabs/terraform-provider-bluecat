// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"ip_address": {
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
			"reverse_record": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "To create a reverse record for the pass host",
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
	ipAddress := d.Get("ip_address").(string)
	ttl := d.Get("ttl").(int)
	reverseRecord := d.Get("reverse_record").(string)
	properties := d.Get("properties").(string)
	fqdnName, err := updatePTR(m, configuration, view, zone, name, ipAddress, reverseRecord, properties, ttl)
	if err != nil {
		return err
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
	ipAddress := d.Get("ip_address").(string)
	ttl := d.Get("ttl").(int)
	reverseRecord := d.Get("reverse_record").(string)
	properties := d.Get("properties").(string)

	fqdnName, err := updatePTR(m, configuration, view, zone, name, ipAddress, reverseRecord, properties, ttl)
	if err != nil {
		return err
	}
	d.Set("name", fqdnName)
	log.Debugf("Completed to update PTR record %s", d.Get("name"))
	return getPTRRecord(d, m)
}

// updatePTRRecord Update the existing PTR record
// Update the PTR, just set the reverseRecord flag
func updatePTR(m interface{}, configuration, view, zone, name, ip4Address, reverseRecord, properties string, ttl int) (string, error) {
	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	fqdnName := name

	if len(zone) > 0 {
		fqdnName = getFQDN(fqdnName, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}

	// Get the host
	_, err := objMgr.GetHostRecord(configuration, view, fqdnName)
	if err != nil {
		msg := fmt.Sprintf("Getting Host record %s failed: %s", fqdnName, err)
		log.Debug(msg)
		return "", fmt.Errorf(msg)
	}

	// Update the host record
	reverseValues := []string{"yes", "true", "1"}
	notReversedValues := []string{"no", "false", "0", ""}
	isReverse := contains(reverseValues, strings.ToLower(strings.Trim(reverseRecord, " ")))
	isNotReverse := contains(notReversedValues, strings.ToLower(strings.Trim(reverseRecord, " ")))

	properties = removeAttributeFromProperties("reverseRecord", properties)
	if isReverse || isNotReverse {
		properties = fmt.Sprintf("%s|reverseRecord=%t", properties, isReverse)
	} else {
		msg := fmt.Sprintf("invalid reverse_record value (must be either 'true' or 'false'): '%s'", reverseRecord)
		log.Debug(msg)
		return "", fmt.Errorf(msg)
	}

	var immutableProperties = []string{"parentId", "parentType"} // these properties will raise error on the rest-api
	properties = utils.RemoveImmutableProperties(properties, immutableProperties)

	_, err = objMgr.UpdateHostRecord(configuration, view, zone, fqdnName, ip4Address, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating PTR record %s: %s", fqdnName, err)
		log.Debug(msg)
		return "", fmt.Errorf(msg)
	}
	return fqdnName, nil
}

// deletePTRRecord Delete the PTR record
// To delete the PTR, just set the reverseRecord flag to False
func deletePTRRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete PTR record %s", d.Get("name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	ipAddress := d.Get("ip_address").(string)
	ttl := d.Get("ttl").(int)
	reverseRecord := "false"
	properties := d.Get("properties").(string)

	_, err := updatePTR(m, configuration, view, zone, name, ipAddress, reverseRecord, properties, ttl)
	if err != nil {
		return err
	}
	d.SetId("")
	log.Debugf("Completed to delete PTR record %s", d.Get("name"))
	return nil
}
