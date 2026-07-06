// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceParseId(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected ip_address/cidr", id)
	}
	address := parts[0]
	cidr := parts[1]

	return address, cidr, nil
}

func getIpVersion(d *schema.ResourceData, address string) string {
	ip := net.ParseIP(address)
	ip_from_id := net.ParseIP(d.Id())
	ipVersion := d.Get("ip_version").(string)
	if ipVersion == "" {
		if ip.To4() != nil || ip_from_id.To4() != nil {
			ipVersion = "ipv4"
		} else {
			ipVersion = "ipv6"
		}
	}
	return ipVersion
}

// ResourceBlock The IPv4 Block
func ResourceBlock() *schema.Resource {
	return &schema.Resource{
		Create: createBlock,
		Read:   getBlock,
		Update: updateBlock,
		Delete: deleteBlock,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// Next-available mode resolves address/cidr during Create, so mark as computed at plan time.
			address := d.Get("address").(string)
			cidr := d.Get("cidr").(string)
			parentBlock := d.Get("parent_block").(string)
			size := d.Get("size").(string)
			if address == "" && cidr == "" && parentBlock != "" && size != "" {
				d.SetNewComputed("address")
				d.SetNewComputed("cidr")
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Configuration. Creating the Block in the default Configuration if doesn't specify",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Block name",
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IPv4 Block's address",
			},
			"cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Block prefix length",
			},
			"parent_block": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parent Block. Specified to creating the child Block. THe Block in CIDR format",
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The size of the block expressed in the power of 2",
			},
			"allocated_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The allocated id of the next available block",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != ""
				},
			},
			"deployment_options": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The deployment options for the block.",
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
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Block IP version: ipv4 or ipv6",
				Default:     "ipv4",
			},
		},
		Importer: &schema.ResourceImporter{
			State: resourceImporter,
		},
	}
}

