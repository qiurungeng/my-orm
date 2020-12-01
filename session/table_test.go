package session

import "testing"

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().SetRefTableByModel(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}

func TestSession_Model(t *testing.T) {
	s := NewSession().SetRefTableByModel(&User{})
	table := s.GetRefTable()
	s.SetRefTableByModel(&Session{})
	if table.Name != "User" || s.GetRefTable().Name != "Session" {
		t.Fatal("Failed to change model")
	}
}