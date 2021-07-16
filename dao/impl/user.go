package impl

import(
	"github.com/Peanuttown/dd_contacts/ent"
	"github.com/Peanuttown/dd_contacts/ent/user"
	up "github.com/Peanuttown/dd_contacts/ent/userpropertyindept"
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
			userCreate := tx.User.Create()
			return this.setFields(
				ctx,
				tx,
				userCreate.Mutation(),
				func(ctx context.Context)error{
					_,err := userCreate.Save(ctx)
					return err
				},
				requiredFields,
				optionlFields...,
			)
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

func (this *daoUser) Upsert(ctx context.Context,requiredFields *models.UserRequiredFields,optionlFields ...models.UserOptionalField)error{
	exist,err := this.client.User.Query().Where(user.IDEQ(requiredFields.UserId)).Exist(ctx)
	if err != nil{
		return err
	}
	return this.client.TxDoIfClientNotTx(
		ctx,
		func(ctx context.Context,tx *ent.ClientWrapper)error{
			if !exist{
				return this.CreateUser(ctx,requiredFields,optionlFields...)
			}
			userUpdate := tx.User.Update().Where(user.IDEQ(requiredFields.UserId))
			return this.setFields(
				ctx,
				tx,
				userUpdate.Mutation(),
				func(ctx context.Context)error{
					_,err := userUpdate.Save(ctx)
					return err
				},
				requiredFields,
				optionlFields...,
			)
		},
	)
}

func (this *daoUser) setFields(ctx context.Context,tx *ent.ClientWrapper,userMutation *ent.UserMutation,userMutationToSave func(ctx context.Context)error,requiredFields *models.UserRequiredFields,optionlFields ...models.UserOptionalField)error{
	// steps:
	// 1. insert user 
	// 2. insert user's properties in depts;
	deptIds := make([]uint,0,len(requiredFields.PropertiesInDepts))
	for _,p := range requiredFields.PropertiesInDepts{
		deptIds = append(deptIds,p.DeptId)
	}
	// add to dept if not in it
	depts,err := tx.User.Query().Where(user.IDEQ(requiredFields.UserId)).QueryDepts().All(ctx)
	OUT:
	for _,deptToAdd := range deptIds{
		for _,hasIn := range depts{
			if hasIn.ID == deptToAdd{
				continue OUT
			}
		}
		userMutation.AddDeptIDs(deptToAdd)
	}
	userMutation.SetID(requiredFields.UserId)
	userMutation.SetName(requiredFields.Name)
	userMutation.SetPhone(requiredFields.Phone)
	for _,v := range optionlFields{
		v(userMutation)
	}
	err = userMutationToSave(ctx)
	if err != nil{
		return fmt.Errorf("Create user failed: %w",err)
	}


	for _,p := range requiredFields.PropertiesInDepts{
		exist,err := tx.UserPropertyInDept.Query().Where(up.And(up.DeptIDEQ(p.DeptId),up.UserIDEQ(requiredFields.UserId))).Exist(ctx)
		if err != nil{
			return err
		}
		if exist{
			_,err =tx.UserPropertyInDept.Update().Where(up.And(up.DeptIDEQ(p.DeptId),up.UserIDEQ(requiredFields.UserId))).SetDeptID(p.DeptId).SetUserID(requiredFields.UserId).SetIsLeader(p.IsDeptLeader).Save(ctx)
		}else{
			fmt.Printf("Create property, %v, %v \n",p.DeptId,requiredFields.UserId)
			_,err =tx.UserPropertyInDept.Create().SetDeptID(p.DeptId).SetUserID(requiredFields.UserId).SetIsLeader(p.IsDeptLeader).Save(ctx)
		}
		if err != nil{
			return fmt.Errorf("Create UserPropertyInDept failed: %w",err)
		}
	}
	return nil
}


func (this *daoUser) FindByNotGeneration(ctx context.Context,generation uint)(ids []string,err error){
	ids,err = this.client.User.Query().
	Where(user.Or(user.GenerationNEQ(generation),user.GenerationIsNil())).Select(user.FieldID).Strings(ctx)
	return ids,err

	}
