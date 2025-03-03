// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/logging"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sirupsen/logrus"
)

var log logrus.Logger

func init() {
	log = *logging.GetLogger()
}

// ResourceHostRecord The Host record
func ResourceHostRecord() *schema.Resource {
	return &schema.Resource{
		Create: createHostRecord,
		Read:   getHostRecord,
		Update: updateHostRecord,
		Delete: deleteHostRecord,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the Host record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, record will be created under default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone in which you want to update a host record. If not provided, the absolute name must be FQDN ones",
			},
			"absolute_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Host record. Must be FQDN if the Zone is not provided",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					zone := d.Get("zone").(string)
					return checkDiffName(old, new, zone)
				},
			},
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP address that will be linked to the Host record",
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
				Description: "Host record's properties. Example: attribute=value|",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return checkDiffProperties(old, new)
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: recordImporter,
		},
	}
}

func recordParseId(id string) (string, string, error) {
	// this func will be used for host and CNAME records
	parts := strings.SplitN(id, ".", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of host record ID (%s), expected zone.host_record_name", id)
	}
	recordName := parts[0]
	zoneName := parts[1]

	return zoneName, recordName, nil
}

// createHostRecord Create the new Host record
func createHostRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create Host record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	absoluteName := d.Get("absolute_name").(string)
	ipAddress := d.Get("ip_address").(string)
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

	_, err := objMgr.CreateHostRecord(configuration, view, zone, fqdnName, ipAddress, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating Host record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.Set("absolute_name", fqdnName)
	log.Debugf("Completed to create Host record %s", d.Get("absolute_name"))
	return getHostRecord(d, m)
}

// getHostRecord Get the Host record
func getHostRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get Host record: %s", d.Get("absolute_name"))
	absoluteName, err := getAbsoluteName(d)
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	hostRecord, err := objMgr.GetHostRecord(configuration, view, absoluteName)
	if err != nil {
		if d.Id() != "" {
			err := createHostRecord(d, m)
			if err != nil {
				msg := fmt.Sprintf("Something gone wrong: %v", err)
				return fmt.Errorf(msg)
			}
		} else {
			msg := fmt.Sprintf("Getting Host record %s failed: %s", absoluteName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	}
	d.SetId(hostRecord.AbsoluteName)
	d.Set("absolute_name", hostRecord.AbsoluteName)
	d.Set("properties", hostRecord.Properties)
	// for import functionality ip4_address must be set for the host_record - required attribute
	d.Set("ip_address", parseRecordPropertyValue(hostRecord.Properties, "addresses"))
	log.Debugf("Completed reading Host record %s", d.Get("absolute_name"))
	return nil
}

// updateHostRecord Update the existing Host record
func updateHostRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update Host record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	absoluteName := d.Get("absolute_name").(string)
	ipAddress := d.Get("ip_address").(string)
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

	_, err := objMgr.UpdateHostRecord(configuration, view, zone, fqdnName, ipAddress, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating Host record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.Set("absolute_name", fqdnName)
	log.Debugf("Completed to update Host record %s", d.Get("absolute_name"))
	return getHostRecord(d, m)
}

// deleteHostRecord Delete the Host record
func deleteHostRecord(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete Host record %s", d.Get("absolute_name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	absoluteName := d.Get("absolute_name").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	// Check the host exist or not
	_, err := objMgr.GetHostRecord(configuration, view, absoluteName)
	if err != nil {
		log.Debugf("Host record %s not found", absoluteName)
	} else {
		_, err := objMgr.DeleteHostRecord(configuration, view, absoluteName)
		if err != nil {
			msg := fmt.Sprintf("Delete Host record %s failed: %s", absoluteName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
	}
	d.SetId("")
	log.Debugf("Completed to delete Host record %s", d.Get("absolute_name"))
	return nil
}

func getFQDN(rrName, zone string) string {
	if !strings.HasSuffix(rrName, ".") && len(zone) > 0 && !strings.HasSuffix(rrName, zone) {
		return fmt.Sprintf("%s.%s", rrName, zone)
	}
	return rrName
}

func getZoneFromRRName(rrName string) (zoneFQDN string) {
	zoneFQDN = ""
	index := strings.Index(rrName, ".")
	if index > 0 {
		zoneFQDN = rrName[index+1:]
	}
	return
}

func checkDiffProperties(old string, new string) bool {
	newProperties := strings.Split(new, "|")
	for i := 0; i < len(newProperties); i++ {
		if newProperties[i] != "" && !strings.Contains(fmt.Sprintf("|%s|", old), fmt.Sprintf("|%s|", newProperties[i])) {
			return false
		}
	}
	return true
}

func checkDiffName(old string, new string, zone string) bool {
	if old == getFQDN(new, zone) {
		return true
	}
	return false
}

func removeAttributeFromProperties(attributeName string, props string) string {
	listProperties := strings.Split(props, "|")
	properties := ""
	for i := 0; i < len(listProperties); i++ {
		prop := strings.Split(listProperties[i], "=")
		if prop[0] != attributeName && prop[0] != "" {
			properties += listProperties[i] + "|"
		}
	}
	return properties
}
