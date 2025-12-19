// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceDHCPRange The IPv4 ResourceDHCPRange
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
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the DHCP Range",
			},
			"properties": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v interface{}) string {
					return utils.JoinProperties(utils.ParseProperties(v.(string)))
				},
				DiffSuppressFunc: suppressWhenRemoteHasSuperset,
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IPv4 Template",
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DHCPRange's IP version",
			},
		},
	}
}

// createDHCPRange Create the new DHCP Range
func createDHCPRange(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	dhcpRange := entities.DHCPRange{}
	if !dhcpRange.InitRange(d) {
		log.Error(dhcpRange.InitRange)
		return fmt.Errorf(dhcpRange.InitError)
	}

	log.Debugf("Beginning to create DHCP Range (%s - %s) in the network %s", dhcpRange.Start, dhcpRange.End, dhcpRange.Network)

	//TODO: Check Network?

	_, err := objMgr.CreateDHCPRange(dhcpRange)
	if err != nil {
		msg := fmt.Sprintf("Error creating DHCP Range (%s - %s) in the network %s: %s", dhcpRange.Start, dhcpRange.End, dhcpRange.Network, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	log.Debugf("Successful to create DHCP Range (%s - %s) in the network %s", dhcpRange.Start, dhcpRange.End, dhcpRange.Network)

	return getDHCPRange(d, m)
}

// getDHCPRange Get the DHCP Range
func getDHCPRange(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	dhcpRange := entities.DHCPRange{}
	if !dhcpRange.InitRange(d) {
		log.Error(dhcpRange.InitRange)
		return fmt.Errorf(dhcpRange.InitError)
	}

	log.Debugf("Beginning to get DHCP Range (%s - %s)", dhcpRange.Start, dhcpRange.End)

	dhcpRangeEntity, err := objMgr.GetDHCPRange(dhcpRange)
	if err != nil {
		msg := fmt.Sprintf("Getting DHCP Range (%s - %s) failed: %s", dhcpRangeEntity.Start, dhcpRangeEntity.End, err)
		msg += fmt.Sprintf("Subpath: %s", dhcpRangeEntity.SubPath())
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	// --- Parse both server and config properties ---
	bamProps := utils.ParseProperties(dhcpRange.Properties)
	cfgProps := utils.ParseProperties(d.Get("properties").(string))

	// --- Filter server properties using keys from config ---
	filteredProperties := utils.FilterProperties(bamProps, cfgProps)

	d.SetId(dhcpRangeEntity.Start + "-" + dhcpRangeEntity.End)
	d.Set("network", dhcpRangeEntity.Network)
	d.Set("template", dhcpRangeEntity.Template)
	d.Set("properties", utils.JoinProperties(filteredProperties))
	log.Debugf("Completed getting DHCP Range (%s - %s)", dhcpRangeEntity.Start, dhcpRangeEntity.End)
	return nil
}

// updateDHCPRange Update the existing DHCP Range
func updateDHCPRange(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	dhcpRange := entities.DHCPRange{}
	if !dhcpRange.InitRange(d) {
		log.Error(dhcpRange.InitRange)
		return fmt.Errorf(dhcpRange.InitError)
	}

	log.Debugf("Beginning to update DHCP Range (%s - %s)", dhcpRange.Start, dhcpRange.End)

	dhcpRange.Template = d.Get("template").(string)
	if dhcpRange.Template == "" {
		dhcpRange.Template = " "
	}

	_, err := objMgr.UpdateDHCPRange(dhcpRange)
	if err != nil {
		msg := fmt.Sprintf("Error updating DHCP Range (%s - %s): %s", dhcpRange.Start, dhcpRange.End, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	startAfterUpdate := getAttributeFromProperties("start", dhcpRange.Properties)
	endAfterUpdate := getAttributeFromProperties("end", dhcpRange.Properties)

	if startAfterUpdate != "" && endAfterUpdate != "" {
		dhcpRange.Start = startAfterUpdate
		dhcpRange.End = endAfterUpdate
	}

	d.Set("start", dhcpRange.Start)
	d.Set("end", dhcpRange.End)

	log.Debugf("Completed to update DHCP Range (%s - %s)", dhcpRange.Start, dhcpRange.End)
	return getDHCPRange(d, m)
}

// deleteDHCPRange Delete the DHCP Range
func deleteDHCPRange(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	dhcpRange := entities.DHCPRange{}
	if !dhcpRange.InitRange(d) {
		log.Error(dhcpRange.InitRange)
		return fmt.Errorf(dhcpRange.InitError)
	}

	log.Debugf("Beginning to delete DHCP Range (%s - %s)", dhcpRange.Start, dhcpRange.End)

	_, err := objMgr.DeleteDHCPRange(dhcpRange)
	if err != nil {
		msg := fmt.Sprintf("Delete DHCP Range (%s - %s) failed: %s", dhcpRange.Start, dhcpRange.End, err)
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
