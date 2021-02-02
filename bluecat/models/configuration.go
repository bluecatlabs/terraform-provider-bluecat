// Copyright 2020 BlueCat Networks. All rights reserved

package models

import "terraform-provider-bluecat/bluecat/entities"

func NewConfiguration(configuration entities.Configuration) *entities.Configuration {
	res := configuration
	res.SetObjectType("configurations")
	res.SetSubPath("")
	return &res
}

func Configuration(configuration entities.Configuration) *entities.Configuration {
	res := configuration
	res.SetObjectType("")
	res.SetSubPath("/configurations/" + configuration.Name)
	return &res
}
