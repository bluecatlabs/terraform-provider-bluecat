// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/logging"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	log = *logging.GetLogger()
}

// ResourceExternalHostRecord The ExternalHost record
func ResourceExternalHostRecord() *schema.Resource {
	return &schema.Resource{
		Create: createExternalHostRecord,
		Read:   getExternalHostRecord,
		Update: updateExternalHostRecord,
		Delete: deleteExternalHostRecord,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the External Host record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the record. If not provided, record will be created under default view",
			},
			"absolute_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ExternalHost record. Must be an FQDN.",
			},
			"addresses": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP addresses that will be linked to the External Host record",
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
				Description: "Whether or not to selectively deploy the Host record",
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

// createExternalHostRecord Create the new ExternalHost record
func createExternalHostRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create ExternalHost record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	absoluteName := d.Get("absolute_name").(string)
	addresses := d.Get("addresses").(string)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	externalHostRecord, err := objMgr.CreateExternalHostRecord(configuration, view, addresses, absoluteName, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating ExternalHost record %s: %s", absoluteName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	deploy := utils.ParseDeploymentValue(d.Get("to_deploy").(string))
	if deploy {
		externalHostRecord.BatchMode = d.Get("batch_mode").(string)
		res, err := objMgr.Connector.DeployObject(externalHostRecord)
		if err != nil {
			msg := fmt.Sprintf("Error deploying External Host record %s: %s", absoluteName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		log.Debugf("Successfully deployed. %s", res)
	}
	d.Set("absolute_name", absoluteName)
	log.Debugf("Completed to create ExternalHost record %s", d.Get("absolute_name"))
	return getExternalHostRecord(d, m)
}

// getExternalHostRecord Get the ExternalHost record
func getExternalHostRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get ExternalHost record: %s", d.Get("absolute_name"))
	absoluteName, err := getAbsoluteName(d)
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	externalHostRecord, err := objMgr.GetExternalHostRecord(configuration, view, absoluteName)
	if err != nil {
		if utils.IsNotFoundErr(err) {
			if d.Id() != "" {
				// If the record is missing remotely, remove from state so Terraform plans a create.
				log.Warnf("External Host record %q not found; removing from state to trigger recreation", d.Id())
				d.SetId("")
				return nil
			}
			// If we don't have an ID yet (e.g., during import resolution) surface the not-found
			return fmt.Errorf("External Host record %s not found: %w", absoluteName, err)
		}
		// Any other error is a real failure
		return fmt.Errorf("Getting External Host Record %s failed: %w", absoluteName, err)
	}

	// --- Parse both server and config properties ---
	bamProps := utils.ParseProperties(externalHostRecord.Properties)
	cfgProps := utils.ParseProperties(d.Get("properties").(string))

	// --- Filter server properties using keys from config ---
	filteredProperties := utils.FilterProperties(bamProps, cfgProps)
	d.SetId(externalHostRecord.AbsoluteName)
	d.Set("absolute_name", externalHostRecord.AbsoluteName)
	d.Set("properties", utils.JoinProperties(filteredProperties))
	// for import functionality ip4_address must be set for the host_record - required attribute
	d.Set("addresses", externalHostRecord.Addresses)
	log.Debugf("Completed reading ExternalHost record %s", d.Get("absolute_name"))
	return nil
}

// updateExternalHostRecord Update the existing ExternalHost record
func updateExternalHostRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update ExternalHost record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	addresses := d.Get("addresses").(string)
	absoluteName := d.Get("absolute_name").(string) // new absolute name
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	var immutableProperties = []string{"parentId", "parentType"} // these properties will raise error on the rest-api
	properties = utils.RemoveImmutableProperties(properties, immutableProperties)

	externalHostRecord, err := objMgr.UpdateExternalHostRecord(configuration, view, addresses, absoluteName, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating ExternalHost record %s: %s", absoluteName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	deploy := utils.ParseDeploymentValue(d.Get("to_deploy").(string))
	if deploy {
		externalHostRecord.BatchMode = d.Get("batch_mode").(string)
		res, err := objMgr.Connector.DeployObject(externalHostRecord)
		if err != nil {
			msg := fmt.Sprintf("Error deploying External Host record %s: %s", absoluteName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		log.Debugf("Successfully deployed. %s", res)
	}
	d.Set("absolute_name", absoluteName)
	log.Debugf("Completed to update ExternalHost record %s", d.Get("absolute_name"))
	return getExternalHostRecord(d, m)
}

// deleteExternalHostRecord Delete the ExternalHost record
func deleteExternalHostRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete ExternalHost record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	absoluteName := d.Get("absolute_name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	// Check the host exist or not
	_, err := objMgr.GetExternalHostRecord(configuration, view, absoluteName)
	if err != nil {
		log.Debugf("ExternalHost record %s not found", absoluteName)
	} else {
		_, err := objMgr.DeleteExternalHostRecord(configuration, view, absoluteName)
		if err != nil {
			msg := fmt.Sprintf("Delete ExternalHost record %s failed: %s", absoluteName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	}
	d.SetId("")
	log.Debugf("Completed to delete ExternalHost record %s", d.Get("absolute_name"))
	return nil
}
