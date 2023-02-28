package models

import (
	"time"

	"leedoc/conf"
)

type Template struct {
	TemplateId   int    `orm:"column(template_id);pk;auto;unique;" json:"template_id"`
	TemplateName string `orm:"column(template_name);size(500);" json:"template_name"`
	MemberId     int    `orm:"column(member_id);index" json:"member_id"`
	BookId       int    `orm:"column(book_id);index" json:"book_id"`
	BookName     string `orm:"-" json:"book_name"`
	//是否是全局模板：0 否/1 是; 全局模板在所有项目中都可以使用；否则只能在创建模板的项目中使用
	IsGlobal        int       `orm:"column(is_global);default(0)" json:"is_global"`
	TemplateContent string    `orm:"column(template_content);type(text);null" json:"template_content"`
	CreateTime      time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	CreateName      string    `orm:"-" json:"create_name"`
	ModifyTime      time.Time `orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	ModifyAt        int       `orm:"column(modify_at);type(int)" json:"-"`
	ModifyName      string    `orm:"-" json:"modify_name"`
	Version         int64     `orm:"type(bigint);column(version)" json:"version"`
}

// TableName 获取对应数据库表名.
func (m *Template) TableName() string {
	return "templates"
}

// TableEngine 获取数据使用的引擎.
func (m *Template) TableEngine() string {
	return "INNODB"
}

func (m *Template) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewTemplate() *Template {
	return &Template{}
}
