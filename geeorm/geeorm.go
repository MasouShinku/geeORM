// Engine模块
// 交互前的准备工作以及交互后的收尾工作
package geeorm

import (
	"database/sql"
	"geeorm/log"
	"geeorm/session"
)

// Engine 结构体定义
type Engine struct {
	db *sql.DB
}

// NewEngine 创建实例方法
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Errorln(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Errorln(err)
		return
	}
	e = &Engine{db: db}
	log.Infoln("数据库连接成功!")
	return
}

// Close 关闭数据库连接
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Errorln("数据库关闭失败!")
	}
	log.Infoln("数据库关闭成功!")
}

// NewSession 创建会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}
