// Copyright 2020 BlueCat Networks. All rights reserved

package entities

// Configuration Configuration entity
type Configuration struct {
	BAMBase    `json:"-"`
	Name       string `json:"name"`
	Properties string `json:"properties,omitempty"`
}
