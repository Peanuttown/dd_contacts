package models

import (
"github.com/Peanuttown/dd_contacts/ent"
)

type UserRequiredFields struct{
	UserId string
	Name string
	PropertiesInDepts []UserPropertiesInDepts
	Phone string
}

type UserPropertiesInDepts struct{
	DeptId uint
	IsDeptLeader bool
}

func NewUserRequiredFields(
	userId string,
	name string,
	propertiesInDepts []UserPropertiesInDepts,
	phone string,
)*UserRequiredFields{
	return &UserRequiredFields{
		UserId:userId,
		Name:name,
		PropertiesInDepts:propertiesInDepts,
		Phone:phone,
	}
}

type UserOptionalField func(userMutation *ent.UserMutation)

func UserOptionalFieldName(name string)UserOptionalField{
	return func(userMutation *ent.UserMutation){
		userMutation.SetName(name)
	}
}

func UserOptionlGeneration(generation uint)UserOptionalField{
	return func(mut *ent.UserMutation){
		mut.SetGeneration(generation)
	}
}
