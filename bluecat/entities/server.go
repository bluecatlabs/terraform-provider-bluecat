// Copyright 2023 BlueCat Networks. All rights reserved

package entities

type Server struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	Name          string `json:"name"`
	ServerFQDN    string `json:"fullHostName"`
	Type          string `json:"type"`
	Properties    string `json:"properties,omitempty"`
	ServerId      int    `json:"id,omitempty"`
}
