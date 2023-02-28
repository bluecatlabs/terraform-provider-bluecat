// Copyright 2023 BlueCat Networks. All rights reserved
package models

import (
	"fmt"
	"terraform-provider-bluecat/bluecat/entities"
)

// Server Initialize the Server to be loaded, updated or deleted
func Server(server entities.Server) *entities.Server {
	res := server
	res.SetObjectType("")
	res.SetSubPath(fmt.Sprintf("%s/server_fqdn/%s", getPath(res.Configuration), server.ServerFQDN))
	return &res
}
