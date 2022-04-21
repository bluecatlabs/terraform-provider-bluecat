// Copyright 2020 BlueCat Networks. All rights reserved

package models

import "terraform-provider-bluecat/bluecat/entities"

// RestLogin Initialize the Rest credentials
func RestLogin(cred entities.RestLogin) *entities.RestLogin {
	res := cred
	res.SetObjectType("")
	res.SetSubPath("/token")
	return &res
}
