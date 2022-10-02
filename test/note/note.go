package note

import (
	"github.com/thaianhsoft/go-ghost/edge"
	"github.com/thaianhsoft/go-ghost/field"
	"github.com/thaianhsoft/go-ghost/schema"
	"github.com/thaianhsoft/go-ghost/test/schemas"
)

type Note struct {
	*schema.Schema
}


func (n *Note) DefineSchema() []field.IField {
	return []field.IField{
		field.VarCharType("Id", 20).Unique(),
		field.JSONType("Content").Nullable(true),
	}
}

func (n *Note) DefineEdges() []edge.IEdge {
	return []edge.IEdge{
		edge.PointBack("PartOfBook", schemas.BookType).RefOn("ContainNotes", schemas.NoteType).Unique(),
	}
}
