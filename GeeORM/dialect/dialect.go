package dialect

import (
	"geeorm/log"
	"reflect"
)

var dialectsMap = map[string]Dialect{}

// Dialect 方言接口，包括类型映射方法以及判断某个表tableName是否存在的SQL语句
// 用来为不同的数据库的对象映射提供统一的接口
type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{}) 
}

// RegisterDialect 注册方言方法，将方言及名称存放在hash表中
func RegisterDialect(name string, d Dialect) {
	dialectsMap[name] = d
}

// 获取方言方法，根据方言名称获取具体方言
func GetDialect(name string) (d Dialect, ok bool) {
	if d, ok = dialectsMap[name]; !ok {
		log.Errorf("The %s dialect is not exists", name)
	}
	return
}


