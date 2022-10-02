package book

import (
	"github.com/thaianhsoft/go-ghost/edge"
	"github.com/thaianhsoft/go-ghost/field"
	"github.com/thaianhsoft/go-ghost/schema"
	"github.com/thaianhsoft/go-ghost/test/schemas"
)

type Book struct {
	*schema.Schema
}

func (b *Book) DefineSchema() []field.IField {
	return []field.IField{
		field.VarCharType("Id", 20).Unique(),
		field.VarCharType("BookName", 40).Nullable(true),
	}
}

func (b *Book) DefineEdges() []edge.IEdge {
	return []edge.IEdge{
		edge.PointTo("ContainNotes", schemas.NoteType),
		edge.PointBack("Owner", schemas.StudentType).RefOn("HasBooks", schemas.BookType).Unique(),
	}
}

