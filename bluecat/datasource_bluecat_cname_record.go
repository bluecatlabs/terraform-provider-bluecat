package bluecat

import (
	"fmt"
	"strconv"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceCNAMERecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCNAMERecordRead,
		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Getting the CNAME record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, record will be got under default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone in which you want to get a CNAME record",
			},
			"canonical": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the CNAME record. Must be FQDN if the Zone is not provided",
			},
			"linked_record": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The record name that will be linked to the CNAME record",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The TTL value",
			},
			"properties": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CNAME record's properties",
			},
		},
	}
}

func dataSourceCNAMERecordRead(d *schema.ResourceData, m interface{}) error {

	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	name := d.Get("canonical").(string)
	linkedRecord := d.Get("linked_record").(string)
	zone := d.Get("zone").(string)

	canonical := name
	if len(zone) > 0 {
		canonical = getFQDN(name, zone)
	}

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	cnameRecord, err := objMgr.GetCNAMERecord(configuration, view, canonical)
	if err != nil {
		msg := fmt.Sprintf("Getting CNAME record %s failed: %s", canonical, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	currentLinkedRecord := getPropertyValue("linkedRecordName", cnameRecord.Properties)

	if linkedRecord != currentLinkedRecord {
		msg := fmt.Sprintf("Getting CNAME record %s failed: linkedRecordName %s isn't matching", canonical, linkedRecord)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	if len(zone) == 0 {
		zone = getZoneFromFQDN(cnameRecord.Name, canonical)
	}

	ttl := -1
	ttlStr := getPropertyValue("ttl", cnameRecord.Properties)
	if ttlStr != "" {
		if ttlInt, err := strconv.Atoi(ttlStr); err == nil {
			ttl = ttlInt
		}
	}

	d.SetId(strconv.Itoa(cnameRecord.CNameId))

	d.Set("zone", zone)
	d.Set("ttl", ttl)
	d.Set("properties", cnameRecord.Properties)

	return nil
}
