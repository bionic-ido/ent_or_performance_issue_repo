// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AdminsColumns holds the columns for the "admins" table.
	AdminsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "age", Type: field.TypeInt},
		{Name: "name", Type: field.TypeString},
	}
	// AdminsTable holds the schema information for the "admins" table.
	AdminsTable = &schema.Table{
		Name:       "admins",
		Columns:    AdminsColumns,
		PrimaryKey: []*schema.Column{AdminsColumns[0]},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "age", Type: field.TypeInt},
		{Name: "name", Type: field.TypeString},
		{Name: "member_admin_id", Type: field.TypeInt, Nullable: true},
		{Name: "lead_admin_id", Type: field.TypeInt, Unique: true, Nullable: true},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "users_admins_team_members",
				Columns:    []*schema.Column{UsersColumns[3]},
				RefColumns: []*schema.Column{AdminsColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "users_admins_team_leader",
				Columns:    []*schema.Column{UsersColumns[4]},
				RefColumns: []*schema.Column{AdminsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "user_member_admin_id",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[3]},
			},
			{
				Name:    "user_lead_admin_id",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[4]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AdminsTable,
		UsersTable,
	}
)

func init() {
	UsersTable.ForeignKeys[0].RefTable = AdminsTable
	UsersTable.ForeignKeys[1].RefTable = AdminsTable
}
