package models

import (
	"time"

	"leedoc/conf"
)

var loggerQueue = &logQueue{channel: make(chan *Logger, 100), isRuning: 0}

type logQueue struct {
	channel  chan *Logger
	isRuning int32
}

// Logger struct .
type Logger struct {
	LoggerId int64 `orm:"pk;auto;unique;column(log_id)" json:"log_id"`
	MemberId int   `orm:"column(member_id);type(int)" json:"member_id"`
	// 日志类别：operate 操作日志/ system 系统日志/ exception 异常日志 / document 文档操作日志
	Category     string    `orm:"column(category);size(255);default(operate)" json:"category"`
	Content      string    `orm:"column(content);type(text)" json:"content"`
	OriginalData string    `orm:"column(original_data);type(text)" json:"original_data"`
	PresentData  string    `orm:"column(present_data);type(text)" json:"present_data"`
	CreateTime   time.Time `orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	UserAgent    string    `orm:"column(user_agent);size(500)" json:"user_agent"`
	IPAddress    string    `orm:"column(ip_address);size(255)" json:"ip_address"`
}

// TableName 获取对应数据库表名.
func (m *Logger) TableName() string {
	return "logs"
}

// TableEngine 获取数据使用的引擎.
func (m *Logger) TableEngine() string {
	return "INNODB"
}
func (m *Logger) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewLogger() *Logger {
	return &Logger{}
}
