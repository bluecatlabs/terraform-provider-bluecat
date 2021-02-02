// Copyright 2020 BlueCat Networks. All rights reserved

package entities

import "fmt"

// Block IPv4 Block entity
type Block struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	ParentBlock   string `json:"-"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	CIDR          string `json:"cidr_notation"`
	Properties    string `json:"properties,omitempty"`
}

// AddressCIDR Get the Block address in CIDR format
func (block *Block) AddressCIDR() string {
	return fmt.Sprintf("%s/%s", block.Address, block.CIDR)
}

// Network IPv4 Network entity
type Network struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	BlockAddr     string `json:"-"`
	Name          string `json:"name"`
	CIDR          string `json:"cidr"`
	Gateway       string `json:"gateway"`
	Properties    string `json:"properties,omitempty"`
}

// IPAddress The IPv4 Address entity
type IPAddress struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	Action        string `json:"action,omitempty"`
	CIDR          string `json:"network,omitempty"`
	Address       string `json:"ipv4addr,omitempty"`
	Mac           string `json:"mac_address,omitempty"`
	Name          string `json:"name,omitempty"`
	Properties    string `json:"properties,omitempty"`
}
