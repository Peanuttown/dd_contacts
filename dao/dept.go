package dao

import(
		"context"
"github.com/Peanuttown/dd_contacts/dao/models"
"github.com/Peanuttown/dd_contacts/ent"
)

type DaoDeptI interface{
	CreateDept(ctx context.Context,requiredFields *models.DeptRequriedFields ,optionalFields ...models.DeptOptionalFields)error
	FindSubDeptIds(ctx context.Context,deptId uint)([]uint,error)
	FindDept(ctx context.Context,deptId uint)(*ent.Dept,error)
	Upsert(ctx context.Context,requiredFields *models.DeptRequriedFields ,optionalFields ...models.DeptOptionalFields)error
	Delete(ctx context.Context,deptId uint)(error)
	FindByNotGeneration(ctx context.Context,generation uint)(deptIds []uint,err error)
}
