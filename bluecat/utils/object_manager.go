// Copyright 2020 BlueCat Networks. All rights reserved

package utils

import (
	"encoding/json"
	"fmt"
	"terraform-provider-bluecat/bluecat/entities"
	"terraform-provider-bluecat/bluecat/models"
)

// ObjectManager The BlueCat object manager
type ObjectManager struct {
	Connector BCConnector
}

// Host record

// CreateHostRecord Create the Host record
func (objMgr *ObjectManager) CreateHostRecord(configuration string, view string, zone string, absoluteName string, ip4Address string, ttl int, properties string) (*entities.HostRecord, error) {

	hostRecord := models.NewHostRecord(entities.HostRecord{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		IP4Address:    ip4Address,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
	})

	_, err := objMgr.Connector.CreateObject(hostRecord)
	return hostRecord, err
}

// GetHostRecord Get the Host record
func (objMgr *ObjectManager) GetHostRecord(configuration string, view string, absoluteName string) (*entities.HostRecord, error) {

	hostRecord := models.HostRecord(entities.HostRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	err := objMgr.Connector.GetObject(hostRecord, &hostRecord)
	return hostRecord, err
}

// UpdateHostRecord Update the Host record
func (objMgr *ObjectManager) UpdateHostRecord(configuration string, view string, zone string, absoluteName string, ip4Address string, ttl int, properties string) (*entities.HostRecord, error) {

	hostRecord := models.HostRecord(entities.HostRecord{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		IP4Address:    ip4Address,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
	})

	err := objMgr.Connector.UpdateObject(hostRecord, &hostRecord)
	return hostRecord, err
}

// DeleteHostRecord Delete the Host record
func (objMgr *ObjectManager) DeleteHostRecord(configuration string, view string, absoluteName string) (string, error) {

	hostRecord := models.HostRecord(entities.HostRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	return objMgr.Connector.DeleteObject(hostRecord)
}

// CNAME record

// CreateCNAMERecord Create the CNAME record
func (objMgr *ObjectManager) CreateCNAMERecord(configuration string, view string, zone string, absoluteName string, linkedRecord string, ttl int, properties string) (*entities.CNAMERecord, error) {

	cnameRecord := models.NewCNAMERecord(entities.CNAMERecord{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		LinkedRecord:  linkedRecord,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
	})

	_, err := objMgr.Connector.CreateObject(cnameRecord)
	return cnameRecord, err
}

// GetCNAMERecord Get the CNAME record
func (objMgr *ObjectManager) GetCNAMERecord(configuration string, view string, absoluteName string) (*entities.CNAMERecord, error) {

	cnameRecord := models.CNAMERecord(entities.CNAMERecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	err := objMgr.Connector.GetObject(cnameRecord, &cnameRecord)
	return cnameRecord, err
}

// UpdateCNAMERecord Update the CNAME record
func (objMgr *ObjectManager) UpdateCNAMERecord(configuration string, view string, zone string, absoluteName string, linkedRecord string, ttl int, properties string) (*entities.CNAMERecord, error) {

	cnameRecord := models.CNAMERecord(entities.CNAMERecord{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		LinkedRecord:  linkedRecord,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
	})

	err := objMgr.Connector.UpdateObject(cnameRecord, &cnameRecord)
	return cnameRecord, err
}

// DeleteCNAMERecord Delete the CNAME record
func (objMgr *ObjectManager) DeleteCNAMERecord(configuration string, view string, absoluteName string) (string, error) {

	cnameRecord := models.CNAMERecord(entities.CNAMERecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	return objMgr.Connector.DeleteObject(cnameRecord)
}

// Configuration

// CreateConfiguration Create a new Configuration
func (objMgr *ObjectManager) CreateConfiguration(name string, properties string) (*entities.Configuration, error) {

	configuration := models.NewConfiguration(entities.Configuration{
		Name:       name,
		Properties: properties,
	})

	_, err := objMgr.Connector.CreateObject(configuration)
	return configuration, err
}

// GetConfiguration Get the Configuration info
func (objMgr *ObjectManager) GetConfiguration(name string) (*entities.Configuration, error) {

	configuration := models.Configuration(entities.Configuration{
		Name: name,
	})

	err := objMgr.Connector.GetObject(configuration, &configuration)
	return configuration, err
}

// UpdateConfiguration Update the Configuration info
func (objMgr *ObjectManager) UpdateConfiguration(name string, properties string) (*entities.Configuration, error) {

	configuration := models.Configuration(entities.Configuration{
		Name:       name,
		Properties: properties,
	})

	err := objMgr.Connector.UpdateObject(configuration, &configuration)
	return configuration, err
}

// DeleteConfiguration Delete the configuration
func (objMgr *ObjectManager) DeleteConfiguration(name string) (string, error) {

	configuration := models.Configuration(entities.Configuration{
		Name: name,
	})

	return objMgr.Connector.DeleteObject(configuration)
}

// Block

// CreateBlock Create a new Block
func (objMgr *ObjectManager) CreateBlock(configuration string, name string, address string, cidr string, parentBlock string, properties string) (*entities.Block, error) {

	block := models.NewBlock(entities.Block{
		Configuration: configuration,
		Name:          name,
		Address:       address,
		CIDR:          cidr,
		ParentBlock:   parentBlock,
		Properties:    properties,
	})

	_, err := objMgr.Connector.CreateObject(block)
	return block, err
}

// GetBlock Get the Block info
func (objMgr *ObjectManager) GetBlock(configuration string, address string, cidr string) (*entities.Block, error) {

	block := models.Block(entities.Block{
		Configuration: configuration,
		Address:       address,
		CIDR:          cidr,
	})

	err := objMgr.Connector.GetObject(block, &block)
	return block, err
}

// UpdateBlock Update the Block info
func (objMgr *ObjectManager) UpdateBlock(configuration string, name string, address string, cidr string, parentBlock string, properties string) (*entities.Block, error) {

	block := models.Block(entities.Block{
		Configuration: configuration,
		Name:          name,
		Address:       address,
		CIDR:          cidr,
		ParentBlock:   parentBlock,
		Properties:    properties,
	})

	err := objMgr.Connector.UpdateObject(block, &block)
	return block, err
}

// DeleteBlock Delete the Block
func (objMgr *ObjectManager) DeleteBlock(configuration string, address string, cidr string) (string, error) {

	block := models.Block(entities.Block{
		Configuration: configuration,
		Address:       address,
		CIDR:          cidr,
	})

	return objMgr.Connector.DeleteObject(block)
}

// Network

func generateNetworkProperties(props string, gateway string, allocatedId string) string {
	result := props
	if len(gateway) > 0 {
		result = fmt.Sprintf("%s|gateway=%s", result, gateway)
	}
	if len(allocatedId) > 0 {
		result = fmt.Sprintf("%s|allocatedId=%s", result, allocatedId)
	}
	return result
}

// CreateNetwork Create a new Network
func (objMgr *ObjectManager) CreateNetwork(configuration string, block string, name string, cidr string, gateway string, properties string, template string) (*entities.Network, error) {

	network := models.NewNetwork(entities.Network{
		Configuration: configuration,
		BlockAddr:     block,
		Name:          name,
		CIDR:          cidr,
		Gateway:       gateway,
		Properties:    generateNetworkProperties(properties, gateway, ""),
		Template:      template,
	})
	_, err := objMgr.Connector.CreateObject(network)
	return network, err
}

// CreateNextAvailableNetwork Create a next available Network
func (objMgr *ObjectManager) CreateNextAvailableNetwork(configuration string, block string, name string, gateway string, properties string, template string, size string, allocatedId string) (*entities.Network, string, error) {

	network := models.NewNextAvailableNetwork(entities.Network{
		Configuration: configuration,
		BlockAddr:     block,
		Name:          name,
		Gateway:       gateway,
		Properties:    generateNetworkProperties(properties, gateway, allocatedId),
		Template:      template,
		Size:          size,
		AllocatedId:   allocatedId,
	})
	ref, err := objMgr.Connector.CreateObject(network)
	return network, ref, err
}

// GetNetwork Get the Network info
func (objMgr *ObjectManager) GetNetwork(configuration string, cidr string) (*entities.Network, error) {

	network := models.Network(entities.Network{
		Configuration: configuration,
		CIDR:          cidr,
	})

	err := objMgr.Connector.GetObject(network, &network)
	return network, err
}

// GetNetworkByAllocatedId Get the Network info by allocated id
func (objMgr *ObjectManager) GetNetworkByAllocatedId(configuration string, block string, allocatedId string) (*entities.Network, error) {

	network := models.Network(entities.Network{
		Configuration: configuration,
		BlockAddr:     block,
		AllocatedId:   allocatedId,
	})

	err := objMgr.Connector.GetObject(network, &network)
	return network, err
}

// UpdateNetwork Update the Network info
func (objMgr *ObjectManager) UpdateNetwork(configuration string, name string, cidr string, gateway string, properties string) (*entities.Network, error) {

	network := models.Network(entities.Network{
		Configuration: configuration,
		Name:          name,
		CIDR:          cidr,
		Properties:    generateNetworkProperties(properties, gateway, ""),
	})

	err := objMgr.Connector.UpdateObject(network, &network)
	return network, err
}

// DeleteNetwork Delete the Network
func (objMgr *ObjectManager) DeleteNetwork(configuration string, cidr string) (string, error) {

	network := models.Network(entities.Network{
		Configuration: configuration,
		CIDR:          cidr,
	})

	return objMgr.Connector.DeleteObject(network)
}

// DHCP Range

// CreateDHCPRange Create a new DHCP Range
func (objMgr *ObjectManager) CreateDHCPRange(configuration string, template string, network string, start string, end string, properties string) (*entities.DHCPRange, error) {

	dhcpRange := models.NewDHCPRange(entities.DHCPRange{
		Configuration: configuration,
		Template:      template,
		Network:       network,
		Start:         start,
		End:           end,
		Properties:    properties,
	})
	_, err := objMgr.Connector.CreateObject(dhcpRange)
	return dhcpRange, err
}

// GetDHCPRange Get the DHCP Range info
func (objMgr *ObjectManager) GetDHCPRange(configuration string, network string, start string, end string) (*entities.DHCPRange, error) {

	dhcpRange := models.DHCPRange(entities.DHCPRange{
		Configuration: configuration,
		Network:       network,
		Start:         start,
		End:           end,
	})

	err := objMgr.Connector.GetObject(dhcpRange, &dhcpRange)
	return dhcpRange, err
}

// GetDeploymentRoles Get all Deployment role on the Zone
func (objMgr *ObjectManager) GetDeploymentRoles(configuration string, view string, zone string) (*entities.DeploymentRoles, error) {

	deploymentRoles := models.GetDeploymentRoles(entities.DeploymentRoles{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
	})

	err := objMgr.Connector.GetObject(deploymentRoles, &deploymentRoles)
	return deploymentRoles, err
}

// UpdateDHCPRange Update the DHCP Range info
func (objMgr *ObjectManager) UpdateDHCPRange(configuration string, template string, network string, start string, end string, properties string) (*entities.DHCPRange, error) {

	dhcpRange := models.DHCPRange(entities.DHCPRange{
		Configuration: configuration,
		Template:      template,
		Network:       network,
		Start:         start,
		End:           end,
		Properties:    properties,
	})

	err := objMgr.Connector.UpdateObject(dhcpRange, &dhcpRange)
	return dhcpRange, err
}

// DeleteDHCPRange Delete the DHCP Range
func (objMgr *ObjectManager) DeleteDHCPRange(configuration string, network string, start string, end string) (string, error) {

	dhcpRange := models.DHCPRange(entities.DHCPRange{
		Configuration: configuration,
		Network:       network,
		Start:         start,
		End:           end,
	})

	return objMgr.Connector.DeleteObject(dhcpRange)
}

// IP

// ReserveIPAddress Create the new IP address for later use
func (objMgr *ObjectManager) ReserveIPAddress(configuration string, network string) (*entities.IPAddress, error) {
	return objMgr.CreateIPAddress(configuration, network, "", "", "", models.AllocateReserved, "", "")
}

// CreateStaticIP Create the new static IP address
func (objMgr *ObjectManager) CreateStaticIP(configuration string, network string, address string, macAddress string, name string, properties string) (*entities.IPAddress, error) {
	return objMgr.CreateIPAddress(configuration, network, address, macAddress, name, models.AllocateStatic, properties, "")
}

// createIPAddress Create the new IP address. Allocate the next available on the network if IP address is not provided
func (objMgr *ObjectManager) CreateIPAddress(configuration string, cidr string, address string, macAddress string, name string, addrType string, properties string, template string) (*entities.IPAddress, error) {
	if len(addrType) == 0 {
		addrType = models.AllocateStatic
	}

	addrEntity := entities.IPAddress{
		Configuration: configuration,
		CIDR:          cidr,
		Name:          name,
		Address:       address,
		Mac:           macAddress,
		Action:        addrType,
		Properties:    properties,
		Template:      template,
	}

	ipAddr := new(entities.IPAddress)
	if len(address) > 0 {
		ipAddr = models.IPAddress(addrEntity)
	} else {
		ipAddr = models.GetNextIPAddress(addrEntity)
		log.Debugf("Requesting the new IP address in the network %s", cidr)
	}
	res, err := objMgr.Connector.CreateObject(ipAddr)
	if err == nil {
		err = json.Unmarshal([]byte(res), &ipAddr)
		if err == nil {
			log.Debugf("Failed to decode the IP object %s", err)
		}
	}
	return ipAddr, err
}

// GetIPAddress Get the IP Address info
func (objMgr *ObjectManager) GetIPAddress(configuration string, address string) (*entities.IPAddress, error) {

	ipAddr := models.IPAddress(entities.IPAddress{
		Configuration: configuration,
		Address:       address,
	})

	err := objMgr.Connector.GetObject(ipAddr, &ipAddr)
	return ipAddr, err
}

// SetMACAddress Update the MAC address for the existing IP address
func (objMgr *ObjectManager) SetMACAddress(configuration string, address string, macAddress string) (*entities.IPAddress, error) {
	ipAddr := models.IPAddress(entities.IPAddress{
		Configuration: configuration,
		Address:       address,
		Mac:           macAddress,
	})
	err := objMgr.Connector.UpdateObject(ipAddr, &ipAddr)
	return ipAddr, err
}

// UpdateIPAddress Update the IP address info
func (objMgr *ObjectManager) UpdateIPAddress(configuration string, address string, macAddress string, name string, addrType string, properties string) (*entities.IPAddress, error) {
	ipAddr := models.IPAddress(entities.IPAddress{
		Configuration: configuration,
		Name:          name,
		Address:       address,
		Mac:           macAddress,
		Action:        addrType,
		Properties:    properties,
	})
	err := objMgr.Connector.UpdateObject(ipAddr, &ipAddr)
	return ipAddr, err
}

// DeleteIPAddress Delete the existing IP address
func (objMgr *ObjectManager) DeleteIPAddress(configuration string, address string) (string, error) {
	ipAddr := models.IPAddress(entities.IPAddress{
		Configuration: configuration,
		Address:       address,
	})
	return objMgr.Connector.DeleteObject(ipAddr)
}

// CreateTXTRecord Create the TXT record
func (objMgr *ObjectManager) CreateTXTRecord(configuration string, view string, zone string, absoluteName string, text string, ttl int, properties string) (*entities.TXTRecord, error) {

	txtRecord := models.NewTXTRecord(entities.TXTRecord{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		Text:          text,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
	})

	_, err := objMgr.Connector.CreateObject(txtRecord)
	return txtRecord, err
}

// GetTXTRecord Get the TXT record
func (objMgr *ObjectManager) GetTXTRecord(configuration string, view string, absoluteName string) (*entities.TXTRecord, error) {

	txtRecord := models.TXTRecord(entities.TXTRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	err := objMgr.Connector.GetObject(txtRecord, &txtRecord)
	return txtRecord, err
}

// UpdateTXTRecord Update the TXT record
func (objMgr *ObjectManager) UpdateTXTRecord(configuration string, view string, zone string, absoluteName string, text string, ttl int, properties string) (*entities.TXTRecord, error) {

	txtRecord := models.TXTRecord(entities.TXTRecord{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		Text:          text,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
	})

	err := objMgr.Connector.UpdateObject(txtRecord, &txtRecord)
	return txtRecord, err
}

// DeleteTXTRecord Delete the TXT record
func (objMgr *ObjectManager) DeleteTXTRecord(configuration string, view string, absoluteName string) (string, error) {

	txtRecord := models.TXTRecord(entities.TXTRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	return objMgr.Connector.DeleteObject(txtRecord)
}

// CreateGenericRecord Create the Generic record
func (objMgr *ObjectManager) CreateGenericRecord(configuration string, view string, zone string, typerr string, absoluteName string, data string, ttl int, properties string) (*entities.GenericRecord, error) {

	genericRecord := models.NewGenericRecord(entities.GenericRecord{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		TypeRR:        typerr,
		Data:          data,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
	})

	_, err := objMgr.Connector.CreateObject(genericRecord)
	return genericRecord, err
}

// GetGenericRecord Get the Generic record
func (objMgr *ObjectManager) GetGenericRecord(configuration string, view string, absoluteName string) (*entities.GenericRecord, error) {

	genericRecord := models.GenericRecord(entities.GenericRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	err := objMgr.Connector.GetObject(genericRecord, &genericRecord)
	return genericRecord, err
}

// UpdateGenericRecord Update the Generic record
func (objMgr *ObjectManager) UpdateGenericRecord(configuration string, view string, zone string, typerr string, absoluteName string, data string, ttl int, properties string) (*entities.GenericRecord, error) {

	genericRecord := models.GenericRecord(entities.GenericRecord{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		TypeRR:        typerr,
		Data:          data,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
	})

	err := objMgr.Connector.UpdateObject(genericRecord, &genericRecord)
	return genericRecord, err
}

// DeleteGenericRecord Delete the Generic record
func (objMgr *ObjectManager) DeleteGenericRecord(configuration string, view string, absoluteName string) (string, error) {

	genericRecord := models.GenericRecord(entities.GenericRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	return objMgr.Connector.DeleteObject(genericRecord)
}

// Zone

// CreateZone Create a new Zone
func (objMgr *ObjectManager) CreateZone(configuration string, view string, zone string, properties string) (*entities.Zone, error) {

	zoneObj := models.NewZone(entities.Zone{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		Properties:    properties,
	})
	_, err := objMgr.Connector.CreateObject(zoneObj)
	return zoneObj, err
}

// GetZone Get the Zone info
func (objMgr *ObjectManager) GetZone(configuration string, view string, zone string) (*entities.Zone, error) {

	zoneObj := models.Zone(entities.Zone{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
	})

	err := objMgr.Connector.GetObject(zoneObj, &zoneObj)
	return zoneObj, err
}

// UpdateZone Update the Zone info
func (objMgr *ObjectManager) UpdateZone(configuration string, view string, zone string, properties string) (*entities.Zone, error) {

	zoneObj := models.Zone(entities.Zone{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		Properties:    properties,
	})

	err := objMgr.Connector.UpdateObject(zoneObj, &zoneObj)
	return zoneObj, err
}

// DeleteZone Delete the Zone
func (objMgr *ObjectManager) DeleteZone(configuration string, view string, zone string) (string, error) {

	zoneObj := models.Zone(entities.Zone{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
	})

	return objMgr.Connector.DeleteObject(zoneObj)
}

// Deployment role

// CreateDeploymentRole Create the Deployment role
func (objMgr *ObjectManager) CreateDeploymentRole(configuration string, view string, zone string, serverFQDN string, roleType string, role string, properties string, secondaryFQDN string) (*entities.DeploymentRole, error) {

	deploymentRole := models.NewDeploymentRole(entities.DeploymentRole{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		ServerFQDN:    serverFQDN,
		RoleType:      roleType,
		Role:          role,
		Properties:    properties,
		SecondaryFQDN: secondaryFQDN,
	})

	_, err := objMgr.Connector.CreateObject(deploymentRole)
	return deploymentRole, err
}

// GetDeploymentRole Get the Deployment role
func (objMgr *ObjectManager) GetDeploymentRole(configuration string, view string, zone string, serverFQDN string) (*entities.DeploymentRole, error) {

	deploymentRole := models.DeploymentRole(entities.DeploymentRole{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		ServerFQDN:    serverFQDN,
	})

	err := objMgr.Connector.GetObject(deploymentRole, &deploymentRole)
	return deploymentRole, err
}

// UpdateDeploymentRole Update the Deployment role
func (objMgr *ObjectManager) UpdateDeploymentRole(configuration string, view string, zone string, serverFQDN string, roleType string, role string, properties string, secondaryFQDN string) (*entities.DeploymentRole, error) {

	deploymentRole := models.DeploymentRole(entities.DeploymentRole{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		ServerFQDN:    serverFQDN,
		RoleType:      roleType,
		Role:          role,
		Properties:    properties,
		SecondaryFQDN: secondaryFQDN,
	})

	err := objMgr.Connector.UpdateObject(deploymentRole, &deploymentRole)
	return deploymentRole, err
}

// DeleteDeploymentRole Delete the Deployment role
func (objMgr *ObjectManager) DeleteDeploymentRole(configuration string, view string, zone string, serverFQDN string) (string, error) {

	deploymentRole := models.DeploymentRole(entities.DeploymentRole{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		ServerFQDN:    serverFQDN,
	})

	return objMgr.Connector.DeleteObject(deploymentRole)
}

// GetServer Get the Server info
func (objMgr *ObjectManager) GetServerByFQDN(configuration string, serverFQDN string) (*entities.Server, error) {

	server := models.Server(entities.Server{
		Configuration: configuration,
		ServerFQDN:    serverFQDN,
	})

	err := objMgr.Connector.GetObject(server, &server)
	return server, err
}
