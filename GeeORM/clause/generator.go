package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

// init clause包初始化函数，注册生成器表
func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderby
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// genBinVars 用来为插入的数据创建占位符字符串
func genBinVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}

// _insert 构建INSERT语句
// "INSERT INTO %s (%v)"
func _insert(values ...interface{}) (string, []interface{}) {
	name := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", name, fields), []interface{}{}
}

// _values 构造INSERT语句中的VALUES字段
// "VALUES (?, ?, ...,?), (...) v, v, ..., v"
func _values(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	var sqlvars []interface{}
	var binStr string
	sql.WriteString("VALUES ")
	for i, value := range values {
		// TODO: interface to string
		varlist := value.([]interface{})
		if binStr == "" {
			binStr = genBinVars(len(varlist))
		}
		sql.WriteString(fmt.Sprintf("(%v)", binStr))
		if i + 1 < len(values) {
			sql.WriteString(",")
		}
		sqlvars = append(sqlvars, varlist...)
	}
	return sql.String(), sqlvars
}

// _select 构造SELECT语句
// "SLEECT %v FROM %s"
func _select(values ...interface{}) (string, []interface{}) {
	name := values[0]
	vars := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", vars, name), []interface{}{}
}

// _limit 构造LIMIT语句
// “LIMIT ?”
func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

// _where 构造WHERE语句
// "WHERE %s"
func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

// _orderby 构造 ORDER BY 语句
//  "ORDER BY %s"
func _orderby(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// _update 构造UPDATE语句 
// "UPDATE %s SET %s"
func _update(values ...interface{}) (string, []interface{}) {
	name := values[0]
	var vars []interface{}
	var keys []string
	m := values[1].(map[string]interface{})
	for k, v := range m {
		vars = append(vars, v)
		keys = append(keys, k + " = ?")
	}
	return fmt.Sprintf("UPDATE %s SET %s", name, strings.Join(keys, ",")), vars
}

// _delete 构造DELETE语句
// "DELETE FROM %s"
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

// _count 构造COUNT语句
// "SELECT count(*) FROM %s"
func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
