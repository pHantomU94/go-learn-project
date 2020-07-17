package dialect

import (
	"fmt"
	"reflect"
	"time"
)

// 创建一个匿名变量
var _ Dialect = (*sqlite3)(nil)

// 在dialect包的初始化函数中注册符合sqlite3规则的方言接口
func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

// sqlite3 的方言接口载体，为其实现符合sqlite3规则的类型映射方法及其他相关功能
type sqlite3 struct{}

// DataTypeOf 为sqlite3实现类型映射方法
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
} 

// TableExistSQL 为sqlite3实现判断某个表tableName是否存在的SQL语句
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}



