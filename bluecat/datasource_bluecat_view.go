package bluecat

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceView() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceViewRead,
		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Getting the Zone in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, zone will be got under default view",
			},
			"deployable": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Zone's deployable property",
			},
			"server_roles": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of server roles. The format of each server role will be 'role type, server fqdn'",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
		},
	}
}

func dataSourceViewRead(d *schema.ResourceData, m interface{}) error {

	configuration := d.Get("configuration").(string)
	viewName := d.Get("view").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	viewObj, err := objMgr.GetView(configuration, viewName)
	if err != nil {
		msg := fmt.Sprintf("Getting Zone %s failed: %s", viewName, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	// Parse BAM properties
	bamProps := utils.ParseProperties(viewObj.Properties)
	d.Set("properties_raw", viewObj.Properties)

	filtered := utils.FilterDataSouceProperties(d, bamProps)

	// Write clean properties string back
	if err := d.Set("properties", utils.JoinProperties(filtered)); err != nil {
		return fmt.Errorf("setting properties failed: %w", err)
	}

	d.SetId(viewObj.Name)

	deployable := utils.GetPropertyValue("deployable", viewObj.Properties)
	if deployable == "true" {
		d.Set("deployable", "True")
	} else {
		d.Set("deployable", "False")
	}

	return nil
}
