package session

import (
	"fmt"
	"myorm/log"
	"myorm/schema"
	"reflect"
	"strings"
)

func (s *Session) SetRefTableByModel(value interface{}) *Session {
	// 无关联Table 或 新值类型与关联Table的Model的类型不一致, 则更新RefTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) GetRefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.GetRefTable()
	var columns []string
	for _, field := range table.Fields{
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s)", table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.GetRefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	existSQL, values := s.dialect.TableExistSQL(s.GetRefTable().Name)
	row := s.Raw(existSQL, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.GetRefTable().Name
}
