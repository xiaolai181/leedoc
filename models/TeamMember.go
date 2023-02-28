package models

import (
	"leedoc/conf"
)

type TeamMember struct {
	TeamMemberId int `orm:"column(team_member_id);pk;auto;unique;" json:"team_member_id"`
	TeamId       int `orm:"column(team_id);type(int)" json:"team_id"`
	MemberId     int `orm:"column(member_id);type(int)" json:"member_id"`
	// RoleId 角色：0 创始人(创始人不能被移除) / 1 管理员/2 编辑者/3 观察者
	RoleId   conf.BookRole `orm:"column(role_id);type(int)" json:"role_id"`
	RoleName string        `orm:"-" json:"role_name"`
	Account  string        `orm:"-" json:"account"`
	RealName string        `orm:"-" json:"real_name"`
	Avatar   string        `orm:"-" json:"avatar"`
	Lang     string        `orm:"-"`
}
