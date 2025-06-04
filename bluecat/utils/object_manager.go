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
func (objMgr *ObjectManager) CreateBlock(block entities.Block) (*entities.Block, error) {

	// default value for the ipVersion is ipv4
	if block.IPVersion == "" {
		block.IPVersion = entities.IPV4
	}
	block = models.NewBlock(
		entities.Block{
			Configuration: block.Configuration,
			Name:          block.Name,
			Address:       block.Address,
			CIDR:          block.CIDR,
			ParentBlock:   block.ParentBlock,
			Properties:    block.Properties,
		},
		block.IPVersion,
	)

	_, err := objMgr.Connector.CreateObject(&block)
	return &block, err
}

// GetBlock Get the Block info
func (objMgr *ObjectManager) GetBlock(configuration string, address string, cidr string, ipVersion string) (*entities.Block, error) {

	// default value for the ipVersion is ipv4
	if ipVersion == "" {
		ipVersion = entities.IPV4
	}
	block := models.IPBlock(
		entities.Block{
			Configuration: configuration,
			Address:       address,
			CIDR:          cidr,
			IPVersion:     ipVersion,
		},
	)

	err := objMgr.Connector.GetObject(block, &block)
	return block, err
}

// UpdateBlock Update the Block info
func (objMgr *ObjectManager) UpdateBlock(block entities.Block) (*entities.Block, error) {

	// default value for the ipVersion is ipv4
	if block.IPVersion == "" {
		block.IPVersion = entities.IPV4
	}
	blockEntity := models.IPBlock(
		entities.Block{
			Configuration: block.Configuration,
			Name:          block.Name,
			Address:       block.Address,
			CIDR:          block.CIDR,
			ParentBlock:   block.ParentBlock,
			Properties:    block.Properties,
			IPVersion:     block.IPVersion,
		},
	)

	err := objMgr.Connector.UpdateObject(blockEntity, &block)
	return blockEntity, err
}

// DeleteBlock Delete the Block
func (objMgr *ObjectManager) DeleteBlock(configuration string, address string, cidr string, ipVersion string) (string, error) {

	// default value for the ipVersion is ipv4
	if ipVersion == "" {
		ipVersion = entities.IPV4
	}
	block := models.IPBlock(
		entities.Block{
			Configuration: configuration,
			Address:       address,
			CIDR:          cidr,
			IPVersion:     ipVersion,
		},
	)

	return objMgr.Connector.DeleteObject(block)
}

// Network

func generateNetworkProperties(props string, gateway string) string {
	result := props
	if len(gateway) > 0 {
		result = fmt.Sprintf("%s|gateway=%s", result, gateway)
	}
	return result
}

// CreateNetwork Create a new Network
func (objMgr *ObjectManager) CreateNetwork(network entities.Network) (*entities.Network, error) {

	networkEntity := entities.Network{
		Configuration: network.Configuration,
		BlockAddr:     network.BlockAddr,
		Name:          network.Name,
		CIDR:          network.CIDR,
		Properties:    network.Properties,
		Template:      network.Template,
		IPVersion:     network.IPVersion,
	}
	if networkEntity.IPVersion == entities.IPV4 || networkEntity.IPVersion == "" {
		networkEntity.Gateway = network.Gateway
		networkEntity.Properties = generateNetworkProperties(network.Properties, network.Gateway)
	}

	network = models.NewNetwork(networkEntity)
	_, err := objMgr.Connector.CreateObject(&network)
	return &network, err
}

