package main

import (
	"fmt"
	"geeorm"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	engine, _  := geeorm.NewEngine("sqlite3", "gee.db")
	defer engine.Close()
	sess := engine.NewSession()
	_, _ = sess.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = sess.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = sess.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := sess.Raw("INSERT INTO User('Name') values (?), (?)", "Tom", "Jack").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}