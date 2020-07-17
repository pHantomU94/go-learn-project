package schema

import (
	"geeorm/dialect"
	"testing"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestParse(t *testing.T) {
	user := User{Name: "Jack", Age: 10}
	dial, _ := dialect.GetDialect("sqlite3")
	s := Parse(user, dial) 
	if len(s.FieldNames) != 2 || s.Name != "User" {
		t.Fatalf("failed to parse User struct")
	}
	if s.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatalf("Parse primary key failed")
	}
}