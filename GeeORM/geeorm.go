package geeorm

import (
	"database/sql"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
)

// Engine 数据库访问引擎
type Engine struct {
	db *sql.DB
	dial dialect.Dialect
}

// NewEngine 创建新的数据库访问连接，并Ping数据库
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s is not exists", driver)
		return
	}
	e = &Engine{db: db, dial: dial}
	log.Info("Connect database success")
	return &Engine{db: db}, err
}

// Close 关闭数据库访问连接
func (e *Engine) Close() {
	err := e.db.Close()
	if err != nil {
		log.Errorf("Failed to close database")
	}	
	log.Info("Close database success")
}

// NewSession 创建新的数据库访问会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dial)
}
