package session

import (
	"fmt"
	"myorm/log"
	"reflect"
	"testing"
)

// test grammar
func TestReflect(t *testing.T) {
	strs := make([]string, 0)
	addr := &strs
	elem := reflect.Indirect(reflect.ValueOf(addr)).Type().Elem()
	fmt.Println(elem)
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().SetRefTableByModel(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		log.Error(len(users))
		t.Fatal("failed to query all")
	}
}