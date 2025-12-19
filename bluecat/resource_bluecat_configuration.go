// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceConfiguration The Configuration resource
func ResourceConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: createConfiguration,
		Read:   getConfiguration,
		Update: updateConfiguration,
		Delete: deleteConfiguration,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Configuration name.",
			},
			"properties": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v interface{}) string {
					return utils.JoinProperties(utils.ParseProperties(v.(string)))
				},
				DiffSuppressFunc: suppressWhenRemoteHasSuperset,
			},
		},
	}
}

// createConfiguration Create the new Configuration
func createConfiguration(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create Configuration %s", d.Get("name"))
	name := d.Get("name").(string)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.CreateConfiguration(name, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating Configuration (%s): %s", name, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	log.Debugf("Completed to create Configuration %s", d.Get("name"))
	return getConfiguration(d, m)
}

// getConfiguration Get the Configuration
func getConfiguration(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get Configuration %s", d.Get("name"))
	name := d.Get("name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	config, err := objMgr.GetConfiguration(name)
	if err != nil {
		msg := fmt.Sprintf("Getting Configuration %s failed: %s", name, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	// --- Parse both server and config properties ---
	bamProps := utils.ParseProperties(config.Properties)
	cfgProps := utils.ParseProperties(d.Get("properties").(string))

	// --- Filter server properties using keys from config ---
	filteredProperties := utils.FilterProperties(bamProps, cfgProps)
	d.SetId(config.Name)
	d.Set("name", config.Name)
	d.Set("properties", utils.JoinProperties(filteredProperties))
	log.Debugf("Completed getting Configuration %s", d.Get("name"))
	return nil
}

// updateConfiguration Update the existing Configuration
func updateConfiguration(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update Configuration %s", d.Get("name"))
	name := d.Get("name").(string)
	props := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	config, err := objMgr.UpdateConfiguration(name, props)
	if err != nil {
		msg := fmt.Sprintf("Updating Configuration %s failed: %s", name, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.Set("name", config.Name)
	log.Debugf("Completed to update Configuration %s", d.Get("name"))
	return getConfiguration(d, m)
}

// deleteConfiguration Delete the Configuration
func deleteConfiguration(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete Configuration %s", d.Get("name"))
	name := d.Get("name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteConfiguration(name)
	if err != nil {
		msg := fmt.Sprintf("Delete Configuration %s failed: %s", name, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Deletion of Configuration complete")
	return nil
}
