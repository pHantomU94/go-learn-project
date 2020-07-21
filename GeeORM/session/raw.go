package session

import (
	"database/sql"
	"geeorm/clause"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/schema"
	"strings"
)

// Session 数据库访问会话
type Session struct {
	db *sql.DB	// 数据库指针
	dial dialect.Dialect // 所连接数据库类型的方言
	refTable *schema.Schema // 会话当前维护的数据库表
	sql strings.Builder // 数据库操作语句
	sqlVars []interface{} // 数据库操作占位符对应的参数
	clause clause.Clause
}

// New 用于创建一个新的数据库访问会话
func New(db *sql.DB, dial dialect.Dialect) *Session {
	return &Session{
		db: db,
		dial: dial,
	}
} 

// Clear 清理会话中的sql语句与参数，使请求可以复用
func (sess *Session) Clear() {
	sess.sql.Reset()
	sess.sqlVars = nil
	sess.clause = clause.Clause{}
}

// DB 返回会话的数据库指针
func (sess *Session) DB() *sql.DB {
	return sess.db
}

// Raw 构建数据库访问原始请求
func (sess *Session) Raw(sql string, sqlVars ...interface{}) *Session {
	sess.sql.WriteString(sql)
	sess.sql.WriteString(" ")
	sess.sqlVars = append(sess.sqlVars, sqlVars...)
	return sess
}

// Exec 数据库原始Exec操作
func (sess *Session) Exec() (result sql.Result, err error) {
	defer sess.Clear()
	log.Info(sess.sql.String(), sess.sqlVars)
	if result, err = sess.DB().Exec(sess.sql.String(), sess.sqlVars...); err != nil {
		log.Error(err)
	}
	return result, err
}

// QueryRow 数据库查询一行QueryRaw操作
func (sess *Session) QueryRow() (*sql.Row) {
	defer sess.Clear()
	log.Info(sess.sql.String(), sess.sqlVars)
	return sess.DB().QueryRow(sess.sql.String(), sess.sqlVars...)
}

// QueryRows 数据库查询多行QueryRaws操作
func (sess *Session) QueryRows() (*sql.Rows, error) {
	defer sess.Clear()
	log.Info(sess.sql.String(), sess.sqlVars)
	rows, err := sess.DB().Query(sess.sql.String(), sess.sqlVars...)
	if err != nil {
		log.Error(err)
	}
	return rows, err
}