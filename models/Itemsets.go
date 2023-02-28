package models

import (
	"time"
)

//项目空间
type Itemsets struct {
	ItemId      int       `orm:"column(item_id);pk;auto;unique" json:"item_id"`
	ItemName    string    `orm:"column(item_name);size(500)" json:"item_name"`
	ItemKey     string    `orm:"column(item_key);size(100);unique" json:"item_key"`
	Description string    `orm:"column(description);type(text);null" json:"description"`
	MemberId    int       `orm:"column(member_id);size(100)" json:"member_id"`
	CreateTime  time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	ModifyTime  time.Time `orm:"column(modify_time);type(datetime);null;auto_now" json:"modify_time"`
	ModifyAt    int       `orm:"column(modify_at);type(int)" json:"modify_at"`

	BookNumber       int    `orm:"-" json:"book_number"`
	CreateTimeString string `orm:"-" json:"create_time_string"`
	CreateName       string `orm:"-" json:"create_name"`
}
