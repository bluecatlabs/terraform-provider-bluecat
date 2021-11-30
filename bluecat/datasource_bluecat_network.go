package bluecat

import (
	"fmt"
	"strconv"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Optional:    true,
				Description: "IPv4 Network's properties",
			},
		},
	}
}

func dataSourceIPv4NetworkRead(d *schema.ResourceData, m interface{}) error {

	configuration := d.Get("configuration").(string)
	cidr := d.Get("cidr").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	network, err := objMgr.GetNetwork(configuration, cidr)
	if err != nil {
		msg := fmt.Sprintf("Getting Network %s failed: %s", cidr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	gateway := getPropertyValue("gateway", network.Properties)

	d.SetId(strconv.Itoa(network.NetWorkId))

	d.Set("name", network.Name)
	d.Set("gateway", gateway)
	d.Set("properties", network.Properties)

	return nil
}
