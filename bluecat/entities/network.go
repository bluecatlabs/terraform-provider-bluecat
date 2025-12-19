// Copyright 2020 BlueCat Networks. All rights reserved

package entities

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const IPV4 = "ipv4"
const IPV6 = "ipv6"

// Block IPv4/IPv6 Block entity
type Block struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	ParentBlock   string `json:"-"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	CIDR          string `json:"cidr_notation"`
	Properties    string `json:"properties,omitempty"`
	BlockId       int    `json:"id,omitempty"`
	IPVersion     string `json:"ip_version,omitempty"`
}

func (block *Block) InitBlock(blockMap *schema.ResourceData) {
	block.Configuration = blockMap.Get("configuration").(string)
	block.Name = blockMap.Get("name").(string)
	block.Address = blockMap.Get("address").(string)
	block.ParentBlock = blockMap.Get("parent_block").(string)
	block.CIDR = blockMap.Get("cidr").(string)
	block.Properties = blockMap.Get("properties").(string)
	block.IPVersion = blockMap.Get("ip_version").(string)
}

// GetIPv6BlockFromPropsPrefix func will get the block cidr notation from properties prefix attribute
func (block *Block) GetIPv6BlockFromPropsPrefix() string {
	propertiesMap := GetPropertiesFromString(block.Properties)
	return propertiesMap["prefix"]
}

// AddressCIDR Get the Block address in CIDR format
func (block *Block) AddressCIDR() string {
	return fmt.Sprintf("%s/%s", block.Address, block.CIDR)
}

func getResourceIPVersion(resourceMap *schema.ResourceData) string {
	ipVersion := resourceMap.Get("ip_version")
	if ipVersion == "" || ipVersion == IPV4 {
		// setting default value
		return IPV4
	} else if ipVersion == IPV6 {
		return IPV6
	} else {
		return ""
	}
}

// Network IPv4 Network entity
type Network struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	BlockAddr     string `json:"-"`
	Name          string `json:"name"`
	CIDR          string `json:"cidr"`
	Gateway       string `json:"gateway"`
	Properties    string `json:"properties"`
	Template      string `json:"template"`
	ParentBlock   string `json:"parent_block"`
	Size          string `json:"size"`
	NetWorkId     int    `json:"id"`
	AllocatedId   string `json:"allocatedId"`
	IPVersion     string `json:"ip_version" default:"ipv4"`
	InitError     string `json:"nil"`
}

func (network *Network) InitNetwork(networkMap *schema.ResourceData) bool {

	// Using reflection we populate network struct from map based on json tag which represents key in the map.
	structType := reflect.TypeOf(*network)
	structValue := reflect.ValueOf(network).Elem()
	// Iterate over the fields of the struct
	for i := 0; i < structType.NumField(); i++ {

		field := structType.Field(i)
		fieldName := field.Name
		jsonTag := field.Tag.Get("json")

		// If the field exists in the ResourceData, set its value
		if value, ok := networkMap.GetOk(jsonTag); ok {
			// Set the field value
			fieldValue := structValue.FieldByName(fieldName)
			if fieldValue.IsValid() && fieldValue.CanSet() {
				fieldValue.SetString(value.(string))
			}
		}

	}

	network.IPVersion = getResourceIPVersion(networkMap)
	if network.IPVersion == "" {
		network.InitError = fmt.Sprintf(
			"Unknown ip_version. If you want to create IPv6 resource please set ip_version = 'ipv6'",
		)
		return false
	}

	network.Configuration = networkMap.Get("configuration").(string)

	return true
}

// IPAddress The IPv4 Address entity
type IPAddress struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	Template      string `json:"template,omitempty"`
	Action        string `json:"action,omitempty"`
	CIDR          string `json:"network,omitempty"`
	Address       string `json:"address,omitempty"`
	Mac           string `json:"mac_address,omitempty"`
	Name          string `json:"name,omitempty"`
	Properties    string `json:"properties,omitempty"`
	IPVersion     string `json:"ip_version,omitempty"`

	InitError string `json:"nil"`
}

func (ipAddress *IPAddress) InitIPAddress(ipAddressMap *schema.ResourceData) bool {

	ipAddress.IPVersion = getResourceIPVersion(ipAddressMap)
	if ipAddress.IPVersion == "" {
		ipAddress.InitError = fmt.Sprintf(
			"Unknown ip_version. If you want to create IPv6 resource please set ip_version = 'ipv6'",
		)
		return false
	}

	ipAddress.Configuration = ipAddressMap.Get("configuration").(string)
	ipAddress.Name = ipAddressMap.Get("name").(string)
	ipAddress.Address = ipAddressMap.Get("ip_address").(string)
	ipAddress.Mac = ipAddressMap.Get("mac_address").(string)
	ipAddress.Properties = ipAddressMap.Get("properties").(string)

	if ipAddressMap.Get("action") != nil {
		ipAddress.Action = ipAddressMap.Get("action").(string)
	}
	if ipAddressMap.Get("template") != nil {
		ipAddress.Template = ipAddressMap.Get("template").(string)
	}
	return true
}

// SetAction will convert status that we get from the response of the GET method into valid value for the PATCH
func (ipAddress *IPAddress) SetAction() {
	// convert action to be valid for the REST-API when we update IP Address
	switch ipAddress.Action {
	case "STATIC":
		ipAddress.Action = AllocateStatic
	case "RESERVED":
		ipAddress.Action = AllocateReserved
	case "DHCP_RESERVED":
		ipAddress.Action = AllocateDHCPReserved
	}
}

type DHCPRange struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	Template      string `json:"template,omitempty"`
	Network       string `json:"network,omitempty"`
	Start         string `json:"start,omitempty"`
	End           string `json:"end,omitempty"`
	Name          string `json:"name,omitempty"`
	Properties    string `json:"properties,omitempty"`
	IPVersion     string `json:"ip_version" default:"ipv4"`

	InitError string `json:"nil"`
}

func (dhcpRange *DHCPRange) InitRange(rangeMap *schema.ResourceData) bool {

	dhcpRange.IPVersion = getResourceIPVersion(rangeMap)
	if dhcpRange.IPVersion == "" {
		dhcpRange.InitError = fmt.Sprintf(
			"Unknown ip_version. If you want to create IPv6 resource please set ip_version = 'ipv6'",
		)
		return false
	}

	dhcpRange.Configuration = rangeMap.Get("configuration").(string)
	dhcpRange.Network = rangeMap.Get("network").(string)
	dhcpRange.Start = rangeMap.Get("start").(string)
	dhcpRange.End = rangeMap.Get("end").(string)
	dhcpRange.Name = rangeMap.Get("name").(string)
	dhcpRange.Properties = rangeMap.Get("properties").(string)
	dhcpRange.Template = rangeMap.Get("template").(string)

	return true
}

// View entity
type View struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	Name          string `json:"name"`
	Properties    string `json:"properties,omitempty"`
}

// TODO: If new function works, replace this with that function
func GetPropertiesFromString(input string) map[string]string {
	result := make(map[string]string)
	properties := strings.Split(input, "|")
	for _, prop := range properties {
		parts := strings.Split(prop, "=")
		if len(parts) == 2 {
			value := strings.Trim(parts[1], "'")
			result[strings.TrimSpace(parts[0])] = value
		}
	}
	return result
}
