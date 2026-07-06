// Copyright 2022 BlueCat Networks. All rights reserved

package entities

// DeploymentRole the Deployment role entity
type DeploymentRole struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	View          string `json:"-"`
	Zone          string `json:"-"`
	ServerFQDN    string `json:"server_fqdn,omitempty"`
	RoleType      string `json:"role_type,omitempty"`
	Role          string `json:"role,omitempty"`
	Properties    string `json:"properties,omitempty"`
	SecondaryFQDN string `json:"secondary_fqdn,omitempty"`
}

// DeploymentRoles the list Deployment role entity
type DeploymentRoles struct {
	BAMBase       `json:"-"`
	Configuration string           `json:"-"`
	View          string           `json:"-"`
	Zone          string           `json:"-"`
	ServerRoles   []DeploymentRole `json:"deployment_roles, omitempty"`
}

// DeploymentOption the Deployment option entity
type DeploymentOption struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	View          string `json:"-"`
	Zone          string `json:"-"`
	ResourceType  string `json:"-"`
	ResourceRef   string `json:"-"`
	IPVersion     string `json:"-"`
	Name          string `json:"name,omitempty"`
	Value         string `json:"value,omitempty"`
	ServerID      int    `json:"-"`
	Properties    string `json:"properties,omitempty"`
}
