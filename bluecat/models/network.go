// Copyright 2020 BlueCat Networks. All rights reserved

package models

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/entities"
)

const (
	// AllocateStatic Allocate the static IP Address
	AllocateStatic string = "MAKE_STATIC"
	// AllocateReserved Reserve the IP Address
	AllocateReserved string = "MAKE_RESERVED"
	// AllocateDHCPReserved Allocate the IP Address for DHCP
	AllocateDHCPReserved string = "MAKE_DHCP_RESERVED"
)

func getPath(configuration string) string {
	result := ""
	if len(configuration) > 0 {
		result = "/configurations/" + configuration
	}
	return result
}

func getIPPath(configuration string) string {
	result := ""
	if len(configuration) > 0 {
		result = "/configurations/" + configuration
	} else {
		result = "/ipv4_addresses"
	}
	return result
}

// NewBlock Initialize the new IPv4 Block to be added
func NewBlock(block entities.Block) *entities.Block {
	res := block
	res.SetObjectType("ipv4_blocks")

	path := getPath(res.Configuration)

	if len(block.ParentBlock) == 0 {
		res.SetSubPath(path)
	} else {
		res.SetSubPath(fmt.Sprintf("%s/ipv4_blocks/%s", path, block.ParentBlock))
	}
	return &res
}

// Block Initialize the IPv4 Block to be loaded, updated or deleted
func Block(block entities.Block) *entities.Block {
	res := block
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/ipv4_blocks/%s", getPath(res.Configuration), block.AddressCIDR()))
	return &res
}

// Network

// NewNetwork Initialize the new IPv4 Network to be added
func NewNetwork(network entities.Network) *entities.Network {
	res := network
	res.SetObjectType("create_network")
	res.SetSubPath(fmt.Sprintf("%s/ipv4_blocks/%s", getPath(res.Configuration), network.BlockAddr))

	return &res
}

// NewNextAvailableNetwork Initialize the new next available IPv4 Network to be added
func NewNextAvailableNetwork(network entities.Network) *entities.Network {
	res := network
	res.SetObjectType("get_next_network")
	res.SetSubPath(fmt.Sprintf("%s/ipv4_blocks/%s", getPath(res.Configuration), network.BlockAddr))

	return &res
}

// Network Initialize the IPv4 Network to be loaded, updated or deleted
func Network(network entities.Network) *entities.Network {
	res := network
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/ipv4_networks/%s", getPath(res.Configuration), network.CIDR))

	return &res
}

// NetworkByAllocatedId Initialize the IPv4 Network to be loaded by allocated id
func NetworkByAllocatedId(network entities.Network) *entities.Network {
	res := network
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/ipv4_blocks/%s/get_network_by_allocated_id/%s", getPath(res.Configuration), network.BlockAddr, network.AllocatedId))

	return &res
}

// IP Address

// GetNextIPAddress Initialize the new IPv4 Address for getting next available address
func GetNextIPAddress(ipAddr entities.IPAddress) *entities.IPAddress {
	res := ipAddr
	if len(ipAddr.Action) == 0 {
		res.Action = AllocateStatic
	}
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/ipv4_networks/%s/get_next_ip", getIPPath(res.Configuration), ipAddr.CIDR))

	return &res
}

// IPAddress Initialize the IPv4 Address
func IPAddress(ipAddr entities.IPAddress) *entities.IPAddress {
	res := ipAddr
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/ipv4_address/%s", getIPPath(res.Configuration), ipAddr.Address))

	return &res
}

// DHCP Range

// NewDHCPRange Initialize the new DHCP Range to be added
func NewDHCPRange(dhcpRange entities.DHCPRange) *entities.DHCPRange {
	res := dhcpRange
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/ipv4_networks/%s/dhcp_ranges", getPath(res.Configuration), dhcpRange.Network))

	return &res
}

// DHCPRange Initialize the DHCP Range to be loaded, updated or deleted
func DHCPRange(dhcpRange entities.DHCPRange) *entities.DHCPRange {
	res := dhcpRange
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/ipv4_networks/%s/start/%s/end/%s/dhcp_ranges", getPath(res.Configuration), dhcpRange.Network, dhcpRange.Start, dhcpRange.End))

	return &res
}
