package models

import (
	"leedoc/conf"
)

// Option struct .
type Option struct {
	OptionId    int    `orm:"column(option_id);pk;auto;unique;" json:"option_id"`
	OptionTitle string `orm:"column(option_title);size(500)" json:"option_title"`
	OptionName  string `orm:"column(option_name);unique;size(80)" json:"option_name"`
	OptionValue string `orm:"column(option_value);type(text);null" json:"option_value"`
	Remark      string `orm:"column(remark);type(text);null" json:"remark"`
}

// TableName 获取对应数据库表名.
func (m *Option) TableName() string {
	return "options"
}

// TableEngine 获取数据使用的引擎.
func (m *Option) TableEngine() string {
	return "INNODB"
}

func (m *Option) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewOption() *Option {
	return &Option{}
}
