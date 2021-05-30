// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// DeptsColumns holds the columns for the "depts" table.
	DeptsColumns = []*schema.Column{
		{Name: "dept_id", Type: field.TypeUint, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "generation", Type: field.TypeUint, Nullable: true},
		{Name: "dept_sub_depts", Type: field.TypeUint, Nullable: true},
	}
	// DeptsTable holds the schema information for the "depts" table.
	DeptsTable = &schema.Table{
		Name:       "depts",
		Columns:    DeptsColumns,
		PrimaryKey: []*schema.Column{DeptsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "depts_depts_sub_depts",
				Columns:    []*schema.Column{DeptsColumns[3]},
				RefColumns: []*schema.Column{DeptsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "user_id", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "phone", Type: field.TypeString},
		{Name: "generation", Type: field.TypeUint, Nullable: true},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:        "users",
		Columns:     UsersColumns,
		PrimaryKey:  []*schema.Column{UsersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// UserPropertyInDeptsColumns holds the columns for the "user_property_in_depts" table.
	UserPropertyInDeptsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "is_leader", Type: field.TypeBool},
		{Name: "dept_id", Type: field.TypeUint, Nullable: true},
		{Name: "user_id", Type: field.TypeString, Nullable: true},
	}
	// UserPropertyInDeptsTable holds the schema information for the "user_property_in_depts" table.
	UserPropertyInDeptsTable = &schema.Table{
		Name:       "user_property_in_depts",
		Columns:    UserPropertyInDeptsColumns,
		PrimaryKey: []*schema.Column{UserPropertyInDeptsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "user_property_in_depts_depts_user_properties_in_dept",
				Columns:    []*schema.Column{UserPropertyInDeptsColumns[2]},
				RefColumns: []*schema.Column{DeptsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "user_property_in_depts_users_properties_in_dept",
				Columns:    []*schema.Column{UserPropertyInDeptsColumns[3]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// DeptUsersColumns holds the columns for the "dept_users" table.
	DeptUsersColumns = []*schema.Column{
		{Name: "dept_id", Type: field.TypeUint},
		{Name: "user_id", Type: field.TypeString},
	}
	// DeptUsersTable holds the schema information for the "dept_users" table.
	DeptUsersTable = &schema.Table{
		Name:       "dept_users",
		Columns:    DeptUsersColumns,
		PrimaryKey: []*schema.Column{DeptUsersColumns[0], DeptUsersColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "dept_users_dept_id",
				Columns:    []*schema.Column{DeptUsersColumns[0]},
				RefColumns: []*schema.Column{DeptsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "dept_users_user_id",
				Columns:    []*schema.Column{DeptUsersColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		DeptsTable,
		UsersTable,
		UserPropertyInDeptsTable,
		DeptUsersTable,
	}
)

func init() {
	DeptsTable.ForeignKeys[0].RefTable = DeptsTable
	UserPropertyInDeptsTable.ForeignKeys[0].RefTable = DeptsTable
	UserPropertyInDeptsTable.ForeignKeys[1].RefTable = UsersTable
	DeptUsersTable.ForeignKeys[0].RefTable = DeptsTable
	DeptUsersTable.ForeignKeys[1].RefTable = UsersTable
}
