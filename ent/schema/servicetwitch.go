package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// ServiceTwitch holds the schema definition for the ServiceTwitch entity.
type ServiceTwitch struct {
	ent.Schema
}

// Fields of the ServiceTwitch.
func (ServiceTwitch) Fields() []ent.Field {
	return []ent.Field{
		field.String("access_token"),
		field.String("refresh_token"),
		field.Strings("scopes"),
	}
}

// Edges of the ServiceTwitch.
func (ServiceTwitch) Edges() []ent.Edge {
	return nil
}
