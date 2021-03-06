package myorm

import (
	"errors"
	"myorm/session"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

type User struct {
	Name string `myorm:"PRIMARY KEY"`
	Age  int
}

func TestTransaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		txRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		txCommit(t)
	})
}

func txRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.SetRefTableByModel(&User{}).DropTable()

	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.SetRefTableByModel(&User{}).CreateTable()
		_, err = s.Insert(&User{Name: "Tim", Age: 18})
		return nil, errors.New("my own test error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("fail to rollback")
	}
}

func txCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.SetRefTableByModel(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.SetRefTableByModel(&User{}).CreateTable()
		_, err = s.Insert(&User{Name: "Tom", Age: 22})
		return
	})
	u := &User{}
	err = s.FindFirst(u)
	if err != nil || (u.Name != "Tom" && u.Age != 22) {
		t.Fatal("fail to commit")
	}
}


/*************************** Migrate Test ****************************/
func TestEngine_Migrate(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	_ = engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
