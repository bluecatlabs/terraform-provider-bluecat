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

	// Make sure the reverseRecord property is properly capitalized (if it exists)
	properties, err := fixReverseRecordPropIfExists(properties)

	hostRecord, err := objMgr.CreateHostRecord(configuration, view, zone, fqdnName, ipAddress, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating Host record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	deploy := utils.ParseDeploymentValue(d.Get("to_deploy").(string))
	if deploy {
		hostRecord.BatchMode = d.Get("batch_mode").(string)
		res, err := objMgr.Connector.DeployObject(hostRecord)
		if err != nil {
			msg := fmt.Sprintf("Error deploying Host record %s: %s", fqdnName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		log.Debugf("Successfully deployed. %s", res)
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
		if utils.IsNotFoundErr(err) {
			if d.Id() != "" {
				// If the record is missing remotely, remove from state so Terraform plans a create.
				log.Warnf("Host record %q not found; removing from state to trigger recreation", d.Id())
				d.SetId("")
				return nil
			}
			// If we don't have an ID yet (e.g., during import resolution) surface the not-found
			return fmt.Errorf("host record %s not found: %w", absoluteName, err)
		}
		// Any other error is a real failure
		return fmt.Errorf("getting host record %s failed: %w", absoluteName, err)
	}

	// --- Parse both server and config properties ---
	bamProps := utils.ParseProperties(hostRecord.Properties)
	cfgProps := utils.ParseProperties(d.Get("properties").(string))

	// --- Filter server properties using keys from config ---
	filteredProperties := utils.FilterProperties(bamProps, cfgProps)

	d.SetId(hostRecord.AbsoluteName)
	d.Set("absolute_name", hostRecord.AbsoluteName)
	d.Set("properties", utils.JoinProperties(filteredProperties))
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

	// Make sure the reverseRecord property is properly capitalized (if it exists)
	properties, err := fixReverseRecordPropIfExists(properties)

	var immutableProperties = []string{"parentId", "parentType"} // these properties will raise error on the rest-api
	properties = utils.RemoveImmutableProperties(properties, immutableProperties)

	hostRecord, err := objMgr.UpdateHostRecord(configuration, view, zone, fqdnName, ipAddress, ttl, properties)
	if err != nil {
		msg := fmt.Sprintf("Error updating Host record %s: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	deploy := utils.ParseDeploymentValue(d.Get("to_deploy").(string))
	if deploy {
		hostRecord.BatchMode = d.Get("batch_mode").(string)
		res, err := objMgr.Connector.DeployObject(hostRecord)
		if err != nil {
			msg := fmt.Sprintf("Error deploying Host record %s: %s", fqdnName, err)
			log.Debug(msg)
			return fmt.Errorf(msg)
		}
		log.Debugf("Successfully deployed. %s", res)
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

	_, err := objMgr.DeleteHostRecord(configuration, view, absoluteName)
	if err != nil {
		msg := fmt.Sprintf("Delete Host record %s failed: %s", absoluteName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
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

func suppressWhenRemoteHasSuperset(k, old, new string, d *schema.ResourceData) bool {
	oldProps := utils.ParseProperties(old)
	newProps := utils.ParseProperties(new)
	if len(newProps) == 0 {
		return true
	}
	for key, newVal := range newProps {
		oldVal, exists := oldProps[key]
		if !exists {
			return false
		}
		if oldVal != newVal {
			return false
		}
	}

	return true
}

func fixReverseRecordPropIfExists(properties string) (string, error) {
	if properties == "" {
		return properties, nil
	}

	segments := strings.Split(properties, "|")
	for i, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}

		// Split key=value
		kv := strings.SplitN(seg, "=", 2)
		if len(kv) != 2 {
			continue // ignore malformed segment
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		if strings.EqualFold(key, "reverseRecord") {
			switch strings.ToLower(val) {
			case "yes", "true":
				val = "true"
			case "no", "false":
				val = "false"
			default:
				return properties, fmt.Errorf("invalid value for reverseRecord: %q (must be yes/no/true/false)", val)
			}
			segments[i] = fmt.Sprintf("%s=%s", key, val)
			break
		}
	}

	return strings.Join(segments, "|"), nil
}
