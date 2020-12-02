package model

import (
	"database/sql"
	"myorm/dialect"
	"myorm/session"
	"os"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var(
	A1 interface{}
	A2 reflect.Value
	P2 []reflect.Value
)

func TestMain(m *testing.M) {
	account := &Account{2020, "asdasdasd"}
	s := []interface{}{account}
	A1 = s[0]
	A2 = reflect.ValueOf(s)
	P2 = []reflect.Value{}
	code := m.Run()
	os.Exit(code)
}

func callByAssert() {
	if a, ok := A1.(ForBenchmarkTest); ok {
		_ = a.ForBenchmarkTest()
	}
}

func callByReflect() {
	if f := A2.MethodByName("ForBenchmarkTest"); f.IsValid(){
		_ = f.Call(P2)
	}
}

func BenchmarkAccount_CallByAssert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		callByAssert()
	}
}

func BenchmarkAccount_CllByReflect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		callByReflect()
	}
}

// 测试结果, 反射性能更好一点
//BenchmarkAccount_CallByAssert
//BenchmarkAccount_CallByAssert-16    	 2065189	       533 ns/op
//BenchmarkAccount_CllByReflect
//BenchmarkAccount_CllByReflect-16    	100000000	        12.2 ns/op
//PASS



func TestAfterQuery(t *testing.T) {
	s := newSession().SetRefTableByModel(Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{ID: 1, Password: "123456"},
		&Account{ID: 2, Password: "qwerty"})
	u := &Account{}
	err := s.FindFirst(u)
	if err != nil || u.Password != "******" {
		t.Fatal("Fail to call hooks after query, got", u)
	}
}

func newSession() *session.Session{
	TestDB, _ := sql.Open("sqlite3", "../gee.db")
	TestDial, _ := dialect.GetDialect("sqlite3")
	return session.New(TestDB, TestDial)
}