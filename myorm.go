package myorm

import (
	"database/sql"
	"fmt"
	"myorm/dialect"
	"myorm/log"
	"myorm/session"
	"strings"
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


/********************************事务支持*********************************/

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


/********************************数据迁移*********************************/

// difference returns set(a) - set(b)
func difference(a, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, val := range b{
		mapB[val] = true
	}
	for _, val := range a{
		if !mapB[val]{
			diff = append(diff, val)
		}
	}
	return
}

// Migrate:
// Model 新增或删除字段, 那么它所映射的表结构也要进行变更
// Migrate 方法帮我们实现对表结构的变更
func (e *Engine) Migrate(newTableStruct interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.SetRefTableByModel(newTableStruct).HasTable() {
			log.Error("table %s doesn't exists.", s.GetRefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.GetRefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1;", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("add columns:[%V], delete columns:[%s]", addCols, delCols)

		//add new fields
		for _, col := range addCols {
			field := table.GetField(col)
			s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, field.Name, field.Type))
			if _, err = s.Exec(); err != nil {
				return
			}
		}

		//delete old & create new
		if len(delCols) > 0 {
			tmpTable := "tmp_" + table.Name
			fields := strings.Join(table.FieldNames, ", ")
			s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmpTable, fields, table.Name))
			s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
			s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmpTable, table.Name))
			_, err = s.Exec()
		}
		return
	})
	return err
}