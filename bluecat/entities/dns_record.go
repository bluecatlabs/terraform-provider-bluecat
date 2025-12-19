// Copyright 2020 BlueCat Networks. All rights reserved

package entities

// Zone the Zone entity
type Zone struct {
	BAMBase       `json:"-"`
	Configuration string   `json:"-"`
	View          string   `json:"-"`
	Zone          string   `json:"name,omitempty"`
	Deployable    string   `json:"deployable,omitempty"`
	ServerRoles   []string `json:"server_roles,omitempty"`
	Properties    string   `json:"properties,omitempty"`
	ZoneId        int      `json:"id,omitempty"`
}

// HostRecord Host record entity
type HostRecord struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	View          string `json:"-"`
	Zone          string `json:"-"`
	AbsoluteName  string `json:"absolute_name,omitempty"`
	IP4Address    string `json:"ip4_address,omitempty"`
	TTL           int    `json:"ttl,omitempty"`
	ReverseRecord string `json:"reverse_record,omitempty"`
	Properties    string `json:"properties,omitempty"`
	BAMId         int    `json:"id,omitempty"`
	BAMType       string `json:"type,omitempty"`
	Name          string `json:"name,omitempty"`
	BatchMode     string `json:"batch_mode,omitempty"`
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
	Name          string `json:"name,omitempty"`
	BAMId         int    `json:"id,omitempty"`
	BAMType       string `json:"type,omitempty"`
	BatchMode     string `json:"batch_mode,omitempty"`
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
	BAMId         int    `json:"id,omitempty"`
	BAMType       string `json:"type,omitempty"`
	BatchMode     string `json:"batch_mode,omitempty"`
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
	BAMId         int    `json:"id,omitempty"`
	BatchMode     string `json:"batch_mode,omitempty"`
}

// SRVRecord SRV record entity
type SRVRecord struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	View          string `json:"-"`
	Zone          string `json:"-"`
	AbsoluteName  string `json:"absolute_name,omitempty"`
	LinkedRecord  string `json:"linked_record,omitempty"`
	Priority      int    `json:"priority,omitempty"`
	Port          int    `json:"port,omitempty"`
	TTL           int    `json:"ttl,omitempty"`
	Weight        int    `json:"weight,omitempty"`
	Properties    string `json:"properties,omitempty"`
	Name          string `json:"name,omitempty"`
	BAMId         int    `json:"id,omitempty"`
	BAMType       string `json:"type,omitempty"`
	BatchMode     string `json:"batch_mode,omitempty"`
}

// External Host Record record entity
type ExternalHostRecord struct {
	BAMBase       `json:"-"`
	Configuration string `json:"-"`
	View          string `json:"-"`
	AbsoluteName  string `json:"name,omitempty"`
	Addresses     string `json:"addresses,omitempty"`
	Properties    string `json:"properties,omitempty"`
	BAMId         int    `json:"id,omitempty"`
	BAMType       string `json:"type,omitempty"`
	BatchMode     string `json:"batch_mode,omitempty"`
}
