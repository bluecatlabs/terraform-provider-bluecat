// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceZone The Zone
func ResourceZone() *schema.Resource {

	return &schema.Resource{
		Create: createZone,
		Read:   getZone,
		Update: updateZone,
		Delete: deleteZone,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the Zone in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, zone will be created under default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The absolute name of zone or sub zone",
			},
			"deployable": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The deployable flag is False by default and is optional. To make the zone deployable, set the deployable flag to True.",
			},
			"server_roles": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of server roles. The format of each server role will be 'role type, server fqdn'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"deployment_options": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The deployment options for the zone.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"properties": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v interface{}) string {
					return utils.JoinProperties(utils.ParseProperties(v.(string)))
				},
				DiffSuppressFunc: suppressWhenRemoteHasSuperset,
			},
		},
		Importer: &schema.ResourceImporter{
			State: zoneImporter,
		},
	}
}

// createZone Create the new Zone
// Create the Host record, then server will create the PTR
func createZone(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create Zone %s", d.Get("name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	deployable := d.Get("deployable").(string)
	serverRolesRaw := d.Get("server_roles").([]interface{})
	deploymentOptions := utils.ExpandStringMap(d.Get("deployment_options"))
	properties := d.Get("properties").(string)

	properties, err := updateDeployableProperty(deployable, properties, false)
	if err != nil {
		msg := fmt.Sprintf("Error creating Zone (%s): %s", zone, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err = objMgr.CreateZone(configuration, view, zone, properties)
	if err != nil {
		msg := fmt.Sprintf("Error creating Zone (%s): %s", zone, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	serverRoles := make([]string, len(serverRolesRaw))
	count := 0
	for i, raw := range serverRolesRaw {
		svrRole, ok := raw.(string)
		if !ok {
			log.Warning("Cannot load server role item at: ", i)
		}
		svrRole = strings.TrimSpace(svrRole)
		if len(svrRole) > 0 {
			serverRoles[count] = svrRole
			count += 1
		}
	}
	for _, serverRole := range serverRoles {
		if len(serverRole) > 0 {
			role, serverFQDN, err := validateServerRole(objMgr, configuration, serverRole)
			if err == nil {
				_, err = objMgr.CreateDeploymentRole(configuration, view, zone, serverFQDN, "dns", role, "", "")
			}

			if err != nil {
				msg := fmt.Sprintf("Error creating Zone (%s): %s", zone, err)
				log.Debug(msg)

				_, err := objMgr.DeleteZone(configuration, view, zone)
				if err != nil {
					msg := fmt.Sprintf("Rollback data - Delete Zone %s failed: %s", zone, err)
					log.Debug(msg)
				}
				return fmt.Errorf(msg)
			}
		}
	}

	for optionName, optionValue := range deploymentOptions {
		_, err = objMgr.CreateDeploymentOption(entities.DeploymentOption{
			Configuration: configuration,
			View:          view,
			Zone:          zone,
			Name:          optionName,
			Value:         optionValue,
		})
		if err != nil {
			msg := fmt.Sprintf("Error creating Zone (%s): %s", zone, err)
			log.Debug(msg)

			_, err := objMgr.DeleteZone(configuration, view, zone)
			if err != nil {
				msg := fmt.Sprintf("Rollback data - Delete Zone %s failed: %s", zone, err)
				log.Debug(msg)
			}
			return fmt.Errorf(msg)
		}
	}

	log.Debugf("Completed to create Zone %s", d.Get("zone"))

	return getZone(d, m)
}

// getZone Get the Zone
func getZone(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get Zone: %s", d.Get("name"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	zoneObj, err := objMgr.GetZone(configuration, view, zone)
	if err != nil {
		if utils.IsNotFoundErr(err) {
			if d.Id() != "" {
				// If the record is missing remotely, remove from state so Terraform plans a create.
				log.Warnf("Zone %q not found; removing from state to trigger recreation", d.Id())
				d.SetId("")
				return nil
			}
			// If we don't have an ID yet (e.g., during import resolution) surface the not-found
			return fmt.Errorf("Zone %s not found: %w", zone, err)
		}
		// Any other error is a real failure
		return fmt.Errorf("Getting Zone %s failed: %w", zone, err)
	}

	// --- Parse both server and config properties ---
	bamProps := utils.ParseProperties(zoneObj.Properties)
	cfgProps := utils.ParseProperties(d.Get("properties").(string))

	// --- Filter server properties using keys from config ---
	filteredProperties := utils.FilterProperties(bamProps, cfgProps)
	// During import there are usually no configured property keys yet.
	// Keep full BAM properties so generated config is populated.
	if len(cfgProps) == 0 {
		filteredProperties = bamProps
	}

	d.SetId(zone)
	d.Set("zone", zone)
	d.Set("configuration", configuration)
	d.Set("view", view)
	d.Set("properties", utils.JoinProperties(filteredProperties))

	deployable := utils.GetPropertyValue("deployable", zoneObj.Properties)
	if strings.EqualFold(deployable, "true") {
		d.Set("deployable", "true")
	} else if strings.EqualFold(deployable, "false") {
		d.Set("deployable", "false")
	}

	serverRoles, rolesErr := objMgr.GetDeploymentRoles(configuration, view, zone)
	if rolesErr != nil {
		msg := fmt.Sprintf("error get all deployment roles on the zone: %s", rolesErr)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	serverRolesRaw := make([]string, 0, len(serverRoles.ServerRoles))
	for _, serverRole := range serverRoles.ServerRoles {
		serverRoleRaw := fmt.Sprintf("%s, %s", getRoleNameInTerraform(serverRole.Role), serverRole.ServerFQDN)
		serverRolesRaw = append(serverRolesRaw, serverRoleRaw)
	}
	d.Set("server_roles", serverRolesRaw)

	deploymentOptionNames := utils.GetSortedMapKeys(utils.ExpandStringMap(d.Get("deployment_options")))
	deploymentOptionsRaw, optionsErr := getDeploymentOptions(objMgr, configuration, view, zone, deploymentOptionNames)
	if optionsErr != nil {
		msg := fmt.Sprintf("error get deployment options on the zone: %s", optionsErr)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.Set("deployment_options", utils.FlattenStringMap(deploymentOptionsRaw))

	log.Debugf("Completed reading Zone %s", d.Get("zone"))
	return nil
}

// updateZone Update the existing Zone
func updateZone(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update Zone %s", d.Get("zone"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	deployable := d.Get("deployable").(string)
	serverRolesRaw := d.Get("server_roles").([]interface{})
	properties := d.Get("properties").(string)

	properties, err := updateDeployableProperty(deployable, properties, true)
	if err != nil {
		msg := fmt.Sprintf("Error updating Zone (%s): %s", zone, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	newServerRoles, currentServerRoles, err := prepareServerRoleData(objMgr, serverRolesRaw, configuration, view, zone)
	if err != nil {
		msg := fmt.Sprintf("Error updating Zone (%s): %s", zone, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}

	currentDeploymentOptionsRaw, newDeploymentOptionsRaw := d.GetChange("deployment_options")
	newDeploymentOptions, currentDeploymentOptions := prepareDeploymentOptionData(currentDeploymentOptionsRaw, newDeploymentOptionsRaw)

	trace, err := updateServerRoles(objMgr, currentServerRoles, newServerRoles, configuration, view, zone)
	if err == nil {
		trace, err = updateDeploymentOptions(objMgr, currentDeploymentOptions, newDeploymentOptions, configuration, view, zone, trace)
	}
	if err == nil {
		_, err = objMgr.UpdateZone(configuration, view, zone, properties)
	}

	if err != nil {
		msg := fmt.Sprintf("Error updating Zone (%s): %s", zone, err)
		log.Debug(msg)

		err := rollBackData(objMgr, trace, configuration, view, zone)
		if err != nil {
			msg := fmt.Sprintf("Rollback data failed: %s", err)
			log.Debug(msg)
		}

		return fmt.Errorf(msg)
	}

	return getZone(d, m)
}

// deleteZone Delete the Zone
func deleteZone(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to delete Zone %s", d.Get("zone"))
	configuration := d.Get("configuration").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err := objMgr.DeleteZone(configuration, view, zone)
	if err != nil {
		msg := fmt.Sprintf("Delete Zone %s failed: %s", zone, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Deletion of Zone complete")
	return nil
}

func checkServerExists(objMgr *utils.ObjectManager, configuration string, serverName string) bool {
	_, err := objMgr.GetServerByFQDN(configuration, serverName)
	if err != nil {
		log.Debugf("Getting server %s failed", serverName)
		return false
	}
	return true
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func updateDeployableProperty(deployable string, propertiesRaw string, isUpdateZone bool) (properties string, err error) {
	deployableValues := []string{"yes", "true", "1"}
	notDeployableValues := []string{"no", "false", "0", ""}

	isDeployable := contains(deployableValues, strings.ToLower(strings.Trim(deployable, " ")))
	isNotDeployable := contains(notDeployableValues, strings.ToLower(strings.Trim(deployable, " ")))

	properties = removeAttributeFromProperties("deployable", propertiesRaw)
	if isDeployable {
		properties = fmt.Sprintf("%s|deployable=true", properties)
	} else if isNotDeployable {
		if isUpdateZone {
			properties = fmt.Sprintf("%s|deployable=false", properties)
		}
	} else {
		err = fmt.Errorf("invalid deployable value (must be either 'true' or 'false'): '%s'", deployable)
	}
	return
}

func getRoleNameInRestApi(roleNameInTerraform string) string {
	// "NAME_IN_TERRAFORM": "NAME_IN_REST_API"
	roles := map[string]string{
		"FORWARDER":         "FORWARDER",
		"PRIMARY":           "MASTER",
		"PRIMARY_HIDDEN":    "MASTER_HIDDEN",
		"NONE":              "NONE",
		"RECURSION":         "RECURSION",
		"SECONDARY":         "SLAVE",
		"SECONDARY_STEALTH": "SLAVE_STEALTH",
		"STUB":              "STUB",
	}
	return roles[roleNameInTerraform]
}

func getRoleNameInTerraform(roleNameInRestApi string) string {
	// "NAME_IN_REST_API": "NAME_IN_TERRAFORM"
	roles := map[string]string{
		"FORWARDER":     "FORWARDER",
		"MASTER":        "PRIMARY",
		"MASTER_HIDDEN": "PRIMARY_HIDDEN",
		"NONE":          "NONE",
		"RECURSION":     "RECURSION",
		"SLAVE":         "SECONDARY",
		"SLAVE_STEALTH": "SECONDARY_STEALTH",
		"STUB":          "STUB",
	}
	return roles[roleNameInRestApi]
}

func validateServerRole(objMgr *utils.ObjectManager, configuration string, serverRole string) (role string, serverFQDN string, err error) {

	prop := strings.Split(serverRole, ",")
	if len(prop) != 2 {
		err = fmt.Errorf("invalid format server role: '%s'", serverRole)
		return
	}
	role = strings.Trim(prop[0], " ")
	serverFQDN = strings.Trim(prop[1], " ")

	roleNameInRestApi := getRoleNameInRestApi(strings.ToUpper(role))
	if roleNameInRestApi == "" {
		err = fmt.Errorf("invalid role type: '%s'", role)
		return
	}

	if serverFQDN == "" {
		err = fmt.Errorf("'server_fqdn' is a required property: '%s'", serverRole)
		return
	}

	if !checkServerExists(objMgr, configuration, serverFQDN) {
		err = fmt.Errorf("Server '%s' with role  '%s' doesn't exists", serverFQDN, role)
		return
	}

	role = roleNameInRestApi
	return
}

func prepareServerRoleData(objMgr *utils.ObjectManager, serverRolesRaw []interface{}, configuration string, view string, zone string) (map[string]string, map[string]string, error) {
	newServerRoles := make(map[string]string)
	currentServerRoles := make(map[string]string)

	for i, serverRole := range serverRolesRaw {
		svrRole, ok := serverRole.(string)
		if !ok {
			log.Warning("Cannot load server role item at: ", i)
		}
		svrRole = strings.TrimSpace(svrRole)
		if len(svrRole) > 0 {
			role, serverFQDN, err := validateServerRole(objMgr, configuration, svrRole)
			if err != nil {
				return newServerRoles, currentServerRoles, err
			}
			newServerRoles[serverFQDN] = role
		}

	}

	serverRoles, err := objMgr.GetDeploymentRoles(configuration, view, zone)
	if err != nil {
		err = fmt.Errorf("error get all deployment roles on the zone: %s", err)
		return newServerRoles, currentServerRoles, err
	}

	for _, serverRole := range serverRoles.ServerRoles {
		currentServerRoles[serverRole.ServerFQDN] = getRoleNameInTerraform(serverRole.Role)
	}

	return newServerRoles, currentServerRoles, err
}

func updateServerRoles(objMgr *utils.ObjectManager, currentServerRoles map[string]string, newServerRoles map[string]string, configuration string, view string, zone string) ([][]string, error) {
	trace := make([][]string, 0)

	for currentServerFQDN, currentRole := range currentServerRoles {
		_, ok := newServerRoles[currentServerFQDN]
		if !ok {
			_, err := objMgr.DeleteDeploymentRole(configuration, view, zone, currentServerFQDN)
			if err != nil {
				return trace, err
			}
			trace = append(trace, []string{"server_role", currentServerFQDN, currentRole, "append"})
		}
	}

	for newServerFQDN, newRole := range newServerRoles {
		currentRole, ok := currentServerRoles[newServerFQDN]
		if ok && !strings.EqualFold(currentRole, newRole) && (strings.EqualFold(currentRole, "PRIMARY") || strings.EqualFold(currentRole, "PRIMARY_HIDDEN")) {
			_, err := objMgr.UpdateDeploymentRole(configuration, view, zone, newServerFQDN, "dns", newRole, "", "")
			if err != nil {
				return trace, err
			}
			trace = append(trace, []string{"server_role", newServerFQDN, currentRole, "update"})
			delete(newServerRoles, newServerFQDN)
		}
	}

	for newServerFQDN, newRole := range newServerRoles {
		currentRole, ok := currentServerRoles[newServerFQDN]
		if ok {
			if !strings.EqualFold(currentRole, newRole) {
				_, err := objMgr.UpdateDeploymentRole(configuration, view, zone, newServerFQDN, "dns", newRole, "", "")
				if err != nil {
					return trace, err
				}
				trace = append(trace, []string{"server_role", newServerFQDN, currentRole, "update"})
			}
		} else {
			_, err := objMgr.CreateDeploymentRole(configuration, view, zone, newServerFQDN, "dns", newRole, "", "")
			if err != nil {
				return trace, err
			}
			trace = append(trace, []string{"server_role", newServerFQDN, newRole, "delete"})
		}
	}
	return trace, nil
}

func prepareDeploymentOptionData(currentDeploymentOptionsRaw interface{}, newDeploymentOptionsRaw interface{}) (map[string]string, map[string]string) {
	newDeploymentOptions := utils.ExpandStringMap(newDeploymentOptionsRaw)
	if newDeploymentOptions == nil {
		newDeploymentOptions = make(map[string]string)
	}
	currentDeploymentOptions := utils.ExpandStringMap(currentDeploymentOptionsRaw)
	if currentDeploymentOptions == nil {
		currentDeploymentOptions = make(map[string]string)
	}

	return newDeploymentOptions, currentDeploymentOptions
}

func updateDeploymentOptions(objMgr *utils.ObjectManager, currentDeploymentOptions map[string]string, newDeploymentOptions map[string]string, configuration string, view string, zone string, trace [][]string) ([][]string, error) {
	for currentOptionName, currentOptionValue := range currentDeploymentOptions {
		_, ok := newDeploymentOptions[currentOptionName]
		if !ok {
			_, err := objMgr.DeleteDeploymentOption(entities.DeploymentOption{
				Configuration: configuration,
				View:          view,
				Zone:          zone,
				Name:          currentOptionName,
				ServerID:      utils.DeploymentOptionAllServersID,
			})
			if err != nil {
				return trace, err
			}
			trace = append(trace, []string{"deployment_option", currentOptionName, currentOptionValue, "append"})
		}
	}

	for newOptionName, newOptionValue := range newDeploymentOptions {
		currentOptionValue, ok := currentDeploymentOptions[newOptionName]
		if ok {
			if !strings.EqualFold(currentOptionValue, newOptionValue) {
				_, err := objMgr.DeleteDeploymentOption(entities.DeploymentOption{
					Configuration: configuration,
					View:          view,
					Zone:          zone,
					Name:          newOptionName,
					Value:         currentOptionValue,
					ServerID:      utils.DeploymentOptionAllServersID,
				})
				if err != nil {
					return trace, err
				}
				trace = append(trace, []string{"deployment_option", newOptionName, currentOptionValue, "append"})
				_, err = objMgr.CreateDeploymentOption(entities.DeploymentOption{
					Configuration: configuration,
					View:          view,
					Zone:          zone,
					Name:          newOptionName,
					Value:         newOptionValue,
				})
				if err != nil {
					return trace, err
				}
				trace = append(trace, []string{"deployment_option", newOptionName, newOptionValue, "delete"})
			}
		} else {
			_, err := objMgr.CreateDeploymentOption(entities.DeploymentOption{
				Configuration: configuration,
				View:          view,
				Zone:          zone,
				Name:          newOptionName,
				Value:         newOptionValue,
			})
			if err != nil {
				return trace, err
			}
			trace = append(trace, []string{"deployment_option", newOptionName, newOptionValue, "delete"})
		}
	}

	return trace, nil
}

func rollBackData(objMgr *utils.ObjectManager, trace [][]string, configuration string, view string, zone string) (err error) {
	for len(trace) > 0 {
		item := trace[len(trace)-1]
		itemType, itemKey, itemValue, action := item[0], item[1], item[2], item[3]
		if itemType == "server_role" {
			if action == "append" {
				_, err = objMgr.CreateDeploymentRole(configuration, view, zone, itemKey, "dns", itemValue, "", "")
			} else if action == "delete" {
				_, err = objMgr.DeleteDeploymentRole(configuration, view, zone, itemKey)
			} else if action == "update" {
				_, err = objMgr.UpdateDeploymentRole(configuration, view, zone, itemKey, "dns", itemValue, "", "")
			}
		} else if itemType == "deployment_option" {
			if action == "append" {
				_, err = objMgr.CreateDeploymentOption(entities.DeploymentOption{
					Configuration: configuration,
					View:          view,
					Zone:          zone,
					Name:          itemKey,
					Value:         itemValue,
				})
			} else if action == "delete" {
				_, err = objMgr.DeleteDeploymentOption(entities.DeploymentOption{
					Configuration: configuration,
					View:          view,
					Zone:          zone,
					Name:          itemKey,
					ServerID:      utils.DeploymentOptionAllServersID,
				})
			}
		}
		if err != nil {
			return
		}
		trace = trace[:len(trace)-1]
	}
	return
}

func getDeploymentOptions(objMgr *utils.ObjectManager, configuration string, view string, zone string, optionNames []string) (map[string]string, error) {
	deploymentOptions := make(map[string]string, len(optionNames))
	if len(optionNames) == 0 {
		return deploymentOptions, nil
	}

	for _, optionName := range optionNames {
		deploymentOption, err := objMgr.GetDeploymentOption(entities.DeploymentOption{
			Configuration: configuration,
			View:          view,
			Zone:          zone,
			Name:          optionName,
			ServerID:      utils.DeploymentOptionLookupServerID,
		})
		if err != nil {
			return nil, err
		}
		deploymentOptions[optionName] = deploymentOption.Value
	}

	return deploymentOptions, nil
}
