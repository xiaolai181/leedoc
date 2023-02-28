package models

import (
	"time"

	"leedoc/conf"
)

type TeamRelationship struct {
	TeamRelationshipId int       `orm:"column(team_relationship_id);pk;auto;unique;" json:"team_relationship_id"`
	BookId             int       `orm:"column(book_id)" json:"book_id"`
	TeamId             int       `orm:"column(team_id)" json:"team_id"`
	CreateTime         time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	TeamName           string    `orm:"-" json:"team_name"`
	MemberCount        int       `orm:"-" json:"member_count"`
	BookMemberId       int       `orm:"-" json:"book_member_id"`
	BookMemberName     string    `orm:"-" json:"book_member_name"`
	BookName           string    `orm:"-" json:"book_name"`
}

// TableName 获取对应数据库表名.
func (m *TeamRelationship) TableName() string {
	return "team_relationship"
}
func (m *TeamRelationship) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

// TableEngine 获取数据使用的引擎.
func (m *TeamRelationship) TableEngine() string {
	return "INNODB"
}

// 联合唯一键
func (m *TeamRelationship) TableUnique() [][]string {
	return [][]string{{"team_id", "book_id"}}
}
