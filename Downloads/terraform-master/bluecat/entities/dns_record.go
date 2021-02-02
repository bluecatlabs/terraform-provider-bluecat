// Copyright 2020 BlueCat Networks. All rights reserved

package entities

// HostRecord Host record entity
type HostRecord struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	View          string `json:"-"`
	Zone          string `json:"-"`
	AbsoluteName  string `json:"absolute_name,omitempty"`
	IP4Address    string `json:"ip4_address,omitempty"`
	TTL           int    `json:"ttl,omitempty"`
	Properties    string `json:"properties,omitempty"`
}

// CNAMERecord CNAME record entity
type CNAMERecord struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	View          string `json:"-"`
	Zone          string `json:"-"`
	AbsoluteName  string `json:"absolute_name,omitempty"`
	LinkedRecord  string `json:"linked_record,omitempty"`
	TTL           int    `json:"ttl,omitempty"`
	Properties    string `json:"properties,omitempty"`
}
