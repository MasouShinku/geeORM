/*
	原生交互部分
*/

package session

import (
	"database/sql"
	"geeorm/log"
	"strings"
)

// Session 会话结构体
type Session struct {
	db      *sql.DB         //数据库指针
	sql     strings.Builder // 用于拼接sql语句
	sqlVars []interface{}   // 用于存储拼接sql语句时的占位符
}

// New 创建会话对象
func New(db *sql.DB) *Session {
	return &Session{
		db: db,
	}
}

// Clear 清空sql对象
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

// GetDB 获取DB对象
func (s *Session) GetDB() *sql.DB {
	return s.db
}

// Raw 创建sql语句
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec sql的执行
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Infoln(s.sql.String(), s.sqlVars)
	if result, err = s.GetDB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Errorln(err)
	}
	return
}

// QueryRow 获取单条结果
func (s *Session) QueryRow() (row *sql.Row) {
	defer s.Clear()
	log.Infoln(s.sql.String(), s.sqlVars)
	row = s.GetDB().QueryRow(s.sql.String(), s.sqlVars...)
	return
}

// QueryRows 获取多条结果
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Infoln(s.sql.String(), s.sqlVars)
	if rows, err = s.GetDB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Errorln(err)
	}
	return
}
