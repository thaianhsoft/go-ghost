package edge

import (
	"github.com/thaianhsoft/go-ghost/field"
	"runtime"
)

var stringJoinRelKey = "Rel"
type ISchema interface {
	DefineSchema() []field.IField
	DefineEdges() []IEdge
	Migrate()
}

type EdgeName string
type IEdge interface{
	Unique() IEdge
	RefOn(refEdgeName string, refSchema EdgeName) IEdge
}

type edge struct {
	edgeName string
	unique  bool
	keyHere string
	toTable string
	refEdge *edge
}

func (e *edge) Unique() IEdge {
	e.unique = true
	return e
}

func (e *edge) RefOn(refEdgeName string, refSchema EdgeName) IEdge {
	defer func() {
		go func() {
			defer manager.getGroup().Done()
			for manager.getEdge(refEdgeName, string(refSchema)) == nil {
				runtime.Gosched()
			}
			refEdge := manager.getEdge(refEdgeName, string(refSchema))
			e.refEdge, refEdge.refEdge = refEdge, e

			e.keyHere = e.edgeName + stringJoinRelKey + "Id"
			refEdge.keyHere = "Id"
			return
		}()
	}()
	return e
}

func PointTo(edgeName string, schemaClassName EdgeName) *edge {
	e := &edge{
		edgeName: edgeName,
		unique:  false,
		keyHere: "",
		toTable: string(schemaClassName),
	}
	manager.addEdge(edgeName, string(schemaClassName), e)
	return e
}


func PointBack(edgeName string, schemaClassName EdgeName) *edge {
	manager.group.Add(1)
	e := &edge{
		edgeName: edgeName,
		unique:  false,
		keyHere: "",
		toTable: string(schemaClassName),
	}
	manager.addEdge(edgeName, string(schemaClassName), e)
	return e
}

