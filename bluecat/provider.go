// Copyright 2020 BlueCat Networks. All rights reserved

package bluecat

import (
	"terraform-provider-bluecat/bluecat/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Required:    true,
				Description: "Password to authenticate with BlueCat Gateway server.",
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"bluecat_host_record":    ResourceHostRecord(),
			"bluecat_configuration":  ResourceConfiguration(),
			"bluecat_ipv4block":      ResourceBlock(),
			"bluecat_ipv4network":    ResourceNetwork(),
			"bluecat_cname_record":   ResourceCNAMERecord(),
			"bluecat_ip_allocation":  ResourceIPAllocation(),
			"bluecat_ip_association": ResourceIPAssociation(),
			"bluecat_ptr_record":     ResourcePTRRecord(),
			"bluecat_txt_record":     ResourceTXTRecord(),
			"bluecat_generic_record": ResourceGenericRecord(),
			"bluecat_dhcp_range":     ResourceDHCPRange(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bluecat_ipv4network":  DataSourceIPv4Network(),
			"bluecat_cname_record": DataSourceCNAMERecord(),
			"bluecat_host_record":  DataSourceHostRecord(),
			"bluecat_ipv4block":    DataSourceBlock(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	hostConfig := utils.HostConfig{
		Host:      d.Get("server").(string),
		Port:      d.Get("port").(string),
		Transport: d.Get("transport").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
		Version:   d.Get("api_version").(string),
	}

	requestBuilder := &utils.APIRequestBuilder{}
	requester := &utils.APIHttpRequester{}

	conn, err := utils.NewConnector(hostConfig, requestBuilder, requester)
	if err != nil {
		log.Debugf("Failed to initialize the provider: %s", err)
		return nil, err
	}
	return conn, err
}
