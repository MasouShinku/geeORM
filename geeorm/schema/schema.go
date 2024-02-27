// 实现go对象与数据库表的转换
package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

// Field 数据库属性条目
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema 数据库表格对象
// 只存放属性信息，实例对象在之后添加
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string // 单独提取列名，防止之后需要遍历对象
	fieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

// 将任意对象解析为schema实例
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		temp := modelType.Field(i)
		if !temp.Anonymous && ast.IsExported(temp.Name) {
			field := &Field{
				Name: temp.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(temp.Type))),
			}
			if v, ok := temp.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, temp.Name)
			schema.fieldMap[temp.Name] = field
		}
	}
	return schema
}
