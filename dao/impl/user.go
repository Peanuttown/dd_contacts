package impl

import(
"github.com/Peanuttown/dd_contacts/ent"
"fmt"
"github.com/Peanuttown/dd_contacts/dao/models"
"context"
)


type daoUser struct{
	client *ent.ClientWrapper
}

func NewDaoUser(client *ent.ClientWrapper) *daoUser{
	return &daoUser{
		client:client,
	}
}

	func (this *daoUser)CreateUser(ctx context.Context,requiredFields *models.UserRequiredFields,optionlFields ...models.UserOptionalField)(error){
		return this.client.TxDoIfClientNotTx(
			ctx,
			func(ctx context.Context,tx *ent.ClientWrapper)error{
				// steps:
				// 1. insert user 
				// 2. insert user's properties in depts;
				deptIds := make([]uint,0,len(requiredFields.PropertiesInDepts))
				for _,p := range requiredFields.PropertiesInDepts{
					deptIds = append(deptIds,p.DeptId)
				}
				_,err := tx.User.Create().
					AddDeptIDs(deptIds...).
					SetID(requiredFields.UserId).
					SetName(requiredFields.Name).
					SetPhone(requiredFields.Phone).Save(ctx)
					if err != nil{
						return fmt.Errorf("Create user failed: %w",err)
					}

				for _,p := range requiredFields.PropertiesInDepts{
					_,err =tx.UserPropertyInDept.Create().SetDeptID(p.DeptId).SetUserID(requiredFields.UserId).Save(ctx)
					if err != nil{
						return fmt.Errorf("Create UserPropertyInDept failed: %w",err)
					}
				}
				return nil
			},
		)
	}

	func (this *daoUser)DeleteUser(ctx context.Context,uid string)(error){
		return this.client.TxDoIfClientNotTx(
			ctx,
			func(ctx context.Context,tx *ent.ClientWrapper)error{
				return tx.User.DeleteOneID(uid).Exec(ctx)
			},
		)
	}
