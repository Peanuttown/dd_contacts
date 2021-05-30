package dao

import(
	"github.com/Peanuttown/dd_contacts/dao/models"
	"context"
)

type DaoUserI interface{
	CreateUser(ctx context.Context,requiredFields *models.UserRequiredFields,optionlFields ...models.UserOptionalField)(error)
	DeleteUser(ctx context.Context,uid string)(error)
	Upsert(ctx context.Context,requiredFields *models.UserRequiredFields,optionlFields ...models.UserOptionalField)error
	FindByNotGeneration(ctx context.Context,generation uint)(ids []string,err error)
}
