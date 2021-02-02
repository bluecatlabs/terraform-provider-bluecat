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
