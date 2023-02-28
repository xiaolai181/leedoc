package models

import (
	"time"

	"leedoc/conf"
)

type DocumentHistory struct {
	HistoryId    int       `gorm:"column(history_id);pk;auto;unique" json:"history_id"`
	Action       string    `gorm:"column(action);size(255)" json:"action"`
	ActionName   string    `gorm:"column(action_name);size(255)" json:"action_name"`
	DocumentId   int       `gorm:"column(document_id);type(int);index" json:"doc_id"`
	DocumentName string    `gorm:"column(document_name);size(500)" json:"doc_name"`
	ParentId     int       `gorm:"column(parent_id);type(int);index;default(0)" json:"parent_id"`
	Markdown     string    `gorm:"column(markdown);type(text);null" json:"markdown"`
	Content      string    `gorm:"column(content);type(text);null" json:"content"`
	MemberId     int       `gorm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime   time.Time `gorm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt     int       `gorm:"column(modify_at);type(int)" json:"-"`
	Version      int64     `gorm:"type(bigint);column(version)" json:"version"`
	IsOpen       int       `gorm:"column(is_open);type(int);default(0)" json:"is_open"`
}

type DocumentHistorySimpleResult struct {
	HistoryId  int       `json:"history_id"`
	ActionName string    `json:"action_name"`
	MemberId   int       `json:"member_id"`
	Account    string    `json:"account"`
	ModifyAt   int       `json:"modify_at"`
	ModifyName string    `json:"modify_name"`
	ModifyTime time.Time `json:"modify_time"`
	Version    int64     `json:"version"`
}

// TableName 获取对应数据库表名.
func (m *DocumentHistory) TableName() string {
	return "document_history"
}

// TableEngine 获取数据使用的引擎.
func (m *DocumentHistory) TableEngine() string {
	return "INNODB"
}

func (m *DocumentHistory) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewDocumentHistory() *DocumentHistory {
	return &DocumentHistory{}
}
