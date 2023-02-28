package models

import (
	"time"

	"leedoc/conf"
)

// Document struct.
type Document struct {
	DocumentId   int    `gorm:"pk;auto;unique;column(document_id)" json:"doc_id"`
	DocumentName string `gorm:"column(document_name);size(500)" json:"doc_name"`
	// Identify 文档唯一标识
	Identify  string `gorm:"column(identify);size(100);index;null;default(null)" json:"identify"`
	BookId    int    `gorm:"column(book_id);type(int);index" json:"book_id"`
	ParentId  int    `gorm:"column(parent_id);type(int);index;default(0)" json:"parent_id"`
	OrderSort int    `gorm:"column(order_sort);default(0);type(int);index" json:"order_sort"`
	// Markdown markdown格式文档.
	Markdown string `gorm:"column(markdown);type(text);null" json:"markdown"`
	// Release 发布后的Html格式内容.
	Release string `gorm:"column(release);type(text);null" json:"release"`
	// Content 未发布的 Html 格式内容.
	Content    string    `gorm:"column(content);type(text);null" json:"content"`
	CreateTime time.Time `gorm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	MemberId   int       `gorm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime time.Time `gorm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt   int       `gorm:"column(modify_at);type(int)" json:"-"`
	Version    int64     `gorm:"column(version);type(bigint);" json:"version"`
	//是否展开子目录：0 否/1 是 /2 空间节点，单击时展开下一级
	IsOpen     int           `gorm:"column(is_open);type(int);default(0)" json:"is_open"`
	ViewCount  int           `gorm:"column(view_count);type(int)" json:"view_count"`
	AttachList []*Attachment `gorm:"-" json:"attach"`
	//i18n
	Lang string `gorm:"-"`
}

// 多字段唯一键
func (item *Document) TableUnique() [][]string {
	return [][]string{
		[]string{"book_id", "identify"},
	}
}

// TableName 获取对应数据库表名.
func (item *Document) TableName() string {
	return "documents"
}

// TableEngine 获取数据使用的引擎.
func (item *Document) TableEngine() string {
	return "INNODB"
}

func (item *Document) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + item.TableName()
}

func NewDocument() *Document {
	return &Document{
		Version: time.Now().Unix(),
	}
}

//根据文档ID查询指定文档.
