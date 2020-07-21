package session

import (
	"fmt"
	"geeorm/log"
	"geeorm/schema"
	"reflect"
	"strings"
)

// Model 为会话创建或更新维护的表信息
func (sess *Session) Model(value interface{}) *Session{
	// 当会话记录的表为nil时创建或者表类型发生变化时更新
	if sess.refTable == nil || reflect.TypeOf(sess.refTable.Model) != reflect.TypeOf(value) {
		sess.refTable = schema.Parse(value, sess.dial)
	}
	return sess
}

// GetrefTable 返回当前会话维持的表信息
func (sess *Session) GetrefTable() *schema.Schema {
	if sess.refTable == nil {
		log.Errorf("Model is not set")
	}
	return sess.refTable
}

// CreateTable 在数据库中创建一个新的表
func (sess *Session) CreateTable() error {
	table := sess.GetrefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := sess.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// DropTable 根据表名从数据库中删除一张表
func (sess *Session) DropTable() error {
	_, err := sess.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", sess.GetrefTable().Name)).Exec()
	return err
}

// HasTable 检查数据库中是否存在当前会话中维持的数据表
func (sess *Session) HasTable() bool {
	sql, values := sess.dial.TableExistSQL(sess.GetrefTable().Name)
	row := sess.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == sess.GetrefTable().Name
}