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

// TXTRecord TXT record entity
type TXTRecord struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	View          string `json:"-"`
	Zone          string `json:"-"`
	AbsoluteName  string `json:"absolute_name,omitempty"`
	Text          string `json:"text,omitempty"`
	TTL           int    `json:"ttl,omitempty"`
	Properties    string `json:"properties,omitempty"`
}

// GenericRecord Generic record entity
type GenericRecord struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	View          string `json:"-"`
	Zone          string `json:"-"`
	TypeRR        string `json:"type,omitempty"`
	AbsoluteName  string `json:"absolute_name,omitempty"`
	Data          string `json:"data,omitempty"`
	TTL           int    `json:"ttl,omitempty"`
	Properties    string `json:"properties,omitempty"`
}
