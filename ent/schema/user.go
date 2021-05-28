package schema

import (
	"context"
	"fmt"
"entgo.io/ent"
"github.com/Peanuttown/dd_contacts/ent/hook"
up "github.com/Peanuttown/dd_contacts/ent/userpropertyindept"
gen_ent "github.com/Peanuttown/dd_contacts/ent"
"entgo.io/ent/schema/field"
"entgo.io/ent/schema/edge"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("user_id"),
		field.String("name"),
		field.String("phone"),
	}
}

const(
	EDGE_USER_TO_USER_PROPERTIES_IN_DEPT ="properties_in_dept"
)

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("depts",Dept.Type).Ref(EDGE_DEPT_TO_USERS),
		edge.To(EDGE_USER_TO_USER_PROPERTIES_IN_DEPT,UserPropertyInDept.Type),
	}
}

func (User) Hooks()[]ent.Hook{
	return []ent.Hook{
		hook.On(
			func (next ent.Mutator)ent.Mutator{
				return hook.UserFunc(
					func(ctx context.Context,m *gen_ent.UserMutation)(ent.Value,error){
						userId,exists := m.ID()
						if !exists{
							return nil,fmt.Errorf("Only Allow deleted user by user_id")
						}
						tx,err := m.Tx()
						if err != nil{
							return nil,err
						}
						_,err = tx.UserPropertyInDept.Delete().Where(up.UserIDEQ(userId)).Exec(ctx)
						if err != nil{
							return nil,err
						}
						return next.Mutate(ctx,m)

					},
				)
			},
			ent.OpDeleteOne | ent.OpDelete ,
		),
	}
}
