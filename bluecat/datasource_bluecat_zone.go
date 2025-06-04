package bluecat

import (
	"fmt"
	"strconv"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZoneRead,
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
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The absolute name of zone or sub zone",
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
				Optional:    true,
				Description: "Zone's properties.",
			},
		},
	}
}

func dataSourceZoneRead(d *schema.ResourceData, m interface{}) error {

	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	zoneObj, err := objMgr.GetZone(configuration, view, zone)
	if err != nil {
		msg := fmt.Sprintf("Getting Zone %s failed: %s", zone, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	d.SetId(strconv.Itoa(zoneObj.ZoneId))
	d.Set("properties", zoneObj.Properties)

	deployable := utils.GetPropertyValue("deployable", zoneObj.Properties)
	if deployable == "true" {
		d.Set("deployable", "True")
	} else {
		d.Set("deployable", "False")
	}

	serverRoles, err := objMgr.GetDeploymentRoles(configuration, view, zone)
	if err != nil {
		msg := fmt.Sprintf("error get all deployment roles on the zone: %s", err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	var serverRolesRaw []string
	for _, serverRole := range serverRoles.ServerRoles {
		serverRoleRaw := fmt.Sprintf("%s, %s", getRoleNameInTerraform(serverRole.Role), serverRole.ServerFQDN)
		serverRolesRaw = append(serverRolesRaw, serverRoleRaw)
	}

	d.Set("server_roles", serverRolesRaw)

	return nil
}
