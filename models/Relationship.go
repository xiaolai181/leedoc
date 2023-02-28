package models

import (
	"errors"
	"log"

	"leedoc/conf"
)

type Relationship struct {
	RelationshipId int `gorm:"pk;auto;unique;column(relationship_id)" json:"relationship_id"`
	MemberId       int `gorm:"column(member_id);type(int);" json:"member_id"`
	BookId         int `gorm:"column(book_id);type(int)" json:"book_id"`
	// RoleId 角色：0 创始人(创始人不能被移除) / 1 管理员/2 编辑者/3 观察者
	RoleId int `gorm:"column(role_id);type(int)" json:"role_id"`
}

// TableName 获取对应数据库表名.
func (m *Relationship) TableName() string {
	return "relationship"
}
func (m *Relationship) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

// TableEngine 获取数据使用的引擎.
func (m *Relationship) TableEngine() string {
	return "INNODB"
}

// 联合唯一键
func (m *Relationship) TableUnique() [][]string {
	return [][]string{
		{"member_id", "book_id"},
	}
}

func NewRelationship() *Relationship {
	return &Relationship{}
}

//
//根据relationship ID获取relationship
func (m *Relationship) Find(id int) (*Relationship, error) {
	db.Where("id=?", id).Find(m)
	return m, nil
}

//查询指定项目的创始人.
func (m *Relationship) FindFounder(book_id int) (*Relationship, error) {
	founder := NewRelationship()
	err := db.Where("book_id=? ", book_id).Find(founder).Error
	if err != nil {
		return nil, err
	}
	return founder, nil
}

//变更创始人权限
func (m *Relationship) UpdateRoleId(bookId, memberId int, roleId int) (*Relationship, error) {
	book := NewBook()
	book.ID = uint(bookId)
	if err := db.First(book, book.ID).Error; err != nil {
		log.Println("Relationship.UpdateRoleId:", err)
		return nil, err
	}
	err := db.Where("book_id = ? AND member_id = ?", bookId, memberId).First(m).Error
	if err == nil {
		m = NewRelationship()
		m.BookId = bookId
		m.MemberId = memberId
		m.RoleId = roleId
	} else if err != nil {
		return nil, err
	} else if m.RoleId == int(conf.BookFounder) {
		return m, errors.New("不能变更创始人的权限")
	}
	m.RoleId = roleId

	if m.RelationshipId > 0 {
		err = db.Save(m).Error
	} else {
		err = db.Create(m).Error
	}
	return m, err

}

//根据项目ID与成员ID获取项目权限
func (m *Relationship) FindForRoleId(bookId, memberId int) (int, error) {
	relationship := NewRelationship()
	err := db.Where("book_id=? AND member_id=?", bookId, memberId).Find(relationship).Error
	if err != nil {
		return -1, err
	}
	return relationship.RoleId, nil
}

//FindByBookIdAndMemberId 根据项目ID和成员ID获取关系
func (m *Relationship) FindByBookIdAndMemberId(book_id, member_id int) (*Relationship, error) {
	relationship := NewRelationship()
	err := db.Where("book_id=? AND member_id=?", book_id, member_id).Find(relationship).Error
	if err != nil {
		return nil, err
	}
	return relationship, nil
}

//插入新的关系
func (m *Relationship) Insert() error {
	err := db.Create(m).Error
	return err
}

//更新关系
func (m *Relationship) Update() error {
	return db.Save(m).Error
}

//删除关系
func (m *Relationship) DeleteByBookIdAndMemberId(book_id, member_id int) error {
	err := db.Where("book_id=? AND member_id=?", book_id, member_id).Delete(m).Error
	if err != nil {
		log.Println("删除关系失败:", err)
		return err
	}
	return nil
}

//转让项目
func (m *Relationship) Transfer(book_id, founder_id, receive_id int) error {
	//删除原创始人
	err := db.Where("book_id=? AND member_id=?", book_id, founder_id).Delete(m).Error
	if err != nil {
		return err
	}
	//插入新创始人
	m.BookId = book_id
	m.MemberId = receive_id
	m.RoleId = int(conf.BookFounder)
	err = db.Create(m).Error
	if err != nil {
		return err
	}
	return nil
}
