package bluecat

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBlock() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlockRead,
		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Getting the IPv4 Block in the default Configuration if doesn't specify",
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IPv4 Block's CIDR",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Block name",
			},
			"properties": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pipe-separated key=value properties (filtered).",
			},
			"properties_raw": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unfiltered raw properties returned by BAM.",
			},
			"allowed_property_keys": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Optional list of property keys to keep when filtering.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Block IP version: ipv4 or ipv6",
			},
		},
	}
}

func dataSourceBlockRead(d *schema.ResourceData, m interface{}) error {

	configuration := d.Get("configuration").(string)
	cidr := d.Get("cidr").(string)
	ipVersion := d.Get("ip_version").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	if !(strings.Contains(cidr, "/")) {
		msg := fmt.Sprintf("Invalid cidr block %s", cidr)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	cidrList := strings.Split(cidr, "/")
	address, cidr := cidrList[0], cidrList[1]

	block, err := objMgr.GetBlock(configuration, address, cidr, ipVersion)
	if err != nil {
		msg := fmt.Sprintf("Getting Block %s/%s failed: %s", address, cidr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	// Parse BAM properties
	bamProps := utils.ParseProperties(block.Properties)
	d.Set("properties_raw", block.Properties)

	filtered := utils.FilterDataSouceProperties(d, bamProps)

	// Write clean properties string back
	if err := d.Set("properties", utils.JoinProperties(filtered)); err != nil {
		return fmt.Errorf("setting properties failed: %w", err)
	}

	d.SetId(strconv.Itoa(block.BlockId))

	d.Set("name", block.Name)

	return nil
}
