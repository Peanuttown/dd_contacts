package schema

import (
	"fmt"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"github.com/Peanuttown/dd_contacts/ent/hook"
	ent_gen "github.com/Peanuttown/dd_contacts/ent"
	"github.com/Peanuttown/dd_contacts/ent/user"
	up "github.com/Peanuttown/dd_contacts/ent/userpropertyindept"
	//"github.com/Peanuttown/dd_contacts/ent/dept"
	"context"

)

// Dept holds the schema definition for the Dept entity.
type Dept struct {
	ent.Schema
}

// Fields of the Dept.
func (Dept) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").StorageKey("dept_id"),
		field.String("name"),
	}
}

const(
	EDGE_DEPT_TO_USERS = "users"
	EDGE_DEPT_TO_USER_PROPERTIES_IN_DEPT = "user_properties_in_dept"
)

// Edges of the Dept.
func (Dept) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To(EDGE_DEPT_TO_USERS,User.Type),
		edge.To(EDGE_DEPT_TO_USER_PROPERTIES_IN_DEPT,UserPropertyInDept.Type),
	}
}

func (Dept) Hooks()[]ent.Hook{
	return []ent.Hook{
		hook.On(
			func (next ent.Mutator)ent.Mutator{
				return hook.DeptFunc(
					func (ctx context.Context,m *ent_gen.DeptMutation)(ent.Value,error){
						fmt.Println("hoos triger")
						// hook for when dept delete 
												deptId,exists:= m.ID()
												if !exists{
													return nil,fmt.Errorf("Only allow delete dept by id")
												}
						tx,err := m.Tx()
						if err != nil{
							return nil,fmt.Errorf("Extract tx in dept hook onDelete failed: %w",err)
						}
						// < delete user in the dept, if the user only in this dept
						if err != nil{
							return nil,err
						}
						fmt.Println("to delete user who not in dept")
						_,err = tx.Dept.UpdateOneID(deptId).ClearUsers().Save(ctx)
						if err != nil{
							return nil,err
						}
						//						usersInDept,err := tx.User.Query().Where(user.HasDeptsWith(dept.IDEQ(deptId))).Exec(ctx)
						//						if err != nil{
						//							t.
						//						}
						usersToDelete,err := tx.User.Query().Where(user.Not(user.HasDepts())).All(ctx)
						if err != nil{
							return nil,err
						}
						for _,v := range usersToDelete{
							err = tx.User.DeleteOneID(v.ID).Exec(ctx)
							if err != nil{
								return nil,err
							}
						}
						// >
						_,err = tx.UserPropertyInDept.Delete().Where(up.DeptIDEQ(deptId)).Exec(ctx)
						if err != nil{
							return nil,err
						}
						return next.Mutate(ctx,m)
					},
				)
			},
			ent.OpDelete|ent.OpDeleteOne,

		),
	}
}
