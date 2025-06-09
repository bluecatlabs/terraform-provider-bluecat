// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceSRVRecord The SRV record
func ResourceSRVRecord() *schema.Resource {
	return &schema.Resource{
		Create: createSRVRecord,
		Read:   getSRVRecord,
		Update: updateSRVRecord,
		Delete: deleteSRVRecord,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the SRV record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, record will be created under default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone in which you want to update a SRV record. If not provided, the absolute name must be FQDN ones",
			},
			"absolute_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the SRV record. Must be FQDN if the Zone is not provided",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					zone := d.Get("zone").(string)
					return checkDiffName(old, new, zone)
				},
			},
			"linked_record": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The host record that the SRV record links to",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					zone := d.Get("zone").(string)
					return checkDiffName(old, new, zone)
				},
			},
			"weight": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "This is the weight, used to determine which server to connect to if multiple servers have the same priority",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "This is the port number on which the service is listening",
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The priority of the record, a lower value is a higher priority",
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
				Description: "SRV record's properties. Example: attribute=value|",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return checkDiffProperties(old, new)
				},
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SRV Records name, used exclusively for changing the name of the record. For identification, use absolute_name.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: recordImporter,
		},
	}
}

// createSRVRecord Create the new SRV record
func createSRVRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create SRV record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	absoluteName := d.Get("absolute_name").(string)
	linkedRecord := d.Get("linked_record").(string)
	weight := d.Get("weight").(int)
	port := d.Get("port").(int)
	priority := d.Get("priority").(int)
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

	_, err := objMgr.CreateSRVRecord(configuration, view, zone, priority, port, weight, fqdnName, linkedRecord, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating SRV record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.Set("absolute_name", fqdnName)
	log.Debugf("Completed to create SRV record %s", d.Get("absolute_name"))
	return getSRVRecord(d, m)
}

// getSRVRecord Get the SRV record
func getSRVRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get SRV record: %s", d.Get("absolute_name"))
	absoluteName, err := getAbsoluteName(d)
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	srvRecord, err := objMgr.GetSRVRecord(configuration, view, absoluteName)
	if err != nil {
		if d.Id() != "" {
			err := createSRVRecord(d, m)
			if err != nil {
				msg := fmt.Sprintf("Something gone wrong: %v", err)
				return fmt.Errorf(msg)
			}
		} else {
			msg := fmt.Sprintf("Getting SRV record %s failed: %s", absoluteName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	}
	d.SetId(srvRecord.AbsoluteName)
	d.Set("absolute_name", srvRecord.AbsoluteName)
	d.Set("properties", srvRecord.Properties)

	log.Debugf("Completed reading SRV record %s", d.Get("absolute_name"))
	return nil
}

// updateSRVRecord Update the existing SRV record
func updateSRVRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update SRV record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	absoluteName := d.Get("absolute_name").(string)
	linkedRecord := d.Get("linked_record").(string)
	weight := d.Get("weight").(int)
	port := d.Get("port").(int)
	priority := d.Get("priority").(int)
	ttl := d.Get("ttl").(int)
	properties := d.Get("properties").(string)
	name := d.Get("name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	fqdnName := absoluteName

	if len(zone) > 0 {
		fqdnName = getFQDN(absoluteName, zone)
	} else {
		zone = getZoneFromRRName(fqdnName)
	}

	var immutableProperties = []string{"parentId", "parentType"} // these properties will raise error on the rest-api
	properties = utils.RemoveImmutableProperties(properties, immutableProperties)

	_, err := objMgr.UpdateSRVRecord(configuration, view, zone, priority, port, weight, fqdnName, linkedRecord, ttl, properties, name)
	if err != nil {
		msg := fmt.Sprintf("Error updating SRV record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	if name != "" {
		fqdnName = replaceName(fqdnName, name)
	}
	d.Set("absolute_name", fqdnName)
	d.SetId(fqdnName)
	log.Debugf("Completed to update SRV record %s", d.Get("absolute_name"))
	return nil
}

// deleteSRVRecord Delete the SRV record
func deleteSRVRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete SRV record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	absoluteName := d.Get("absolute_name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteSRVRecord(configuration, view, absoluteName)
	if err != nil {
		msg := fmt.Sprintf("Getting SRV record %s failed: %s", absoluteName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Completed to delete SRV record %s", d.Get("absolute_name"))
	return nil
}

func replaceName(url, newName string) string {
	// Find the index of the first dot
	dotIndex := strings.Index(url, ".")
	if dotIndex == -1 {
		// If there's no dot, return the original URL
		return url
	}
	// Replace the part before the first dot
	return newName + url[dotIndex:]
}
