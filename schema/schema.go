package schema

import (
	"go/ast"
	"myorm/dialect"
	"reflect"
)

// Field 表示 db 中的一列
type Field struct {
	Name string // 字段名
	Type string // 类型
	Tag  string // 约束条件
}

// Schema 表示 db 中的一张表
type Schema struct {
	Model      interface{}       // 被映射对象
	Name       string            // 表名
	Fields     []*Field           // 字段
	FieldNames []string          // 所有字段名(列名)
	fieldMap   map[string]*Field // 字段名和Field之间的映射关系
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

// Parse: 将任意对象解析为 Schema
func Parse(dest interface{}, dialect dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model: dest,
		Name: modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// 若该字段 不是匿名字段 且 包外可见
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: dialect.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			// 额外约束条件: 如, PRIMARY KEY
			if v, ok := p.Tag.Lookup("myorm"); ok{
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

// 将 struct 转为 表中记录值数组
func (s *Schema) ToRecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range s.Fields{
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}