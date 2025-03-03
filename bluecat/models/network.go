// Copyright 2020 BlueCat Networks. All rights reserved

package models

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/entities"
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
func NewBlock(block entities.Block, ipVersion string) entities.Block {
	res := block
	// ipVersion should be ipv4 or ipv6
	res.SetObjectType(fmt.Sprintf("%s_blocks", ipVersion))

	path := getPath(res.Configuration)

	if len(block.ParentBlock) == 0 {
		res.SetSubPath(path)
	} else {
		res.SetSubPath(fmt.Sprintf("%s/%s_blocks/%s", path, ipVersion, block.ParentBlock))
	}
	return res
}

// IPBlock Initialize the IPv4/IPv6 Block to be loaded, updated or deleted
func IPBlock(block entities.Block) *entities.Block {
	block.SetObjectType("")
	block.SetSubPath(fmt.Sprintf("%s/%s_blocks/%s", getPath(block.Configuration), block.IPVersion, block.AddressCIDR()))
	return &block
}

// Network

// NewNetwork Initialize the new IPv4 Network to be added
func NewNetwork(network entities.Network) entities.Network {
	res := network
	res.SetObjectType("create_network")
	res.SetSubPath(fmt.Sprintf("%s/%s_blocks/%s", getPath(res.Configuration), network.IPVersion, network.BlockAddr))

	return res
}

// NewNextAvailableNetwork Initialize the new next available IPv4 Network to be added
func NewNextAvailableNetwork(network entities.Network) *entities.Network {
	res := network
	res.SetObjectType("get_next_network")
	res.SetSubPath(fmt.Sprintf("%s/%s_blocks/%s", getPath(res.Configuration), network.IPVersion, network.BlockAddr))

	return &res
}

// Network Initialize the IPv4/IP Network to be loaded, updated or deleted
func Network(network entities.Network) *entities.Network {
	res := network
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/%s_networks/%s", getPath(res.Configuration), network.IPVersion, network.CIDR))

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
		res.Action = entities.AllocateStatic
	}
	res.SetObjectType("")
	res.SetSubPath(
		fmt.Sprintf("%s/%s_networks/%s/get_next_ip", getIPPath(res.Configuration), ipAddr.IPVersion, ipAddr.CIDR),
	)

	return &res
}

// IPAddress Initialize the IPv4 Address
func IPAddress(ipAddr entities.IPAddress) *entities.IPAddress {
	res := ipAddr
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/%s_address/%s", getIPPath(res.Configuration), ipAddr.IPVersion, ipAddr.Address))

	return &res
}

// DHCP Range

// NewDHCPRange Initialize the new DHCP Range to be added
func NewDHCPRange(dhcpRange entities.DHCPRange) *entities.DHCPRange {
	res := dhcpRange
	res.SetObjectType("")
	res.SetSubPath(
		fmt.Sprintf("%s/%s_networks/%s/dhcp_ranges", getPath(res.Configuration), dhcpRange.IPVersion, dhcpRange.Network),
	)

	return &res
}

// DHCPRange Initialize the DHCP Range to be loaded, updated or deleted
func DHCPRange(dhcpRange entities.DHCPRange) *entities.DHCPRange {
	res := dhcpRange
	res.SetObjectType("")
	if dhcpRange.IPVersion == entities.IPV4 {
		res.SetSubPath(
			fmt.Sprintf(
				"%s/ipv4_networks/%s/start/%s/end/%s/dhcp_ranges",
				getPath(res.Configuration),
				dhcpRange.Network,
				dhcpRange.Start,
				dhcpRange.End,
			),
		)
	} else if dhcpRange.IPVersion == entities.IPV6 {
		res.SetSubPath(
			fmt.Sprintf(
				"%s/ipv6_networks/%s/dhcp_range/start/%s/end/%s",
				getPath(res.Configuration),
				dhcpRange.Network,
				dhcpRange.Start,
				dhcpRange.End,
			),
		)
	}

	return &res
}

// NewView Initialize the new View to be added
func NewView(view *entities.View) *entities.View {
	path := getPath(view.Configuration)
	view.SetSubPath(fmt.Sprintf("%s/views", path))
	return view
}

// View Initialize the View to be loaded, updated or deleted
func View(view entities.View) *entities.View {
	res := view
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/views/%s", getPath(res.Configuration), view.Name))
	return &res
}
