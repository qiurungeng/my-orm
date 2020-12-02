package clause

import (
	"fmt"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func testSelect(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"*"})
	clause.Set(WHERE, "Name = ?", "Tom")
	clause.Set(ORDERBY, "Age ASC")
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func testInsert(t *testing.T) {
	var clause Clause
	clause.Set(INSERT, "User", []string{"Name", "Age"})
	//clause.Set(VALUES, []interface{}{User{"Tom", 11}, User{"Jerry", 112}})
	clause.Set(VALUES, []interface{}{"Tom", 1},[]interface{}{"Jerry", 2})
	sql, vars := clause.Build(INSERT, VALUES)
	t.Log(sql, vars)
}

func TestClause_Build(t *testing.T) {
		t.Run("select", func(t *testing.T) {
		testSelect(t)
		testInsert(t)
	})
}

func TestUpdate(t *testing.T)  {
	fmt.Println(_update("User", map[string]interface{}{
		"Name": "Tom",
		"Age": 13,
	}))
}