package field

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"testing"
)




type Place struct {
	Lat float32
	Long float32
}
type Student struct {

}

var v = "{\"Lat\": 3.23, \"Long\": 1.28}"
func (s *Student) DefineSchema() []IField {
	return []IField{
		UIntType("Id").AI().Nullable(false),
		VarCharType("Gmail", 30).Unique().Nullable(false).Default("thaianhsoftmail.com"),
		VarCharType("Password", 30).Nullable(false),
		CharType("IsOnline", 1).Nullable(true).Default("0"),
		JSONType("Content").Nullable(false).Default(v),
	}
}

func OpenDb(config *mysql.Config) *sql.DB {
	if db, err := sql.Open("mysql", config.FormatDSN()); err == nil {
		return db
	} else {
		return nil
	}
}

func TestCreateField(t *testing.T) {
	db := OpenDb(mysqlConfig)
	if db != nil {
		fmt.Println("open datbase ok!!!")
		s := &Student{}
		rv := reflect.Indirect(reflect.ValueOf(s))
		schemaName := rv.Type().Name()
		createStmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%v` (\n", schemaName)
		for i, field := range s.DefineSchema() {
			createStmt += fmt.Sprintf("%*s%v",4, "",field.ToSqlCreateStmt())
			if i < len(s.DefineSchema()) - 1 {
				createStmt += ", \n"
			}
		}
		createStmt += "\n)"
		if rows, err := db.Query(createStmt); err == nil {
			fmt.Println(rows)
		} else {
			fmt.Println("error create table: ", err)
		}
	}
}
