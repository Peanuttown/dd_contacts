package impl

import(
	"strconv"
	"fmt"
	"context"		
	"github.com/Peanuttown/dd_contacts/ent"
	"github.com/Peanuttown/dd_contacts/ent/dept"
	"github.com/Peanuttown/dd_contacts/dao"
	"github.com/Peanuttown/dd_contacts/dao/models"
)

type daoDept struct{
	client *ent.ClientWrapper
}

func NewDaoDept(client *ent.ClientWrapper) dao.DaoDeptI{
	return &daoDept{
		client:client,
	}
}

func (this *daoDept) FindDept(ctx context.Context,deptId uint)(*ent.Dept,error){
	return this.client.Dept.Query().Where(dept.IDEQ(deptId)).Only(ctx)
}

func (this *daoDept) CreateDept(ctx context.Context,requiredFields *models.DeptRequriedFields ,optionalFields ...models.DeptOptionalFields)error{
	createDept := this.client.Dept.Create()
	mut := createDept.Mutation()
	this.setField(ctx,mut,requiredFields,optionalFields...)
	_,err := createDept.Save(ctx)
	return err
}

func (this *daoDept)FindSubDeptIds(ctx context.Context,deptId uint)([]uint,error){
	dept,err := this.FindDept(ctx,deptId)
	if err != nil{
		return nil,err
	}
	subDepts,err := this.client.Dept.QuerySubDepts(dept).All(ctx)
	if err != nil{
		return nil,err
	}
	retIds := make([]uint,0,len(subDepts))
	for _,v := range subDepts{
		retIds = append(retIds,v.ID)
	}
	return retIds,nil
}

func (this *daoDept) Upsert(ctx context.Context,requiredFields *models.DeptRequriedFields ,optionalFields ...models.DeptOptionalFields)error{
	// query if exists
	exist,err := this.client.Dept.Query().Where(dept.IDEQ(requiredFields.DeptId)).Exist(ctx)
	if err != nil{
		return err
	}
	if !exist{
		return this.CreateDept(ctx,requiredFields,optionalFields...)
	}
	deptUpdate := this.client.Dept.Update().Where(dept.IDEQ(requiredFields.DeptId))
	mut := deptUpdate.Mutation()
	this.setField(ctx,mut,requiredFields,optionalFields...)
	_,err = deptUpdate.Save(ctx)
	return err
}

func (this *daoDept) setField(ctx context.Context, mut *ent.DeptMutation,requiredFields *models.DeptRequriedFields ,optionalFields ...models.DeptOptionalFields){
	mut.SetID(requiredFields.DeptId)
	fmt.Println(requiredFields.Name)
	mut.SetName(requiredFields.Name)
	for _,v := range optionalFields{
		 v(mut)
	}
}

func(this *daoDept)Delete(ctx context.Context,deptId uint)(error){
	return this.client.TxDoIfClientNotTx(
		ctx,
		func(ctx context.Context , tx *ent.ClientWrapper)error{
			return tx.Dept.DeleteOneID(deptId).Exec(ctx)
		},
	)

}
func(this *daoDept)FindByNotGeneration(ctx context.Context,generation uint)(deptIds []uint,err error){
	ids,err := this.client.Dept.Query().Where(dept.Or(dept.GenerationNEQ(generation),dept.GenerationIsNil())).Select(dept.FieldID).Strings(ctx)
	if err != nil{
		return nil,err
	}
	depts := make([]uint,0,len(ids))
	for _,id:=range ids{
		idUint,err := strconv.ParseUint(id,10,64)
		if err != nil{
			return nil,fmt.Errorf("Parse deptId %v to uint failed: %w",id,err)
		}
		depts = append(depts,uint(idUint))
	}
	return depts,nil
}
