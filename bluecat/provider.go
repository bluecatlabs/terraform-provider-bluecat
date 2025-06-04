// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"context"
	"fmt"
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider BlueCat provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "BlueCat Gateway IP address.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User to authenticate with BlueCat Gateway server.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Password to authenticate with BlueCat Gateway server. The encrypted file name if encrypt_password set to True",
			},
			"api_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API Version of REST_API workflow server",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Port number used for connection for BlueCat Gateway Server.",
			},
			"transport": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Transport type (HTTP or HTTPS).",
			},
			"encrypt_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Default is false, to indicate if the password is encrypted",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bluecat_host_record":          ResourceHostRecord(),
			"bluecat_configuration":        ResourceConfiguration(),
			"bluecat_ipv4block":            ResourceBlock(),
			"bluecat_ipv6block":            ResourceBlock(),
			"bluecat_ipv4network":          ResourceNetwork(),
			"bluecat_ipv6network":          ResourceNetwork(),
			"bluecat_cname_record":         ResourceCNAMERecord(),
			"bluecat_ip_allocation":        ResourceIPAllocation(),
			"bluecat_ip_association":       ResourceIPAssociation(),
			"bluecat_ptr_record":           ResourcePTRRecord(),
			"bluecat_txt_record":           ResourceTXTRecord(),
			"bluecat_srv_record":           ResourceSRVRecord(),
			"bluecat_external_host_record": ResourceExternalHostRecord(),
			"bluecat_generic_record":       ResourceGenericRecord(),
			"bluecat_dhcp_range":           ResourceDHCPRange(),
			"bluecat_zone":                 ResourceZone(),
			"bluecat_view":                 ResourceView(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bluecat_ipv4network":  DataSourceIPv4Network(),
			"bluecat_ipv6network":  DataSourceIPv4Network(),
			"bluecat_cname_record": DataSourceCNAMERecord(),
			"bluecat_host_record":  DataSourceHostRecord(),
			"bluecat_ipv4block":    DataSourceBlock(),
			"bluecat_ipv6block":    DataSourceBlock(),
			"bluecat_zone":         DataSourceZone(),
			"bluecat_view":         DataSourceView(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	hostConfig := utils.HostConfig{
		Host:            d.Get("server").(string),
		Port:            d.Get("port").(string),
		Transport:       d.Get("transport").(string),
		Username:        d.Get("username").(string),
		Password:        d.Get("password").(string),
		Version:         d.Get("api_version").(string),
		EncryptPassword: d.Get("encrypt_password").(bool),
	}

	requestBuilder := &utils.APIRequestBuilder{}
	requester := &utils.APIHttpRequester{}

	conn, err := utils.NewConnector(hostConfig, requestBuilder, requester)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to initialize the provider: %s", err),
		})
		return nil, diags
	}
	return conn, diags
}

func GetObjManager(m interface{}) *utils.ObjectManager {
	connector := m.(*utils.Connector)
	objMgr := new(utils.ObjectManager)
	objMgr.Connector = connector
	return objMgr
}
