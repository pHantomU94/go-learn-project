package clause

import (
	"strings"
)

// Type 定义Type为int类型，用来进行枚举
type Type int

// 枚举各种数据库操作操作
const(
	INSERT Type=iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// Clause 数据库操作语句，可以包含多种子操作
type Clause struct {
	sql map[Type]string
	sqlVars map[Type][]interface{}
}

// Set 用来为Clause构造一个给定操作的子语句
func (c *Clause) Set(name Type, vars ...interface{}) {
	
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	c.sql[name], c.sqlVars[name] = generators[name](vars...)
}

// Build 用来根据给定的操作顺序构造完整的SQL语句
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}