package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// UserSession holds the schema definition for the UserSession entity.
type UserSession struct {
	ent.Schema
}

// Fields of the UserSession.
func (UserSession) Fields() []ent.Field {
	return []ent.Field{
		field.String("session_id").Unique(),
		field.Int("user_id"),
		field.Time("login_time"),
		field.Time("last_seen_time"),
	}
}

// Edges of the UserSession.
func (UserSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Unique().Field("user_id").Required(),
	}
}
