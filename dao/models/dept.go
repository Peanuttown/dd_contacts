package models

import(
"github.com/Peanuttown/dd_contacts/ent"

)

type DeptRequriedFields struct{

	DeptId uint
	Name string
}

type DeptOptionalFields = func(*ent.DeptMutation)

