package geeorm

import (
	"testing"
	_ "github.com/mattn/go-sqlite3"
)
	

func OpenDB(t *testing.T) *Engine{
	t.Helper()
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("failed to connect to database", err)
	}
	return engine
}

func TestNewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	sess := engine.NewSession()
	sess.Clear()
}