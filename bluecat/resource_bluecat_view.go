// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"errors"
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceZone The Zone
func ResourceView() *schema.Resource {

	return &schema.Resource{
		Create: createView,
		Read:   getView,
		Update: updateView,
		Delete: deleteView,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the View in the default Configuration if doesn't specify",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The View name",
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
		Importer: &schema.ResourceImporter{
			State: viewImporter,
		},
	}
}

// createView creates a new View
func createView(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create View %s", d.Get("address"))
	configuration := d.Get("configuration").(string)
	name := d.Get("name").(string)
	properties := d.Get("properties").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.CreateView(configuration, name, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating View (%s): %s", name, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	log.Debugf("Completed to create View %s", d.Get("name"))
	return getView(d, m)
}

// getView Get the View
func getView(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get Block %s", d.Get("address"))
	var viewName string
	var err error
	if d.Id() != "" {
		viewName = d.Id()
	} else {
		viewName = d.Get("name").(string)
	}

	configuration := d.Get("configuration").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	view, err := objMgr.GetView(configuration, viewName)
	if err != nil {
		if utils.IsNotFoundErr(err) {
			if d.Id() != "" {
				// If the record is missing remotely, remove from state so Terraform plans a create.
				log.Warnf("View %q not found; removing from state to trigger recreation", d.Id())
				d.SetId("")
				return nil
			}
			// If we don't have an ID yet (e.g., during import resolution) surface the not-found
			return fmt.Errorf("View %s not found: %w", viewName, err)
		}
		// Any other error is a real failure
		return fmt.Errorf("Getting View %s failed: %w", viewName, err)
	}

	// --- Parse both server and config properties ---
	bamProps := utils.ParseProperties(view.Properties)
	cfgProps := utils.ParseProperties(d.Get("properties").(string))

	// --- Filter server properties using keys from config ---
	filteredProperties := utils.FilterProperties(bamProps, cfgProps)

	d.Set("configuration", view.Configuration)
	d.Set("name", view.Name)
	d.Set("properties", utils.JoinProperties(filteredProperties))
	d.SetId(view.Name)
	log.Debugf("Completed getting View %s", d.Get("name"))
	return nil
}

// updateView Update the existing View - NOT IMPLEMENTED IN REST API
func updateView(d *schema.ResourceData, m interface{}) error {
	return errors.New("Updating View is not possible since it is not implemented in REST API.")
}

// deleteIP4Block Delete the IPv4 Block
func deleteView(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to Delete View %s", d.Get("name"))
	configuration := d.Get("configuration").(string)
	view_name := d.Get("name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteView(configuration, view_name)
	if err != nil {
		msg := fmt.Sprintf("Delete View %s failed: %s", view_name, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Deletion of View complete ")
	return nil
}
