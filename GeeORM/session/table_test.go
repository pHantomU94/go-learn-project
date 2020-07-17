package session

import "testing"

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestModel(t *testing.T) {
	sess :=NewSession().Model(&User{})
	_ = sess.DropTable()
	_ = sess.CreateTable()
	if !sess.HasTable() {
		t.Fatalf("Create table failed")
	}
}