package dao

import(
		"context"
"github.com/Peanuttown/dd_contacts/dao/models"
)

type DaoDeptI interface{
	CreateDept(ctx context.Context,requiredFields *models.DeptRequriedFields ,optionalFields ...models.DeptOptionalFields)
}
