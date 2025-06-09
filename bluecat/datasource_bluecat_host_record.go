// Copyright 2021 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceHostRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostRecordRead,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Getting the Host record in the default Configuration if doesn't specify.",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, record will be got under default view.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone which contains the details of the Host record.",
			},
			"fqdn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Host record. Must be FQDN if the Zone is not provided",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP Address that will be linked to the Host record",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The TTL value",
			},
			"properties": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host record's properties. Example: attribute=value|",
			},
		},
	}
}

func dataSourceHostRecordRead(d *schema.ResourceData, m interface{}) error {
	configuration := d.Get("configuration").(string)

	view := d.Get("view").(string)
	name := d.Get("fqdn").(string)
	ipAddress := d.Get("ip_address").(string)
	zone := d.Get("zone").(string)

	fqdnName := name
	if len(zone) > 0 {
		fqdnName = getFQDN(name, zone)
	}

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	hostRecord, err := objMgr.GetHostRecord(configuration, view, fqdnName)
	if err != nil {
		msg := fmt.Sprintf("Getting Host record %s failed: %s", fqdnName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	ipLinked := utils.GetPropertyValue("addresses", hostRecord.Properties)

	if !(strings.Contains(ipLinked, ipAddress)) {
		msg := fmt.Sprintf("Getting Host record %s failed: IP Address %s isn't matching", fqdnName, ipAddress)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	if len(zone) == 0 {
		zone = getZoneFromFQDN(hostRecord.Name, fqdnName)
	}

	ttl := -1
	ttlStr := utils.GetPropertyValue("ttl", hostRecord.Properties)
	if ttlStr != "" {
		if ttlInt, err := strconv.Atoi(ttlStr); err == nil {
			ttl = ttlInt
		}
	}

	d.SetId(strconv.Itoa(hostRecord.HostId))
	d.Set("zone", zone)
	d.Set("ttl", ttl)
	d.Set("properties", hostRecord.Properties)
	log.Debugf("Completed reading Host record %s", fqdnName)

	return nil
}

func getZoneFromFQDN(rrName, fqdnName string) string {
	if rrName == "" {
		return fqdnName
	}
	return strings.Split(fqdnName, rrName)[1][1:]
}