// CreateNextAvailableNetwork Create a next available Network
func (objMgr *ObjectManager) CreateNextAvailableNetwork(network entities.Network) (*entities.Network, string, error) {

	networkEntity := models.NewNextAvailableNetwork(entities.Network{
		Configuration: network.Configuration,
		BlockAddr:     network.BlockAddr,
		Name:          network.Name,
		Gateway:       network.Gateway,
		Properties:    network.Properties,
		Template:      network.Template,
		Size:          network.Size,
		AllocatedId:   network.AllocatedId,
		IPVersion:     network.IPVersion,
	})

	if networkEntity.IPVersion == entities.IPV4 || networkEntity.IPVersion == "" {
		networkEntity.Gateway = network.Gateway
		networkEntity.Properties = generateNetworkProperties(network.Properties, network.Gateway)
	}

	ref, err := objMgr.Connector.CreateObject(networkEntity)
	return networkEntity, ref, err
}

// GetNetwork Get the Network info
func (objMgr *ObjectManager) GetNetwork(network *entities.Network) (*entities.Network, error) {

	networkEntity := models.Network(entities.Network{
		Configuration: network.Configuration,
		CIDR:          network.CIDR,
		IPVersion:     network.IPVersion,
	})

	err := objMgr.Connector.GetObject(networkEntity, &network)
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
func (objMgr *ObjectManager) UpdateNetwork(network entities.Network) (*entities.Network, error) {

	networkEntity := models.Network(entities.Network{
		Configuration: network.Configuration,
		Name:          network.Name,
		CIDR:          network.CIDR,
		Properties:    network.Properties,
		IPVersion:     network.IPVersion,
	})

	if networkEntity.IPVersion == entities.IPV4 || networkEntity.IPVersion == "" {
		networkEntity.Gateway = network.Gateway
		networkEntity.Properties = generateNetworkProperties(network.Properties, network.Gateway)
	}

	err := objMgr.Connector.UpdateObject(networkEntity, &network)
	return networkEntity, err
}

// DeleteNetwork Delete the Network
func (objMgr *ObjectManager) DeleteNetwork(network entities.Network) (string, error) {

	networkEntity := models.Network(entities.Network{
		Configuration: network.Configuration,
		CIDR:          network.CIDR,
		IPVersion:     network.IPVersion,
	})

	return objMgr.Connector.DeleteObject(networkEntity)
}

// DHCP Range

// CreateDHCPRange Create a new DHCP Range
func (objMgr *ObjectManager) CreateDHCPRange(dhcpRange entities.DHCPRange) (*entities.DHCPRange, error) {
	dhcpRangeEntity := models.NewDHCPRange(dhcpRange)
	_, err := objMgr.Connector.CreateObject(dhcpRangeEntity)
	return dhcpRangeEntity, err
}

// GetDHCPRange Get the DHCP Range info
func (objMgr *ObjectManager) GetDHCPRange(dhcpRange entities.DHCPRange) (*entities.DHCPRange, error) {
	dhcpRangeEntity := models.DHCPRange(dhcpRange)
	err := objMgr.Connector.GetObject(dhcpRangeEntity, &dhcpRangeEntity)
	return dhcpRangeEntity, err
}

// GetDeploymentRoles Get all Deployment role on the Zone
func (objMgr *ObjectManager) GetDeploymentRoles(configuration string, view string, zone string) (*entities.DeploymentRoles, error) {
	var deploymentRoles *entities.DeploymentRoles
	if zone == "" {
		deploymentRoles = models.GetDeploymentRoles(entities.DeploymentRoles{
			Configuration: configuration,
			View:          view,
		})
	} else {
		deploymentRoles = models.GetDeploymentRoles(entities.DeploymentRoles{
			Configuration: configuration,
			View:          view,
			Zone:          zone,
		})
	}

	err := objMgr.Connector.GetObject(deploymentRoles, &deploymentRoles)
	return deploymentRoles, err
}

// UpdateDHCPRange Update the DHCP Range info
func (objMgr *ObjectManager) UpdateDHCPRange(dhcpRange entities.DHCPRange) (*entities.DHCPRange, error) {
	dhcpRangeEntity := models.DHCPRange(dhcpRange)
	err := objMgr.Connector.UpdateObject(dhcpRangeEntity, &dhcpRangeEntity)
	return dhcpRangeEntity, err
}

// DeleteDHCPRange Delete the DHCP Range
func (objMgr *ObjectManager) DeleteDHCPRange(dhcpRange entities.DHCPRange) (string, error) {
	dhcpRangeEntity := models.DHCPRange(dhcpRange)
	return objMgr.Connector.DeleteObject(dhcpRangeEntity)
}

// IP

// ReserveIPAddress Create the new IP address for later use
func (objMgr *ObjectManager) ReserveIPAddress(configuration string, network string, ipVersion string) (*entities.IPAddress, error) {
	address := entities.IPAddress{
		Configuration: configuration,
		CIDR:          network,
		Name:          "",
		Address:       "",
		Mac:           "",
		Action:        entities.AllocateReserved,
		Properties:    "",
		Template:      "",
		IPVersion:     ipVersion,
	}
	return objMgr.CreateIPAddress(address)
}

// createIPAddress Create the new IP address. Allocate the next available on the network if IP address is not provided
func (objMgr *ObjectManager) CreateIPAddress(address entities.IPAddress) (*entities.IPAddress, error) {
	if len(address.Action) == 0 {
		address.Action = entities.AllocateStatic
	}

	ipAddr := new(entities.IPAddress)
	if len(address.Address) > 0 {
		ipAddr = models.IPAddress(address)
	} else {
		ipAddr = models.GetNextIPAddress(address)
		log.Debugf("Requesting the new IP address in the network %s", address.CIDR)
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
func (objMgr *ObjectManager) GetIPAddress(configuration string, address string, ipVersion string) (*entities.IPAddress, error) {

	ipAddr := models.IPAddress(entities.IPAddress{
		Configuration: configuration,
		Address:       address,
		IPVersion:     ipVersion,
	})

	err := objMgr.Connector.GetObject(ipAddr, &ipAddr)
	return ipAddr, err
}

// SetMACAddress Update the MAC address for the existing IP address
func (objMgr *ObjectManager) SetMACAddress(address entities.IPAddress) (*entities.IPAddress, error) {
	address.Properties = ""
	ipAddr := models.IPAddress(address)
	err := objMgr.Connector.UpdateObject(ipAddr, &ipAddr)
	return ipAddr, err
}

// UpdateIPAddress Update the IP address info
func (objMgr *ObjectManager) UpdateIPAddress(address entities.IPAddress) (*entities.IPAddress, error) {
	ipAddr := models.IPAddress(address)
	ipAddr.SetAction()
	err := objMgr.Connector.UpdateObject(ipAddr, &ipAddr)
	return ipAddr, err
}

// DeleteIPAddress Delete the existing IP address
func (objMgr *ObjectManager) DeleteIPAddress(configuration string, address string, ipVersion string) (string, error) {
	ipAddr := models.IPAddress(entities.IPAddress{
		Configuration: configuration,
		Address:       address,
		IPVersion:     ipVersion,
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

// CreateSRVRecord Create the SRV record
func (objMgr *ObjectManager) CreateSRVRecord(configuration string, view string, zone string, priority int, port int, weight int, absoluteName string, linkedRecord string, ttl int, properties string) (*entities.SRVRecord, error) {

	srvRecord := models.NewSRVRecord(entities.SRVRecord{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		LinkedRecord:  linkedRecord,
		Priority:      priority,
		Port:          port,
		Weight:        weight,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
	})

	_, err := objMgr.Connector.CreateObject(srvRecord)
	return srvRecord, err
}

// GetSRVRecord Get the SRV record
func (objMgr *ObjectManager) GetSRVRecord(configuration string, view string, absoluteName string) (*entities.SRVRecord, error) {

	srvRecord := models.SRVRecord(entities.SRVRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	err := objMgr.Connector.GetObject(srvRecord, &srvRecord)
	return srvRecord, err
}

// UpdateSRVRecord Update the SRV record
func (objMgr *ObjectManager) UpdateSRVRecord(configuration string, view string, zone string, priority int, port int, weight int, absoluteName string, linkedRecord string, ttl int, properties string, name string) (*entities.SRVRecord, error) {

	srvRecord := models.SRVRecord(entities.SRVRecord{
		Configuration: configuration,
		View:          view,
		Zone:          zone,
		LinkedRecord:  linkedRecord,
		Priority:      priority,
		Port:          port,
		Weight:        weight,
		AbsoluteName:  absoluteName,
		TTL:           ttl,
		Properties:    properties,
		Name:          name,
	})

	err := objMgr.Connector.UpdateObject(srvRecord, &srvRecord)
	return srvRecord, err
}

// DeleteSRVRecord Delete the SRV record
func (objMgr *ObjectManager) DeleteSRVRecord(configuration string, view string, absoluteName string) (string, error) {

	srvRecord := models.SRVRecord(entities.SRVRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	return objMgr.Connector.DeleteObject(srvRecord)
}

// CreateSRVRecord Create the SRV record
func (objMgr *ObjectManager) CreateExternalHostRecord(configuration string, view string, addresses string, absoluteName string, properties string) (*entities.ExternalHostRecord, error) {

	externalHostRecord := models.NewExternalHostRecord(entities.ExternalHostRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
		Properties:    properties,
		Addresses:     addresses,
	})

	_, err := objMgr.Connector.CreateObject(externalHostRecord)
	return externalHostRecord, err
}

// GetSRVRecord Get the SRV record
func (objMgr *ObjectManager) GetExternalHostRecord(configuration string, view string, absoluteName string) (*entities.ExternalHostRecord, error) {

	externalHostRecord := models.ExternalHostRecord(entities.ExternalHostRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	err := objMgr.Connector.GetObject(externalHostRecord, &externalHostRecord)
	return externalHostRecord, err
}

// UpdateSRVRecord Update the SRV record
func (objMgr *ObjectManager) UpdateExternalHostRecord(configuration string, view string, addresses string, absoluteName string, properties string) (*entities.ExternalHostRecord, error) {

	externalHostRecord := models.ExternalHostRecord(entities.ExternalHostRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
		Addresses:     addresses,
		Properties:    properties,
	})

	err := objMgr.Connector.UpdateObject(externalHostRecord, &externalHostRecord)
	return externalHostRecord, err
}

// DeleteSRVRecord Delete the SRV record
func (objMgr *ObjectManager) DeleteExternalHostRecord(configuration string, view string, absoluteName string) (string, error) {

	externalHostRecord := models.ExternalHostRecord(entities.ExternalHostRecord{
		Configuration: configuration,
		View:          view,
		AbsoluteName:  absoluteName,
	})

	return objMgr.Connector.DeleteObject(externalHostRecord)
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

// CreateView Create a new View
func (objMgr *ObjectManager) CreateView(configuration string, name string, properties string) (*entities.View, error) {

	view := models.NewView(&entities.View{
		Configuration: configuration,
		Name:          name,
		Properties:    properties,
	})

	_, err := objMgr.Connector.CreateObject(view)
	return view, err
}

func (objMgr *ObjectManager) GetView(configuration string, name string) (*entities.View, error) {

	view := models.View(entities.View{
		Configuration: configuration,
		Name:          name,
	})

	err := objMgr.Connector.GetObject(view, &view)
	return view, err
}

// UpdateView Update the View info
func (objMgr *ObjectManager) UpdateView(configuration string, name string, properties string) (*entities.View, error) {

	view := models.View(entities.View{
		Configuration: configuration,
		Name:          name,
		Properties:    properties,
	})

	err := objMgr.Connector.UpdateObject(view, &view)
	return view, err
}

// DeleteView Delete the View
func (objMgr *ObjectManager) DeleteView(configuration string, name string) (string, error) {

	view := models.View(entities.View{
		Configuration: configuration,
		Name:          name,
	})

	return objMgr.Connector.DeleteObject(view)
}
