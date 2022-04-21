// Copyright 2020 BlueCat Networks. All rights reserved

package entities

type BAMObject interface {
	ObjectType() string
	SetObjectType(objType string)
	SubPath() string
	SetSubPath(subPath string)
}

type BAMBase struct {
	objectType string
	subPath    string
}

func (obj *BAMBase) ObjectType() string {
	return obj.objectType
}

func (obj *BAMBase) SetObjectType(objType string) {
	obj.objectType = objType
}

func (obj *BAMBase) SubPath() string {
	return obj.subPath
}

func (obj *BAMBase) SetSubPath(subPath string) {
	obj.subPath = subPath
}

type RestLogin struct {
	BAMBase
	UserName        string `json:"username"`
	Password        string `json:"password"`
	EncryptPassword bool   `json:"encrypt_password"`
}
type EntityID struct {
	ID int
}
