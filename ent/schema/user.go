package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Indexes of the User
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("member_admin"),
		index.Edges("lead_admin"),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("age"),
		field.String("name"),
		field.Int("member_admin_id").
			Optional(),
		field.Int("lead_admin_id").
			Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("member_admin", Admin.Type).
			Field("member_admin_id").
			Ref("team_members").Unique(),
		edge.From("lead_admin", Admin.Type).
			Field("lead_admin_id").
			Ref("team_leader").Unique(),
	}
}
