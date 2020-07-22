package session

import "geeorm/log"

// Begin 开始事务
func (s *Session) Begin() (err error) {
	log.Info("transcation begin")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

// Commit 事务提交
func (s *Session) Commit() (err error) {
	log.Info("transcation commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
		return
	}
	return
}

// Rollback 事务提交
func (s *Session) Rollback() (err error) {
	log.Info("transcation rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
		return
	}
	return
}