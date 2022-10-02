package changeset

import (
	"fmt"
	"github.com/thaianhsoft/go-ghost/edge"
	"github.com/thaianhsoft/go-ghost/field"
	"github.com/thaianhsoft/go-ghost/schema"
	"testing"
)

type Place struct {
	Lat float32
	Lon float32
}
type MessageCast struct {
	TestEntityAddress string
	TestEntityPhone string
	TestEntityJSONGPS *Place
	TestEntityStartTS uint32
}


type TestEntity struct {
	*schema.Schema
	Id uint32
	Address string
	Phone string
	JSONGPS *Place
	StartTS uint32
}

func (t TestEntity) DefineSchema() []field.IField {
	return []field.IField{
		field.UIntType("Id").AI(),
		field.VarCharType("Address", 40).Nullable(false),
		field.VarCharType("Phone", 10).Nullable(false),
		field.JSONType("JSONGPS").Nullable(false),
		field.JSONType("StartTS").Nullable(false),
	}
}

func (t TestEntity) DefineEdges() []edge.IEdge {
	return nil
}

func (t TestEntity) Migrate() {

}

func TestChangeSet(t *testing.T) {
	msg := &MessageCast{
		TestEntityAddress: "A0(@gmail.com",
		TestEntityStartTS: 0,
	}
	/*msg.TestEntityJSONGPS = &Place{
		Lat: 15.42,
		Lon: 12.43,
	}*/
	cs := CastIntoSchemaClass(msg, &TestEntity{})
	err := cs.ValidateRequired("Address", "Phone").ValidPattern("Address", []byte(`[A-Z]+[0-9]+(\(|\*)+.*@gmail.com$`)).IsValidAll()
	fmt.Println(cs.(*changeSetInternal).notNullIndex)
	fmt.Println(err)
}


