// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceGenericRecord The Generic record
func ResourceGenericRecord() *schema.Resource {
	return &schema.Resource{
		Create: createGenericRecord,
		Read:   getGenericRecord,
		Update: updateGenericRecord,
		Delete: deleteGenericRecord,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the Generic record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, record will be created under default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone in which you want to update a Generic record. If not provided, the absolute name must be FQDN ones",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Type in which you want to create type of Generic record. If not provided, record will be created fail",
			},
			"absolute_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Generic record. Must be FQDN if the Zone is not provided",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					zone := d.Get("zone").(string)
					return checkDiffName(old, new, zone)
				},
			},
			"data": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Data of the Generic record",
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
				Description: "Generic record's properties. Example: attribute=value|",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return checkDiffProperties(old, new)
				},
			},
			"to_deploy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether or not to selectively deploy the Generic record",
				Default:     "no",
			},
		},
		Importer: &schema.ResourceImporter{
			State: recordImporter,
		},
	}
}

// createGenericRecord Create the new Generic record
func createGenericRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create Generic record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	typerr := d.Get("type").(string)
	absoluteName := d.Get("absolute_name").(string)
	data := d.Get("data").(string)
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

	genericRecord, err := objMgr.CreateGenericRecord(configuration, view, zone, typerr, fqdnName, data, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating Generic record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	deploy := utils.ParseDeploymentValue(d.Get("to_deploy").(string))
	if deploy {
		res, err := objMgr.Connector.DeployObject(genericRecord)
		if err != nil {
			msg := fmt.Sprintf("Error deploying Generic record %s: %s", absoluteName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		log.Debugf("Successfully deployed. %s", res)
	}
	d.Set("absolute_name", fqdnName)
	log.Debugf("Completed to create Generic record %s", d.Get("absolute_name"))
	return getGenericRecord(d, m)
}

// getGenericRecord Get the Generic record
func getGenericRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get Generic record: %s", d.Get("absolute_name"))
	absoluteName, err := getAbsoluteName(d)
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	genericRecord, err := objMgr.GetGenericRecord(configuration, view, absoluteName)
	if err != nil {
		if d.Id() != "" {
			err := createGenericRecord(d, m)
			if err != nil {
				msg := fmt.Sprintf("Something gone wrong: %v", err)
				return fmt.Errorf(msg)
			}
		} else {
			msg := fmt.Sprintf("Getting Generic record %s failed: %s", absoluteName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	}
	d.SetId(genericRecord.AbsoluteName)
	d.Set("absolute_name", genericRecord.AbsoluteName)
	d.Set("properties", genericRecord.Properties)
	log.Debugf("Completed reading Generic record %s", d.Get("absolute_name"))
	return nil
}

// updateGenericRecord Update the existing Generic record
func updateGenericRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update Generic record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	typerr := d.Get("type").(string)
	absoluteName := d.Get("absolute_name").(string)
	data := d.Get("data").(string)
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

	var immutableProperties = []string{"parentId", "parentType"} // these properties will raise error on the rest-api
	properties = utils.RemoveImmutableProperties(properties, immutableProperties)

	genericRecord, err := objMgr.UpdateGenericRecord(configuration, view, zone, typerr, fqdnName, data, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating Generic record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	deploy := utils.ParseDeploymentValue(d.Get("to_deploy").(string))
	if deploy {
		res, err := objMgr.Connector.DeployObject(genericRecord)
		if err != nil {
			msg := fmt.Sprintf("Error deploying Generic record %s: %s", absoluteName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		log.Debugf("Successfully deployed. %s", res)
	}
	d.Set("absolute_name", fqdnName)
	log.Debugf("Completed to update Generic record %s", d.Get("absolute_name"))
	return getGenericRecord(d, m)
}

// deleteGenericRecord Delete the Generic record
func deleteGenericRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete Generic record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	absoluteName := d.Get("absolute_name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteGenericRecord(configuration, view, absoluteName)
	if err != nil {
		msg := fmt.Sprintf("Getting Generic record %s failed: %s", absoluteName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Completed to delete Generic record %s", d.Get("absolute_name"))
	return nil
}
