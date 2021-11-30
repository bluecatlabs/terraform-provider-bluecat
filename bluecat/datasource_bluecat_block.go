package bluecat

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Optional:    true,
				Description: "IPv4 Block's properties",
			},
		},
	}
}

func dataSourceBlockRead(d *schema.ResourceData, m interface{}) error {

	configuration := d.Get("configuration").(string)
	cidr := d.Get("cidr").(string)

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

	block, err := objMgr.GetBlock(configuration, address, cidr)
	if err != nil {
		msg := fmt.Sprintf("Getting Block %s/%s failed: %s", address, cidr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	d.SetId(strconv.Itoa(block.BlockId))

	d.Set("name", block.Name)
	d.Set("properties", block.Properties)

	return nil
}
