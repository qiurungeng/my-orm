package session

import (
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