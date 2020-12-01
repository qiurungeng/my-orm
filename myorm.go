package myorm

import (
	"database/sql"
	"myorm/dialect"
	"myorm/log"
	"myorm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil{
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s not found!", driver)
		return
	}
	e = &Engine{
		db: db,
		dialect: dial,
	}
	log.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Fail to close database")
	}else {
		log.Info("Close database success")
	}
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}