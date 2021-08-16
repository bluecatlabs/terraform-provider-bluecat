// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ResourceNetwork The IPv4 Network
func ResourceDHCPRange() *schema.Resource {

	return &schema.Resource{
		Create: createDHCPRange,
		Read:   getDHCPRange,
		Update: updateDHCPRange,
		Delete: deleteDHCPRange,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the Network in the default Configuration if doesn't specify",
			},
			"network": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network address in CIDR format",
			},
			"start": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Start IP of the DHCP Range",
			},
			"end": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "End IP of the DHCP Range",
			},
			"properties": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DHCP Range's properties. Example: attribute=value|",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return checkDiffProperties(old, new)
				},
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 Template",
			},
		},
	}
}

// createDHCPRange Create the new DHCP Range
func createDHCPRange(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create DHCP Range (%s - %s) in the network %s", d.Get("start"), d.Get("end"), d.Get("network"))
	configuration := d.Get("configuration").(string)
	network := d.Get("network").(string)
	start := d.Get("start").(string)
	end := d.Get("end").(string)
	properties := d.Get("properties").(string)
	template := d.Get("template").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	//TODO: Check Network?

	_, err := objMgr.CreateDHCPRange(configuration, template, network, start, end, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating DHCP Range (%s - %s) in the network %s: %s", start, end, network, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	log.Debugf("Successful to create DHCP Range (%s - %s) in the network %s", start, end, network)

	return getDHCPRange(d, m)
}

// getDHCPRange Get the DHCP Range
func getDHCPRange(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get DHCP Range (%s - %s)", d.Get("start"), d.Get("end"))
	configuration := d.Get("configuration").(string)
	network := d.Get("network").(string)
	start := d.Get("start").(string)
	end := d.Get("end").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	dhcpRange, err := objMgr.GetDHCPRange(configuration, network, start, end)
	if err != nil {
		msg := fmt.Sprintf("Getting DHCP Range (%s - %s) failed: %s", start, end, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	d.SetId(dhcpRange.Start + "-" + dhcpRange.End)
	d.Set("network", dhcpRange.Network)
	d.Set("template", dhcpRange.Template)
	d.Set("properties", dhcpRange.Properties)
	log.Debugf("Completed getting DHCP Range (%s - %s)", start, end)
	return nil
}

// updateDHCPRange Update the existing DHCP Range
func updateDHCPRange(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update DHCP Range (%s - %s)", d.Get("start"), d.Get("end"))
	configuration := d.Get("configuration").(string)
	network := d.Get("network").(string)
	start := d.Get("start").(string)
	end := d.Get("end").(string)
	properties := d.Get("properties").(string)
	template := d.Get("template").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.UpdateDHCPRange(configuration, template, network, start, end, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating DHCP Range (%s - %s): %s", start, end, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	d.Set("start", getAttributeFromProperties("start", properties))
	d.Set("end", getAttributeFromProperties("end", properties))
	log.Debugf("Completed to update DHCP Range (%s - %s)", start, end)
	return getNetwork(d, m)
}

// deleteDHCPRange Delete the DHCP Range
func deleteDHCPRange(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete DHCP Range (%s - %s)", d.Get("start"), d.Get("end"))
	configuration := d.Get("configuration").(string)
	network := d.Get("network").(string)
	start := d.Get("start").(string)
	end := d.Get("end").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteDHCPRange(configuration, network, start, end)
	if err != nil {
		msg := fmt.Sprintf("Delete DHCP Range (%s - %s) failed: %s", start, end, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Deletion of DHCP Range complete ")
	return nil
}

func getAttributeFromProperties(attributeName string, props string) string {
	listProperties := strings.Split(props, "|")
	fmt.Println(listProperties)
	for i := 0; i < len(listProperties); i++ {
		prop := strings.Split(listProperties[i], "=")
		fmt.Println(prop)
		if prop[0] == attributeName && len(prop) == 2 {
			return prop[1]
		}
	}
	return ""
}
