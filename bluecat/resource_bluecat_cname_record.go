// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceCNAMERecord The CNAME record
func ResourceCNAMERecord() *schema.Resource {
	return &schema.Resource{
		Create: createCNAMERecord,
		Read:   getCNAMERecord,
		Update: updateCNAMERecord,
		Delete: deleteCNAMERecord,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the CNAME record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, record will be created under default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone in which you want to update a CNAME record. If not provided, the absolute name must be FQDN ones",
			},
			"absolute_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the CNAME record. Must be FQDN if the Zone is not provided",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					zone := d.Get("zone").(string)
					return checkDiffName(old, new, zone)
				},
			},
			"linked_record": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The record name that will be linked to the CNAME record",
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
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v interface{}) string {
					return utils.JoinProperties(utils.ParseProperties(v.(string)))
				},
				DiffSuppressFunc: suppressWhenRemoteHasSuperset,
			},
			"to_deploy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether or not to selectively deploy the CNAME record",
				Default:     "no",
			},
			"batch_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether or not to use batch mode when selectively deploying",
				Default:     "disabled",
			},
		},
		Importer: &schema.ResourceImporter{
			State: recordImporter,
		},
	}
}

// createCNAMERecord Create the new CNAME record
func createCNAMERecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create CNAME record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	absoluteName := d.Get("absolute_name").(string)
	linkedRecord := d.Get("linked_record").(string)
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

	cnameRecord, err := objMgr.CreateCNAMERecord(configuration, view, zone, fqdnName, linkedRecord, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating CNAME record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	deploy := utils.ParseDeploymentValue(d.Get("to_deploy").(string))
	if deploy {
		cnameRecord.BatchMode = d.Get("batch_mode").(string)
		res, err := objMgr.Connector.DeployObject(cnameRecord)
		if err != nil {
			msg := fmt.Sprintf("Error deploying CNAME record %s: %s", fqdnName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		log.Debugf("Successfully deployed. %s", res)
	}
	d.Set("absolute_name", fqdnName)
	log.Debugf("Completed to create CNAME record %s", d.Get("absolute_name"))
	return getCNAMERecord(d, m)
}

// getCNAMERecord Get the CNAME record
func getCNAMERecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get CNAME record: %s", d.Get("absolute_name"))
	absoluteName, err := getAbsoluteName(d)
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	cnameRecord, err := objMgr.GetCNAMERecord(configuration, view, absoluteName)
	if err != nil {
		if utils.IsNotFoundErr(err) {
			if d.Id() != "" {
				// If the record is missing remotely, remove from state so Terraform plans a create.
				log.Warnf("CNAME Record %q not found; removing from state to trigger recreation", d.Id())
				d.SetId("")
				return nil
			}
			// If we don't have an ID yet (e.g., during import resolution) surface the not-found
			return fmt.Errorf("CNAME Record %s not found: %w", absoluteName, err)
		}
		// Any other error is a real failure
		return fmt.Errorf("Getting CNAME Record %s failed: %w", absoluteName, err)
	}
	// --- Parse both server and config properties ---
	bamProps := utils.ParseProperties(cnameRecord.Properties)
	cfgProps := utils.ParseProperties(d.Get("properties").(string))

	// --- Filter server properties using keys from config ---
	filteredProperties := utils.FilterProperties(bamProps, cfgProps)

	d.SetId(cnameRecord.AbsoluteName)
	d.Set("absolute_name", cnameRecord.AbsoluteName)
	d.Set("properties", utils.JoinProperties(filteredProperties))
	// for import functionality linked_record must be set for the cname_record - required attribute
	d.Set("linked_record", parseRecordPropertyValue(cnameRecord.Properties, "linkedRecordName"))
	log.Debugf("Completed reading CNAME record %s", d.Get("absolute_name"))
	return nil
}

// updateCNAMERecord Update the existing CNAME record
func updateCNAMERecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update CNAME record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	absoluteName := d.Get("absolute_name").(string)
	linkedRecord := d.Get("linked_record").(string)
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

	cnameRecord, err := objMgr.UpdateCNAMERecord(configuration, view, zone, fqdnName, linkedRecord, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating CNAME record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	deploy := utils.ParseDeploymentValue(d.Get("to_deploy").(string))
	if deploy {
		cnameRecord.BatchMode = d.Get("batch_mode").(string)
		res, err := objMgr.Connector.DeployObject(cnameRecord)
		if err != nil {
			msg := fmt.Sprintf("Error deploying CNAME record %s: %s", fqdnName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		log.Debugf("Successfully deployed. %s", res)
	}
	d.Set("absolute_name", fqdnName)
	log.Debugf("Completed to update CNAME record %s", d.Get("absolute_name"))
	return getCNAMERecord(d, m)
}

// deleteCNAMERecord Delete the CNAME record
func deleteCNAMERecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete CNAME record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	absoluteName := d.Get("absolute_name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteCNAMERecord(configuration, view, absoluteName)
	if err != nil {
		msg := fmt.Sprintf("Getting CNAME record %s failed: %s", absoluteName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Completed to delete CNAME record %s", d.Get("absolute_name"))
	return nil
}
