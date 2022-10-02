package edge

import (
	"fmt"
	"github.com/thaianhsoft/go-ghost/field"
	"github.com/thaianhsoft/go-ghost/mockdb"
	"github.com/thaianhsoft/go-ghost/schema"
	"testing"
)

type Book struct {
	*schema.Schema
}

func (b *Book) DefineSchema() []field.IField {
	return []field.IField{
		field.UIntType("Id").AI(),
		field.VarCharType("BookName", 40).Nullable(true),
	}
}

func (b *Book) DefineEdges() []IEdge {
	return []IEdge{
		PointTo("ContainNotes", &Note{}),
		PointBack("Owner", &Student{}).RefOn("HasBooks", &Book{}).Unique(),
	}
}

type Note struct {
	*schema.Schema
}


func (n *Note) DefineSchema() []field.IField {
	return []field.IField{
		field.UIntType("Id").AI(),
		field.JSONType("Content").Nullable(true),
	}
}

func (n *Note) DefineEdges() []IEdge {
	return []IEdge{
		PointBack("PartOfBook", &Book{}).RefOn("ContainNotes", &Note{}).Unique(),
	}
}

type Student struct {
	*schema.Schema
}

func (s *Student) DefineSchema() ([]field.IField) {
	return []field.IField{
		field.UIntType("Id").AI(),
		field.VarCharType("Gmail", 30).Unique().Nullable(false).Default("thaianhsoftmail.com"),
		field.VarCharType("Password", 30).Nullable(false),
		field.CharType("IsOnline", 1).Nullable(true).Default("0"),
	}
}

func (s *Student) DefineEdges() []IEdge {
	return []IEdge{
		PointTo("HasBooks", &Book{}),
	}
}

func TestEdge(t *testing.T) {
	query := manager.Migrate(&Note{}, &Student{}, &Book{})
	db := mockdb.OpenDB()
	fmt.Println(*query)
	if _, err := db.Query(*query); err == nil {
		fmt.Println("create tables successfully!")
	} else {
		fmt.Println(err)
	}
}

