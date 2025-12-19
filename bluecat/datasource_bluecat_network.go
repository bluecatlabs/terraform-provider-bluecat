package bluecat

import (
	"fmt"
	"strconv"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIPv4Network() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIPv4NetworkRead,
		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Getting the IPv4 Network in the default Configuration if doesn't specify",
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network address in CIDR format",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Network name",
			},
			"gateway": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Gateway address",
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
				Description: "Network's IP version",
			},
		},
	}
}

func dataSourceIPv4NetworkRead(d *schema.ResourceData, m interface{}) error {

	network := entities.Network{}
	if !network.InitNetwork(d) {
		log.Error(network.InitError)
		return fmt.Errorf(network.InitError)
	}

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	retrievedNetwork, err := objMgr.GetNetwork(&network)
	if err != nil {
		msg := fmt.Sprintf("Getting Network %s failed: %s", retrievedNetwork.CIDR, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	gateway := utils.GetPropertyValue("gateway", retrievedNetwork.Properties)
	// Parse BAM properties
	bamProps := utils.ParseProperties(retrievedNetwork.Properties)
	d.Set("properties_raw", retrievedNetwork.Properties)

	filtered := utils.FilterDataSouceProperties(d, bamProps)

	// Write clean properties string back
	if err := d.Set("properties", utils.JoinProperties(filtered)); err != nil {
		return fmt.Errorf("setting properties failed: %w", err)
	}

	d.SetId(strconv.Itoa(retrievedNetwork.NetWorkId))

	d.Set("name", retrievedNetwork.Name)
	d.Set("gateway", gateway)

	return nil
}
