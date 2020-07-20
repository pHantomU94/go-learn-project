package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

// init clause包初始化函数，注册生成器表
func init() {
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderby
	generators[UPDATE] = _insert
	generators[DELETE] = _insert
	generators[COUNT] = _insert
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
func _insert(values ...interface{}) (string, []interface{}) {
	name := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", name, fields), []interface{}{}
}

// _values 构造INSERT语句中的VALUES字段
func _values(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	var sqlvars []interface{}
	var binStr string
	sql.WriteString("VALUES ")
	for i, value := range values {
		// TODO: interface to string
		varlist := value.([]interface{})
		if binStr == "" {
			binstr = genBinVars(len(varlist))
		}
		sql.WriteString(fmt.Sprintf("(%v)", binStr))
		if i + 1 < len(varlist) {
			sql.WriteString(",")
		}
		sqlvars = append(sqlvars, varlist...)
	}
	return sql.String(), sqlvars
}

func _select(values ...interface{}) (string, []interface{}) {
	name := values[0]
	vars := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", vars, name), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}
func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}
func _orderby(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}
func _update(values ...interface{}) (string, []interface{}) {

}
func _delete(values ...interface{}) (string, []interface{}) {

}
func _count(values ...interface{}) (string, []interface{}) {

}
