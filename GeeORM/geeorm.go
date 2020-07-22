package geeorm

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
	"strings"
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
	return 
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

// TxFunc 事务函数模板
type TxFunc func(*session.Session) (reslut interface{}, err error)

// Transaction 事务接口
func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			_ = s.Commit()
		} 
	}()
	return f(s)
}

// difference 用来比较两个字符串数组的差异 返回 a-b
func difference (a []string, b []string) []string {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	res := make([]string, 0)
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			res = append(res, v)
		}
	}
	return res
} 

// Migrate 数据库迁移操作
// 根据输入的对象类型名查找对应的表，比对条目的差异，先增加新增的列，创建一张新表迁移旧表，再删除旧表，修改新表名
func (e *Engine) Migrate(value interface{}) error {
	_ ,err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s is not exists", s.GetrefTable().Name)
			return nil, s.CreateTable()
		}		
		table := s.GetrefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)
		for _, col := range addCols {
			f := table.GetField(col)
			if _, err = s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name,f.Name, f.Type)).Exec(); err != nil {
				return
			}
		}
		
		if len(delCols) == 0 {
			return
		}
		tmp := "Tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ",")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENMAE TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}