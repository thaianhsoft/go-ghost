package changeset

import (
	"fmt"
	"github.com/thaianhsoft/go-ghost/edge"
	"github.com/thaianhsoft/go-ghost/field"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync/atomic"
)

var manager *changesetManager = &changesetManager{
	reflectSchemas: nil,
}
type tableName string
type changesetManager struct {
	casFreeLock uint32
	reflectSchemas map[tableName]*reflect.Value // one changeset contain one reflect schema value
}

func (m *changesetManager) TrySaveCS(name string, v *reflect.Value) {
	if _, ok := m.reflectSchemas[tableName(name)]; !ok {
		m.reflectSchemas[tableName(name)] = v
	}
}

func (m *changesetManager) GetRvSchema(name string) *reflect.Value {
	for !atomic.CompareAndSwapUint32(&m.casFreeLock, 0, 1){
		runtime.Gosched()
	}
	defer func() {
		fmt.Println("release free lock cas")
		atomic.StoreUint32(&m.casFreeLock, 0)
	}()
	if _, ok := m.reflectSchemas[tableName(name)]; ok {
		newCS := reflect.New((*m.reflectSchemas[tableName(name)]).Type())
		return &newCS
	}
	return nil
}

func (m *changesetManager) GetChangeSetManager() *changesetManager {
	return m
}
type boxFieldCS struct {
	colName string
	val interface{}
	indexNotNull int
}

type validateErr uint8

func (v validateErr) toErrString(fieldName string, need interface{}, have interface{}) string {
	return fmt.Sprintf("[%v]: Validate on Field [%v] failed, need [%v] but have [%v]\n", errTypeToString[v], fieldName, need, have)
}

const (
	validateRequiredErr validateErr = iota
	validatePatternErr
	validateUniqueErr
)

var errTypeToString = [...]string{"RequiredError", "PatternError", "UniqueError"}

type ChangeSet interface{
	ValidateRequired(fieldNames ...string) ChangeSet
	ValidPattern(fieldName string, pattern []byte) ChangeSet
	ValidUnique(fieldName string) ChangeSet
	IsValidAll() error
}

type changeSetInternal struct {
	boxes          map[string]*boxFieldCS
	orderValidFunc []func() string
	notNullIndex   uint32
}

func (c *changeSetInternal) IsValidAll() error {
	allErrString := ""
	for _, validateFunc := range c.orderValidFunc {
		allErrString += validateFunc()
	}
	if allErrString == "" {
		return nil
	}
	return fmt.Errorf(allErrString)
}

func (c *changeSetInternal) ValidateRequired(fieldNames ...string) ChangeSet {
	c.orderValidFunc = append(c.orderValidFunc, func() string {
		if c.notNullIndex == 0 {
			return ""
		} else {
			var errString string = ""
			if len(fieldNames) == 0 {
				// default validate not null all field

				for _, fieldName := range fieldNames {
					if _, ok := c.boxes[fieldName]; ok {
						if (c.notNullIndex & (1 << c.boxes[fieldName].indexNotNull)) != 0 {
							errString += validateRequiredErr.toErrString(fieldName, "Not Null", "Null")
						}
					}
				}
			} else {
				for fieldName, _ := range c.boxes {
					if (c.notNullIndex & (1 << c.boxes[fieldName].indexNotNull)) != 0 {
						errString += validateRequiredErr.toErrString(fieldName, "Not Null", "Null")
					}
				}
			}
			if errString != "" {
				return errString
			} else {
				return ""
			}
		}
	})
	return c
}

func (c *changeSetInternal) ValidPattern(fieldName string, pattern []byte) ChangeSet {
	c.orderValidFunc = append(c.orderValidFunc, func() string {
		r, err := regexp.Compile(string(pattern))
		if err == nil {
			if _, ok := c.boxes[fieldName]; ok {
				fmt.Println(c.boxes[fieldName])
				if _, ok := c.boxes[fieldName].val.(string); ok {
					if r.MatchString(c.boxes[fieldName].val.(string)) {
						return ""
					} else {
						return validatePatternErr.toErrString(fieldName, string(pattern), c.boxes[fieldName].val.(string))
					}
				}
			}
		}
		return validatePatternErr.toErrString(fieldName, "pattern regex valid", string(pattern) + " not valid")
	})
	return c
}

func (c *changeSetInternal) ValidUnique(fieldName string) ChangeSet {
	return c
}

func CastIntoSchemaClass(fromMessage interface{}, toSchemaClass edge.ISchema) ChangeSet {
	rvMessage := reflect.Indirect(reflect.ValueOf(fromMessage))
	rvSchema := reflect.Indirect(reflect.ValueOf(toSchemaClass))
	prefixSchemaName := rvSchema.Type().Name()
	indexBox := 0
	cs := &changeSetInternal{
		boxes:                 map[string]*boxFieldCS{},
		notNullIndex : 0,
	}
	for indexBox < len(toSchemaClass.DefineSchema()) {
		fieldSchema := toSchemaClass.DefineSchema()[indexBox]
		switch f := fieldSchema.(type) {
		case *field.Field:
			if _, ok := cs.boxes[f.GetFieldName()]; !ok {
				newBox := &boxFieldCS{
					colName:      f.GetFieldName(),
					val:          nil,
					indexNotNull: indexBox,
				}
				cs.boxes[f.GetFieldName()] = newBox
				if !f.CanNull() {
					cs.notNullIndex |= 1 << indexBox
				}
				indexBox++
			}
		}
	}
	manager.TrySaveCS(prefixSchemaName, &rvSchema)
	switch rvMessage.Kind() {
	case reflect.Struct:
		for i := 0; i < rvMessage.NumField(); i++ {
			fieldMessageNameHaveSchemaName := rvMessage.Type().Field(i).Name
			rvfieldMessage := rvMessage.Field(i)
			fmt.Println(fieldMessageNameHaveSchemaName, prefixSchemaName)
			if strings.Contains(fieldMessageNameHaveSchemaName, prefixSchemaName) {
				fieldNameSkipPrefixSchemaName := fieldMessageNameHaveSchemaName[len(prefixSchemaName):]
				if _, ok := rvSchema.Type().FieldByName(fieldNameSkipPrefixSchemaName); ok  {
					// exist
					fmt.Println(rvfieldMessage.Type().Name(), rvSchema.FieldByName(fieldNameSkipPrefixSchemaName).Type().Name())
					if rvfieldMessage.Type().Name() == rvSchema.FieldByName(fieldNameSkipPrefixSchemaName).Type().Name() {
						// check same type
						fmt.Println(rvfieldMessage.Type().Name())
						if !rvfieldMessage.IsZero() {
							// check field have value from message exlude [0, "", nil]
							if _, ok := cs.boxes[fieldNameSkipPrefixSchemaName]; ok {
								cs.boxes[fieldNameSkipPrefixSchemaName].val = rvfieldMessage.Interface()
								cs.notNullIndex &= ^(1 << cs.boxes[fieldNameSkipPrefixSchemaName].indexNotNull) // clear
							}
						}
					}
				}
			}
		}
	}
	return cs
}


/*
func (c *changeSetInternal) castRelationEmbedded(schemas ...edge.ISchema) {
	for _, schemaRel := range schemas {
		schemaName := reflect.Indirect(reflect.ValueOf(schemaRel)).Type().Name()

	}
}
 */


// encode err for type validate
// validsErrField => bit mask 32 bit equivalent 32 field index
// []valueFieldErr => each slot is one index equivalent fieldName