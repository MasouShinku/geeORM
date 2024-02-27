package geeorm

import (
	"fmt"
	"geeorm/session"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}

func TestSession_Model(t *testing.T) {
	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	table := s.RefTable()
	s.Model(&session.Session{})
	t.Log(fmt.Sprintf("table name is %s", table.Name))
	t.Log(fmt.Sprintf("refTable name is %s", s.RefTable().Name))
	if table.Name != "User" || s.RefTable().Name != "Session" {
		t.Fatal("Failed to change model")
	}
}
