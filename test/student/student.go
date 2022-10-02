package student

import (
	"github.com/thaianhsoft/go-ghost/edge"
	"github.com/thaianhsoft/go-ghost/field"
	"github.com/thaianhsoft/go-ghost/schema"
	"github.com/thaianhsoft/go-ghost/test/schemas"
)

type Student struct {
	*schema.Schema
}
func (s *Student) DefineSchema() []field.IField {
	return []field.IField{
		field.VarCharType("Id", 40).Unique(),
		field.VarCharType("Gmail", 30).Unique().Nullable(false).Default("thaianhsoftmail.com"),
		field.VarCharType("Password", 30).Nullable(false),
		field.CharType("IsOnline", 1).Nullable(true).Default("0"),
	}
}

func (s *Student) DefineEdges() []edge.IEdge {
	return []edge.IEdge{
		edge.PointTo("HasBooks", schemas.BookType),
	}
}
