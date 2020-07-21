package session

import (
	"geeorm/log"
	"reflect"
)

// 已经注册的钩子函数名
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

// CallMethod 调用钩子函数的入口
func (s *Session) CallMethod(method string, value interface{}) {
	function := reflect.ValueOf(s.GetrefTable().Model).MethodByName(method)
	if value != nil {
		function = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if function.IsValid() {
		if v := function.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return 
}