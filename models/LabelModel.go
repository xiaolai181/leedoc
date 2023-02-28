package models

import (
	"leedoc/conf"
)

type Label struct {
	LabelId    int    `orm:"column(label_id);pk;auto;unique;" json:"label_id"`
	LabelName  string `orm:"column(label_name);size(50);unique" json:"label_name"`
	BookNumber int    `orm:"column(book_number)" json:"book_number"`
}

// TableName 获取对应数据库表名.
func (m *Label) TableName() string {
	return "label"
}

// TableEngine 获取数据使用的引擎.
func (m *Label) TableEngine() string {
	return "INNODB"
}

func (m *Label) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewLabel() *Label {
	return &Label{}
}
