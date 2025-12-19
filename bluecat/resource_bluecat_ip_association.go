// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceIPAssociation The IP Association
func ResourceIPAssociation() *schema.Resource {
	return &schema.Resource{
		Create: createIPAssociation,
		Read:   getIPAssociation,
		Update: updateIPAssociation,
		Delete: deleteIPAssociation,

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Associate the IP address/Host record in the default Configuration if doesn't specify",
			},
			"view": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The view which contains the details of the zone. If not provided, uses the default view",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Zone in which you want to update a host record. If not provided, the absolute name must be FQDN ones",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Host record. Must be FQDN if the Zone is not provided",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					zone := d.Get("zone").(string)
					return checkDiffName(old, new, zone)
				},
			},
			"network": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Network address in CIDR format",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP address",
			},
			"mac_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The MAC address",
			},
			"properties": {
				Type:     schema.TypeString,
				Optional: true,
				StateFunc: func(v interface{}) string {
					return utils.JoinProperties(utils.ParseProperties(v.(string)))
				},
				DiffSuppressFunc: suppressWhenRemoteHasSuperset,
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP Address version: ipv4 or ipv6",
			},
		},
	}
}

// createIPAssociation Associate the IP address/Host record
func createIPAssociation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to associate IP address %s", d.Get("ip_address"))
	err := updateAllocatedResource(d, m)
	if err != nil {
		return err
	}
	log.Debugf("Completed to associate IP address %s", d.Get("ip_address"))
	return getIPAssociation(d, m)
}

// getIPAssociation Get the allocated IP address/Host info
func getIPAssociation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to get IP address: %s", d.Get("ip_address").(string))
	err := getIPAllocation(d, m)
	if err != nil {
		return err
	}
	log.Debugf("Completed reading IP address %s", d.Get("ip_address"))
	return nil
}

// updateIPAssociation Update the association
func updateIPAssociation(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update the association for the IP address %s", d.Get("ip_address"))
	err := updateAllocatedResource(d, m)
	if err != nil {
		return err
	}
	log.Debugf("Completed to update association %s", d.Get("ip_address"))
	return getIPAssociation(d, m)
}

// deleteIPAssociation Delete the association IP address/Host record
func deleteIPAssociation(d *schema.ResourceData, m interface{}) error {

	objMgr := GetObjManager(m)

	address := entities.IPAddress{}
	if !address.InitIPAddress(d) {
		log.Error(address.InitError)
		return fmt.Errorf(address.InitError)
	}
	//address.Properties = utils.RemoveImmutableProperties(address.Properties, []string{"parentId", "parentType", "addresses", "addressesIds"})

	log.Debugf("Beginning to release an association for the IP address %s", address.Address)
	//configuration := d.Get("configuration").(string)
	//ipAddress := d.Get("ip_address").(string)
	view := d.Get("view").(string)
	zone := d.Get("zone").(string)
	//name := d.Get("name").(string)
	//ipVersion := d.Get("ip_version").(string)

	fqdnName := address.Name
	if len(zone) > 0 {
		fqdnName = getFQDN(address.Name, zone)
	}

	log.Debugf("Getting host record %s", fqdnName)
	hostRecord, err := objMgr.GetHostRecord(address.Configuration, view, fqdnName)
	//hostRecord.Properties = utils.RemoveImmutableProperties(hostRecord.Properties, []string{"parentId", "parentType"})
	if err != nil {
		msg := fmt.Sprintf("The Host record %s not found: %s", fqdnName, err)
		log.Debug(msg)
	} else {
		// Checking for existing linked IP address
		properties := hostRecord.Properties

		currentAssociateIPs := utils.GetPropertyValue("addresses", properties)

		if strings.Contains(currentAssociateIPs, address.Address) && len(strings.Split(currentAssociateIPs, ",")) > 1 {
			TTL := utils.GetPropertyValue("ttl", hostRecord.Properties)
			rrTTL, err := strconv.Atoi(TTL)
			if err != nil {
				msg := fmt.Sprintf("Convert Host record TTL %s failed: %s", TTL, err)
				log.Debug(msg)
				rrTTL = -1
			}
			log.Debugf("Removing the IP %s from the Host record %s", address.Address, fqdnName)
			associateIPs := removeIPFromList(currentAssociateIPs, address.Address)
			properties = removeAttributeFromProperties("addresses", properties)
			properties = fmt.Sprintf("%s|addresses=%s", properties, associateIPs)
			log.Debugf("Association destroy properties: %s", properties)
			_, err = objMgr.UpdateHostRecord(address.Configuration, view, zone, fqdnName, associateIPs, rrTTL, properties)
			if err != nil {
				msg := fmt.Sprintf("Error updating Host record %s: %s", fqdnName, err)
				log.Debug(msg)
				return fmt.Errorf(msg)
			}
		}
	}

	address.Mac = "00:00:00:00:00:00"
	_, err = objMgr.SetMACAddress(address)
	if err != nil {
		msg := fmt.Sprintf("Releasing the IP address %s failed: %s", address.Address, err)
		log.Debug(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Completed to release an association for the IP address %s", address.Address)
	return nil
}

func removeIPFromList(ips, ip string) (val string) {
	result := ""
	ipList := strings.Split(ips, ",")
	for i := 0; i < len(ipList); i++ {
		if ipList[i] != ip {
			result = fmt.Sprintf("%s,%s", result, ipList[i])
		}
	}
	return result[1:]
}
