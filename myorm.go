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


// 编程式事务支持:

type TxFunc func(*session.Session) (result interface{}, err error)

func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	//事务新开启一个Session
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		//出现panic时先Rollback再panic
		if p := recover(); p != nil {
			//忽略此处错误, 不要让Rollback的error覆盖了f的error
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	return f(s)
}