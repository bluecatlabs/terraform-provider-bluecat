// Copyright 2020 BlueCat Networks. All rights reserved

package models

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/entities"
)

func getRRPrefixPath(configuration, view string) string {
	result := ""
	if len(configuration) > 0 && len(view) > 0 {
		result = fmt.Sprintf("/configurations/%s/views/%s", configuration, view)
	}

	return result
}

// Zone
// Zone Initialize the new Zone to be added
func NewZone(zone entities.Zone) *entities.Zone {
	res := zone
	res.SetObjectType("zones")
	sPath := getRRPrefixPath(zone.Configuration, zone.View)
	res.SetSubPath(sPath)
	return &res
}

// Zone Initialize the Zone to be loaded, updated or deleted
func Zone(zone entities.Zone) *entities.Zone {
	res := zone
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/zones/%s", getRRPrefixPath(zone.Configuration, zone.View), zone.Zone))
	return &res
}

// NewHostRecord Initialize the new Host record to be added
func NewHostRecord(hostRecord entities.HostRecord) *entities.HostRecord {
	res := hostRecord
	res.SetObjectType("host_records")
	sPath := getRRPrefixPath(hostRecord.Configuration, hostRecord.View)
	if len(hostRecord.Zone) > 0 {
		sPath = fmt.Sprintf("%s/zones/%s", sPath, hostRecord.Zone)
	}
	res.SetSubPath(sPath)
	return &res
}

// HostRecord Initialize the Host record to be loaded, updated or deleted
func HostRecord(hostRecord entities.HostRecord) *entities.HostRecord {
	res := hostRecord
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/host_records/%s", getRRPrefixPath(hostRecord.Configuration, hostRecord.View), hostRecord.AbsoluteName))
	return &res
}

// NewCNAMERecord Initialize the new CNAME record to be added
func NewCNAMERecord(cnameRecord entities.CNAMERecord) *entities.CNAMERecord {
	res := cnameRecord
	res.SetObjectType("cname_records")
	sPath := getRRPrefixPath(cnameRecord.Configuration, cnameRecord.View)
	if len(cnameRecord.Zone) > 0 {
		sPath = fmt.Sprintf("%s/zones/%s", sPath, cnameRecord.Zone)
	}
	res.SetSubPath(sPath)
	return &res
}

// CNAMERecord Initialize the CNAME record to be loaded, updated or deleted
func CNAMERecord(cnameRecord entities.CNAMERecord) *entities.CNAMERecord {
	res := cnameRecord
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/cname_records/%s", getRRPrefixPath(cnameRecord.Configuration, cnameRecord.View), cnameRecord.AbsoluteName))
	return &res
}

// TXTRecord Initialize the TXT record to be loaded, updated or deleted
func TXTRecord(txtRecord entities.TXTRecord) *entities.TXTRecord {
	res := txtRecord
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/text_records/%s", getRRPrefixPath(txtRecord.Configuration, txtRecord.View), txtRecord.AbsoluteName))
	return &res
}

// NewTXTRecord Initialize the new TXT record to be added
func NewTXTRecord(txtRecord entities.TXTRecord) *entities.TXTRecord {
	res := txtRecord
	res.SetObjectType("text_records")
	sPath := getRRPrefixPath(txtRecord.Configuration, txtRecord.View)
	if len(txtRecord.Zone) > 0 {
		sPath = fmt.Sprintf("%s/zones/%s", sPath, txtRecord.Zone)
	}
	res.SetSubPath(sPath)
	return &res
}

// GenericRecord Initialize the Generic record to be loaded, updated or deleted
func GenericRecord(genericRecord entities.GenericRecord) *entities.GenericRecord {
	res := genericRecord
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/generic_records/%s", getRRPrefixPath(genericRecord.Configuration, genericRecord.View), genericRecord.AbsoluteName))
	return &res
}

// NewGenericRecord Initialize the new Generic record to be added
func NewGenericRecord(genericRecord entities.GenericRecord) *entities.GenericRecord {
	res := genericRecord
	res.SetObjectType("generic_records")
	sPath := getRRPrefixPath(genericRecord.Configuration, genericRecord.View)
	if len(genericRecord.Zone) > 0 {
		sPath = fmt.Sprintf("%s/zones/%s", sPath, genericRecord.Zone)
	}
	res.SetSubPath(sPath)
	return &res
}
