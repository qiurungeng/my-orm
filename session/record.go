package session

import (
	"errors"
	"myorm/clause"
	"myorm/log"
	"reflect"
)

// 多次构造 clause 子句， 最后在统一 build 一次

// Insert: 传入struct, 将它们插入所关联的表中, 返回插入记录数
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for i, value := range values{
		refTable := s.SetRefTableByModel(value).GetRefTable()
		if i == 0 {
			s.clause.Set(clause.INSERT, refTable.Name,  refTable.FieldNames)
		}
		recordValues = append(recordValues, refTable.ToRecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		log.Error(err)
	}
	return res.RowsAffected()
}

// Find: 传入一个空切片, 将查询结果都塞进去
func (s *Session) Find(valueSlice interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(valueSlice))
	destType := destSlice.Type().Elem()
	table := s.SetRefTableByModel(reflect.New(destType).Elem().Interface()).GetRefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		log.Error(err)
		return err
	}
	// 将查询结果转成结构体塞入被传入切片当中
	for rows.Next(){
		dest := reflect.New(destType).Elem()
		var values []interface{}
		// 获取每个字段的指针值
		for _, name := range table.FieldNames{
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// 将该条记录填充到各字段指针当中
		if err := rows.Scan(values...); err != nil{
			return err
		}
		//valueSlice = append(valueSlice, dest)
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}


// FindFirst: 查询第一条记录
func (s *Session) FindFirst(resultPtr interface{}) error {
	result := reflect.Indirect(reflect.ValueOf(resultPtr))
	destSlice := reflect.New(reflect.SliceOf(result.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil{
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	result.Set(destSlice.Index(0))
	return nil
}

// Update:
// Session s 必须提前设定 refTable
// 传入参数支持 map, 也支持 field1, value1, field2, value2......
func (s *Session) Update(kv ...interface{}) (int64, error) {
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0 ; i < len(kv); i += 2 {
			m[kv[i].(string) + " = ?"] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.GetRefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete:
// 必须提前设定 refTable
func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.GetRefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Count: SELECT COUNT(*) ... WHERE ...
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.GetRefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var count int64
	if err := row.Scan(&count); err != nil{
		return 0, err
	}
	return count, nil
}


//链式调用

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(condition string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, condition), args...)...)
	return s
}

func (s *Session) OrderBy(field string) *Session {
	s.clause.Set(clause.ORDERBY, field)
	return s
}