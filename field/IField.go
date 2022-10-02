package field

import (
	"fmt"
	"reflect"
)
var maxBitSizeOption = 8
type fieldType uint32
const (
	intFieldType fieldType = iota
	uIntFieldType
	varcharFieldType
	charFieldType
	jsonFieldType
)

var mapFieldTypes = [...]string{"INT", "INT UNSIGNED", "VARCHAR", "CHAR", "JSON"}
type IField interface{
	GetColName() string
	GetSqlTypeAndOptions() string
	Unique() IField
	AI() IField
	Nullable(v bool) IField
	Default(v interface{}) IField
}

type option uint8
const (
	nullableOp = iota
	notNullableOp
	uniqueOp
	autoIncrementOp
	defaultOp
)

var mapFieldOps =  [...]string{
	"NULL",
	"NOT NULL",
	"UNIQUE",
	"AUTO_INCREMENT",
}


type Field struct {
	fType fieldType
	ops option
	name string
	size int
	pk bool
	defaultValue interface{}
}

func (f *Field) PrimaryKey() IField {
	f.pk = true
	return f
}

func (f *Field) Default(v interface{}) IField {
	rv := reflect.Indirect(reflect.ValueOf(v))
	rvDefault := reflect.Indirect(reflect.ValueOf(f.defaultValue))
	fmt.Println(rv.Kind(), rvDefault.Kind())
	if rv.Kind() == rvDefault.Kind() {
		if f.fType == varcharFieldType || f.fType == charFieldType {
			if len(v.(string)) <= f.size {
				f.defaultValue = v
				return f
			}
		} else {
			f.defaultValue = v
		}
	}
	return f
}

func (f *Field) Unique() IField {
	if !f.ops.CheckOp(autoIncrementOp) {
		f.ops.setOps(uniqueOp)
	}
	return f
}

func (f *Field) AI() IField {
	if !f.ops.CheckOp(autoIncrementOp) {
		f.ops.setOps(autoIncrementOp)
		f.ops.setOps(notNullableOp)
	}
	return f
}

func (f *Field) Nullable(v bool) IField {
	if f.ops.CheckOp(autoIncrementOp) {
		return f
	}
	if !v {
		f.ops.setOps(notNullableOp)
	} else {
		f.ops.setOps(nullableOp)
	}
	return f
}


func (o *option) setOps(ops ...option) {
	for _, op := range ops {
		*o |= 1 << op
	}
}

func (o *option) CheckOp(haveOp option) bool {
	return ((*o) & (1<<haveOp)) != 0
}


func IntType(fieldName string) IField {
	return &Field{
		fType: uIntFieldType,
		name: fieldName,
	}
}

func UIntType(fieldName string) IField {
	return &Field{
		defaultValue: uint32(0),
		fType: uIntFieldType,
		ops: 0,
		name: fieldName,
	}
}

func BigIntType(fieldName string) IField {
	return &Field{
		defaultValue: int64(0),
		ops:  0,
		name: fieldName,
	}
}


func BigUIntType(fieldName string) IField {
	return &Field{
		defaultValue: uint64(0),
		name: fieldName,
		ops: 0,
	}
}

func VarCharType(fieldName string, size int) IField {
	return &Field{
		defaultValue: "",
		fType: varcharFieldType,
		size: size,
		name: fieldName,
	}
}

func CharType(fieldName string, size int) IField {
	return &Field{
		defaultValue: "",
		fType: charFieldType,
		size: size,
		name: fieldName,
	}
}

func JSONType(fieldName string) IField {
	return &Field{
		defaultValue: "",
		fType: jsonFieldType,
		name: fieldName,
	}
}


func (f *Field) GetColName() string {
	return f.name
}

func (f *Field) GetSqlTypeAndOptions() string {
	createStmt := fmt.Sprintf("%v", mapFieldTypes[f.fType])
	if f.fType == varcharFieldType || f.fType == charFieldType {
		createStmt += fmt.Sprintf("(%v)", f.size)
	}
	for i := 0; i < maxBitSizeOption; i++ {
		if f.ops.CheckOp(option(i)) {
			createStmt += " " + mapFieldOps[i]
		}
	}
	rvDefault := reflect.Indirect(reflect.ValueOf(f.defaultValue))
	if !rvDefault.IsZero() {
		if rvDefault.Kind() == reflect.String {
			createStmt += " DEFAULT "
			if f.fType == jsonFieldType {
				createStmt += fmt.Sprintf("('%v')", f.defaultValue)
			} else {
				createStmt += fmt.Sprintf("'%v'", f.defaultValue)
			}

		} else {
			createStmt += fmt.Sprintf(" DEFAULT %v", f.defaultValue)
		}
	}
	return createStmt
}

func (f *Field) CanNull() bool {
	if f.ops.CheckOp(autoIncrementOp) {
		return true
	}
	if f.ops.CheckOp(notNullableOp) {
		return false
	}
	return true
}

func (f *Field) GetFieldName() string {
	return f.name
}
