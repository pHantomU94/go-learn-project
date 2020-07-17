package schema

import (
	"geeorm/dialect"
	"geeorm/log"
	"go/ast"
	"reflect"
)

// Field 表字段类型，用来映射一个成员变量与数据库中的一个字段
type Field struct {
	Name string
	Type string
	Tag string
}

// Schema 表概要类型，用来维护一个对象与一张数据库中的表之间的映射关系，存储表中相关数据
type Schema struct {
	Model interface{}
	Name string
	Fields []*Field
	FieldNames []string
	fieldMap map[string]*Field
}

// GetField 根据字段名称获取对应字段
func (s *Schema) GetField(name string) *Field{
	field, ok := s.fieldMap[name]
	if !ok {
		log.Errorf("Field %s is not exists in table %s", name, s.Name)
	}
	return field
}

// Parse 用来将一个对象映射成一个表概要
func Parse(obj interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(obj)).Type()
	s := &Schema{
		Model: obj,
		Name: modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		sf := modelType.Field(i)
		if !sf.Anonymous && ast.IsExported(sf.Name){
			field := &Field{
				Name: sf.Name,
				Type: sf.Type.Name(),
			}
			if tag, ok := sf.Tag.Lookup("geeorm"); ok {
				field.Tag = tag
			}
			s.Fields = append(s.Fields, field)
			s.FieldNames = append(s.FieldNames, field.Name)
			s.fieldMap[field.Name] = field
		}
	}
	// 遍历对象的每一个成员，将其映射成表中的字段
	return s
}
