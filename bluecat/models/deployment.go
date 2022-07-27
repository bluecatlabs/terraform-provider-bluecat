// Copyright 2022 BlueCat Networks. All rights reserved

package models

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/entities"
)

// NewDeploymentRole Initialize the new Deployment role to be added
func NewDeploymentRole(deploymentRole entities.DeploymentRole) *entities.DeploymentRole {
	res := deploymentRole
	res.SetObjectType("deployment_roles")
	sPath := getRRPrefixPath(deploymentRole.Configuration, deploymentRole.View)
	if len(deploymentRole.Zone) > 0 {
		sPath = fmt.Sprintf("%s/zones/%s", sPath, deploymentRole.Zone)
	}
	res.SetSubPath(sPath)
	return &res
}

// GetDeploymentRoles get all Deployment role on the Zone
func GetDeploymentRoles(deploymentRoles entities.DeploymentRoles) *entities.DeploymentRoles {
	res := deploymentRoles
	res.SetObjectType("deployment_roles")
	sPath := getRRPrefixPath(deploymentRoles.Configuration, deploymentRoles.View)
	if len(deploymentRoles.Zone) > 0 {
		sPath = fmt.Sprintf("%s/zones/%s", sPath, deploymentRoles.Zone)
	}
	res.SetSubPath(sPath)
	return &res
}

// DeploymentRole Initialize the Deployment role to be loaded, updated or deleted
func DeploymentRole(deploymentRole entities.DeploymentRole) *entities.DeploymentRole {
	res := deploymentRole
	res.SetObjectType("")
	sPath := getRRPrefixPath(deploymentRole.Configuration, deploymentRole.View)
	res.SetSubPath(fmt.Sprintf("%s/zones/%s/server/%s/deployment_roles", sPath, deploymentRole.Zone, deploymentRole.ServerFQDN))
	return &res
}
