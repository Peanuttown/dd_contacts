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
	IsDeptLeader int
}

type UserOptionalField func(userMutation *ent.UserMutation)

func UserOptionalFieldName(name string)UserOptionalField{
	return func(userMutation *ent.UserMutation){
		userMutation.SetName(name)
	}
}
