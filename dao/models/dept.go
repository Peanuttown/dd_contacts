package models

import(
"github.com/Peanuttown/dd_contacts/ent"

)

type DeptRequriedFields struct{
	DeptId uint
	Name string
}

func NewDeptRequiredFields(deptId uint,name string)*DeptRequriedFields{
	return &DeptRequriedFields{
		DeptId:deptId,
		Name:name,
	}
}

type DeptOptionalFields = func(*ent.DeptMutation)

func DeptOptionalParentDeptId(parentId uint)DeptOptionalFields{
	return func(mut *ent.DeptMutation){
		if parentId  == 0{
			return
		}
		mut.SetParentID(parentId)
	}
}

func DeptOptionalGeneration(generation uint)DeptOptionalFields{
	return func(mut *ent.DeptMutation){
		mut.SetGeneration(generation)
	}

}
