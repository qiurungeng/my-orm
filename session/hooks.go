package session

import (
	"myorm/log"
	"reflect"
)

// Hooks constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// Hook 机制:
// 查看参数 value 或 Model 是否实现了 method 方法，实现了则通过反射调用
func (s *Session) CallMethod(method string, value interface{}) {
	var f reflect.Value
	if value == nil{
		f = reflect.ValueOf(s.GetRefTable().Model).MethodByName(method)
	} else {
		f = reflect.ValueOf(value).MethodByName(method)
	}
	// 将Session自身作为备选参数传入方法中
	param := []reflect.Value{reflect.ValueOf(s)}
	if f.IsValid() {
		if v := f.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}