// createIP4Block Create the new IPv4 Block
func createBlock(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to create Block %s", d.Get("address"))

	block := entities.Block{}
	block.InitBlock(d)

	objMgr := GetObjManager(m)
	if block.Address != "" || block.CIDR != "" {
		// Create specified block
		if block.Address == "" || block.CIDR == "" {
			msg := "'address' and 'cidr' are required to create specified block"
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		_, err := strconv.Atoi(block.CIDR)
		if err != nil {
			msg := fmt.Sprintf("Error converting the CIDR (%s): %s", block.CIDR, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		_, err = objMgr.CreateBlock(block)
		if err != nil {
			msg := fmt.Sprintf("Error creating Block (%s): %s", block.Address, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		log.Debugf("Completed to create Block %s", d.Get("address"))
	} else if block.IPVersion != entities.IPV6 {
		// Create next available block

		if block.ParentBlock == "" {
			msg := "'parent_block' is a required property to get next available block"
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		sizeNumber, err := strconv.Atoi(block.Size)
		if err != nil || !isPowerOfTwo(sizeNumber) {
			msg := "'size' is a required property and must be power of 2 to get next available block"
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		parentBlockAddress := block.ParentBlock
		parentBlockCIDR := "0"
		if strings.Contains(block.ParentBlock, "/") {
			parts := strings.SplitN(block.ParentBlock, "/", 2)
			parentBlockAddress = parts[0]
			parentBlockCIDR = parts[1]
		}

		parentBlockEntity, err := objMgr.GetBlock(block.Configuration, parentBlockAddress, parentBlockCIDR, block.IPVersion)
		if err != nil {
			msg := fmt.Sprintf("Failed to getting the IPv4 Block for (%s): %s", block.ParentBlock, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		block.ParentBlock = parentBlockEntity.AddressCIDR()

		_, ref, err := objMgr.CreateNextAvailableBlock(block)
		if err != nil {
			msg := fmt.Sprintf("Error creating next available block of Block(%s): %s", block.ParentBlock, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}

		log.Debugf("Successful to create next available Block of Block %s", d.Get("parent_block"))

		blockAddressCIDR := getObjectFieldValue("CIDR", ref)
		if blockAddressCIDR != "" && strings.Contains(blockAddressCIDR, "/") {
			addressCIDR := strings.SplitN(blockAddressCIDR, "/", 2)
			block.Address = addressCIDR[0]
			block.CIDR = addressCIDR[1]
		}
		if block.Address == "" || block.CIDR == "" {
			msg := fmt.Sprintf("Error parsing created next available Block response for (%s): %s", block.ParentBlock, ref)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		d.Set("address", block.Address)
		d.Set("cidr", block.CIDR)
	} else {
		msg := "'address' and 'cidr' are required to create IPv6 block"
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	err := utils.CreateDeploymentOptions(objMgr, entities.DeploymentOption{
		Configuration: block.Configuration,
		ResourceType:  "block",
		ResourceRef:   block.AddressCIDR(),
		IPVersion:     block.IPVersion,
	}, utils.ExpandStringMap(d.Get("deployment_options")))
	if err != nil {
		return fmt.Errorf("creating deployment options on Block (%s) failed: %w", block.AddressCIDR(), err)
	}
	return getBlock(d, m)
}

// getBlock Get the IPv4 Block
func getBlock(d *schema.ResourceData, m interface{}) error {
	var address, cidrStr string
	var err error
	if d.Id() != "" {
		address, cidrStr, err = resourceServiceParseId(d.Id())
	}

	if err != nil {
		return err
	}
	log.Debugf("Beginning to get Block %s", d.Get("address"))
	configuration := d.Get("configuration").(string)

	if address == "" {
		address = d.Get("address").(string)
	}
	ipVersion := getIpVersion(d, address)
	if cidrStr == "" {
		cidrStr = d.Get("cidr").(string)
	}
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	block, err := objMgr.GetBlock(configuration, address, cidrStr, ipVersion)

	if err != nil {
		if utils.IsNotFoundErr(err) {
			if d.Id() != "" {
				// If the record is missing remotely, remove from state so Terraform plans a create.
				log.Warnf("Block %q not found; removing from state to trigger recreation", d.Id())
				d.SetId("")
				return nil
			}
			// If we don't have an ID yet (e.g., during import resolution) surface the not-found
			return fmt.Errorf("Block %s not found: %w", cidrStr, err)
		}
		// Any other error is a real failure
		return fmt.Errorf("Getting Block %s failed: %w", cidrStr, err)
	}
	// --- Parse both server and config properties ---
	bamProps := utils.ParseProperties(block.Properties)
	cfgProps := utils.ParseProperties(d.Get("properties").(string))

	// --- Filter server properties using keys from config ---
	filteredProperties := utils.FilterProperties(bamProps, cfgProps)
	deploymentOptions, err := utils.ReadDeploymentOptions(objMgr, entities.DeploymentOption{
		Configuration: configuration,
		ResourceType:  "block",
		ResourceRef:   block.AddressCIDR(),
		IPVersion:     block.IPVersion,
	}, utils.ExpandStringMap(d.Get("deployment_options")))
	if err != nil {
		return fmt.Errorf("getting deployment options on Block %s failed: %w", block.AddressCIDR(), err)
	}

	d.Set("name", block.Name)
	d.Set("address", block.Address)
	d.Set("cidr", block.CIDR)
	d.Set("ip_version", block.IPVersion)
	d.Set("properties", utils.JoinProperties(filteredProperties))
	d.Set("deployment_options", utils.FlattenStringMap(deploymentOptions))
	d.SetId(block.AddressCIDR())
	log.Debugf("Completed getting Block %s", d.Get("address"))
	return nil
}

// updateBlock Update the existing IPv4 Block
func updateBlock(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to update Block %s", d.Get("address"))

	block := entities.Block{}
	address := d.Get("address").(string)
	d.Set("ip_version", getIpVersion(d, address))
	block.InitBlock(d)
	d.Set("ip_version", nil)

	_, err := strconv.Atoi(block.CIDR)
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", block.CIDR, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err = objMgr.UpdateBlock(block)
	if err != nil {
		msg := fmt.Sprintf("Error updating Block (%s): %s", block.Address, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	currentRaw, newRaw := d.GetChange("deployment_options")
	err = utils.UpdateDeploymentOptionsForTarget(objMgr, entities.DeploymentOption{
		Configuration: block.Configuration,
		ResourceType:  "block",
		ResourceRef:   block.AddressCIDR(),
		IPVersion:     block.IPVersion,
	}, currentRaw, newRaw)
	if err != nil {
		return fmt.Errorf("updating deployment options on Block (%s) failed: %w", block.AddressCIDR(), err)
	}
	log.Debugf("Completed to update Block %s", d.Get("address"))
	return getBlock(d, m)
}

// deleteBlock Delete the IPv4 Block
func deleteBlock(d *schema.ResourceData, m interface{}) error {
	log.Debugf("Beginning to Delete Block %s", d.Get("address"))
	configuration := d.Get("configuration").(string)
	address := d.Get("address").(string)
	cidrStr := d.Get("cidr").(string)
	ipVersion := d.Get("ip_version").(string)
	_, err := strconv.Atoi(cidrStr)
	if err != nil {
		msg := fmt.Sprintf("Error converting the CIDR (%s): %s", cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	ipVersion = getIpVersion(d, address)

	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector

	_, err = objMgr.DeleteBlock(configuration, address, cidrStr, ipVersion)
	if err != nil {
		msg := fmt.Sprintf("Delete Block %s/%s failed: %s", address, cidrStr, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	d.SetId("")
	log.Debugf("Deletion of Block complete ")
	return nil
}
