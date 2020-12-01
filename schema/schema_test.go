package schema

import (
	"fmt"
	"myorm/dialect"
	"testing"
)

type User struct {
	Name string `myorm:"PRIMARY KEY"`
	Age int
}

var TestDialect, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&User{}, TestDialect)
	if schema.Name != "User" || len(schema.Fields) != 2{
		fmt.Println(schema.Name)
		fmt.Println(len(schema.Fields))
		t.Fatal("fail to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("fail to parse PRIMARY KEY tag")
	}
}