package schema

import (
"entgo.io/ent"
"entgo.io/ent/schema/field"
"entgo.io/ent/schema/edge"
)

// UserPropertyInDept holds the schema definition for the UserPropertyInDept entity.
type UserPropertyInDept struct {
	ent.Schema
}

const(
	FIELD_USER_PROPERTIES_IN_DEPT_USER_ID ="user_id"
	FIELD_USER_PROPERTIES_IN_DEPT_DEPT_ID ="dept_id"
)

// Fields of the UserPropertyInDept.
func (UserPropertyInDept) Fields() []ent.Field {
	return []ent.Field{
		field.String(FIELD_USER_PROPERTIES_IN_DEPT_USER_ID).Optional(),
		field.Uint(FIELD_USER_PROPERTIES_IN_DEPT_DEPT_ID).Optional(),
		field.Bool("isLeader"),
	}
}

// Edges of the UserPropertyInDept.
func (UserPropertyInDept) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user",User.Type).Ref(EDGE_USER_TO_USER_PROPERTIES_IN_DEPT ).Field(FIELD_USER_PROPERTIES_IN_DEPT_USER_ID).Unique(),
		edge.From("dept",Dept.Type).Ref(EDGE_DEPT_TO_USER_PROPERTIES_IN_DEPT ).Field(FIELD_USER_PROPERTIES_IN_DEPT_DEPT_ID).Unique(),
	}
}

