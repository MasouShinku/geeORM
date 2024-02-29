// 增删改查相关代码
package session

import (
	"geeorm/clause"
	"reflect"
)

// Insert 期望操作形式：
// u1:=&User{Name:"shinku" ,Age:18}
// u2:=&User{Name:"shizuka",Age:20}
// s.Insert(u1,u2)
// 实际sql语句:
// INSERT INTO $tableName($col1,$col2,...) VALUES
// (A1,A2,...),
// (B1,B2,...),
// ...
func (s *Session) Insert(values ...interface{}) (int64, error) {

	// 设置INSERT部分
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	// 设置VALUES部分
	s.clause.Set(clause.VALUES, recordValues...)

	// 拼接生成完整语句
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Find 查询语句
// 期望调用形式：
// var users []User
// s.Find(&users)
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))

	// 获取元素类型，并映射出表结构
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	// 拼接出select语句并执行
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	// 将结果平铺，遍历后按字段赋值
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return rows.Close()

}
