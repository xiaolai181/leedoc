package models

import (
	"time"

	"leedoc/conf"
)

type MemberToken struct {
	TokenId   int       `orm:"column(token_id);pk;auto;unique" json:"token_id"`
	MemberId  int       `orm:"column(member_id);type(int)" json:"member_id"`
	Token     string    `orm:"column(token);size(150);index" json:"token"`
	Email     string    `orm:"column(email);size(255)" json:"email"`
	IsValid   bool      `orm:"column(is_valid)" json:"is_valid"`
	ValidTime time.Time `orm:"column(valid_time);null" json:"valid_time"`
	SendTime  time.Time `orm:"column(send_time);auto_now_add;type(datetime)" json:"send_time"`
}

// TableName 获取对应数据库表名.
func (m *MemberToken) TableName() string {
	return "member_token"
}

// TableEngine 获取数据使用的引擎.
func (m *MemberToken) TableEngine() string {
	return "INNODB"
}

func (m *MemberToken) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewMemberToken() *MemberToken {
	return &MemberToken{}
}
