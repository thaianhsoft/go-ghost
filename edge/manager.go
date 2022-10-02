package edge

import (
	"fmt"
	"github.com/thaianhsoft/go-ghost/field"
	"reflect"
	"strings"
	"sync"
)

var manager *edgeManager = &edgeManager{
	group: &sync.WaitGroup{},
	locker: &sync.Mutex{},
	edges:  map[string]*edge{},
	insertQueries: map[string]string{}, // map insert query for each table
}

func DefaultManager() *edgeManager {
	return manager
}

type edgeManager struct {
	locker        *sync.Mutex
	edges         map[string]*edge
	group         *sync.WaitGroup
	insertQueries map[string]string
}

func (e *edgeManager) getEdge(edgeName, toTableName string) *edge {
	e.locker.Lock()
	defer e.locker.Unlock()
	encodeEdge := edgeName + stringJoinRelKey + toTableName
	if _, ok := e.edges[encodeEdge]; ok {
		return e.edges[encodeEdge]
	}
	return nil
}

func (e *edgeManager) addEdge(edgeName, toTableName string, edgeClass *edge) {
	e.locker.Lock()
	defer e.locker.Unlock()
	encodeEdge := edgeName + stringJoinRelKey + toTableName
	if _, ok := e.edges[encodeEdge]; !ok {
		e.edges[encodeEdge] = edgeClass
	}
}

func GetEdges() *map[string]*edge {
	return &manager.edges
}

func LoopEdges(fn func(edgeName string, e *edge)) {
	for edgeName, e := range manager.edges {
		manager.locker.Lock()
		fn(edgeName, e)
		manager.locker.Unlock()
	}
}

func (e *edgeManager) Migrate(schemas ...ISchema) *string {
	sqlCreateTableStmt := ``
	for _, schema := range schemas {
		schema.DefineEdges()
	}
	e.getGroup().Wait()
	colsOfTableAdded := map[string]map[string]field.IField{}
	for _, schema := range schemas {
		rv := reflect.Indirect(reflect.ValueOf(schema))
		tbName := rv.Type().Name()
		pkStmt := "PRIMARY KEY"
		if _, ok := colsOfTableAdded[tbName]; !ok {
			colsOfTableAdded[tbName] = make(map[string]field.IField)
		}
		sqlCreateTableStmt += fmt.Sprintf("CREATE TABLE `%v`(\n", rv.Type().Name())
		for i, schemaField := range schema.DefineSchema() {
			sqlCreateTableStmt += fmt.Sprintf("%*s`%v` %v", 4, "", schemaField.GetColName(), schemaField.GetSqlTypeAndOptions())
			if strings.Contains(schemaField.GetSqlTypeAndOptions(), "AUTO_INCREMENT") {
				pkStmt += fmt.Sprintf(" (`%v`)", schemaField.GetColName())
			}
			if i < len(schema.DefineSchema()) - 1 {
				sqlCreateTableStmt += ", \n"
			}
			rvSchemaField := reflect.Indirect(reflect.ValueOf(schemaField))
			if _, ok := rvSchemaField.Type().FieldByName("name"); ok {
				colName := rvSchemaField.FieldByName("name").String()
				if _, ok := colsOfTableAdded[tbName][colName]; !ok {
					colsOfTableAdded[tbName][colName] = schemaField
				}
			}
		}
		if pkStmt != "PRIMARY KEY" {
			sqlCreateTableStmt += fmt.Sprintf(", %v", pkStmt)
		}
		sqlCreateTableStmt += ");\n"
	}

	edges := GetEdges()
	fmt.Println(*edges)
	LoopEdges(func(edgeName string, e *edge) {
		fmt.Println("self edge: ", e, " ref edge: ", e.refEdge)
		if e.unique && !e.refEdge.unique{
			fmt.Println("find one alter table")
			// foreign key on e
			foreignTable := e.refEdge.toTable
			foreignKey := e.keyHere
			pkKeyRef := e.refEdge.keyHere
			pkRefTable := e.toTable
			//nameConstraintFK := edgeName
			// check fk and pk exist in script sql
			if field, ok := colsOfTableAdded[pkRefTable][pkKeyRef]; ok {

				q := field.GetSqlTypeAndOptions()
				maxIndex := len("AUTO_INCREMENT")
				for i, char := range q {
					if string(char) == "A" &&  i + maxIndex - 1< len(q){
						fmt.Println("find A", string(char), i)
						if string(q[i:i+maxIndex]) == "AUTO_INCREMENT" {
							fmt.Println("find auto increment")
							q = q[0:i-1] + q[i+maxIndex:] // trim space of i-1
						}
					}
				}
				sqlCreateTableStmt += fmt.Sprintf("ALTER TABLE `%v` ADD COLUMN `%v` %v", foreignTable, foreignKey, string(q))
			}
			sqlCreateTableStmt += fmt.Sprintf(", ADD CONSTRAINT `%v` FOREIGN KEY(`%v`) REFERENCES `%v`(`%v`) ON DELETE CASCADE ON UPDATE CASCADE;\n", edgeName, foreignKey, pkRefTable, pkKeyRef)
			delete(*edges, e.edgeName + stringJoinRelKey + e.toTable)
			delete(*edges, e.refEdge.edgeName + stringJoinRelKey + e.refEdge.toTable)
		}
		if !e.unique && !e.refEdge.unique {
			// many to many create more table
		}
	})

	return &sqlCreateTableStmt
}


func (e *edgeManager) getGroup() *sync.WaitGroup {
	return e.group
}


