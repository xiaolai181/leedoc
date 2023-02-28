package models

import (
	"errors"
	"fmt"
	"leedoc/conf"
	"leedoc/utils"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/beego/i18n"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model

	Account       string          `gorm:"column:account;unique_index" json:"account" form:"account"`
	RealName      string          `gorm:"column:real_name;type:varchar(255)" json:"real_name" form:"real_name"`
	Password      string          `gorm:"column:password;type:varchar(1000)" json:"password" form:"password"`
	AuthMethod    string          `gorm:"column:auth_method;default:'local';type:varchar(50)" json:"auth_method" form:"auth_method"`
	Description   string          `gorm:"column:description ;type:varchar(2000)" json:"description" form:"description"`
	Email         string          `gorm:"column:email;type:varchar(255)" json:"email" form:"email"`
	Phone         string          `gorm:"column:phone;type:varchar(255)" json:"phone" form:"phone"`
	Avatar        string          `gorm:"column:avatar;type:varchar(255)" json:"avatar" form:"avatar"`
	Role          conf.SystemRole `gorm:"column:role;type:int;" json:"role" form:"role"`
	RoleName      string          `gorm:"-" json:"role_name" form:"role_name"`
	Status        int             `gorm:"column:status;type:int;" json:"status" form:"status"`
	CreateId      int             `gorm:"column:create_id;type:int" json:"create_id" form:"create_id"`
	LastLoginTime time.Time       `gorm:"column:last_login_time;type:datetime" json:"last_login_time" form:"last_login_time"`
	//i18n
	Lang string `gorm:"-" json:"lang" form:"lang"`
}

// TableName 获取对应数据库表名.
func (m *Member) TableName() string {
	return "members"
}

// TableEngine 获取数据使用的引擎.
func (m *Member) TableEngine() string {
	return "INNODB"
}

func (m *Member) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewMember() *Member {
	return &Member{}
}

//默认用户
func RegisterAdmin(account, password string) {
	member := Member{
		Account:  account,
		Password: password,
		Role:     0,
	}
	if member.Email == "" {
		member.Email = account + "@leedoc.com"
	}
	if err := member.Add(); err != nil {
		panic(err)
	}
}

//用户登录
func (m *Member) Login() (bool, error) {
	var member Member
	ps := m.Password
	if err := db.Where("account = ? ", m.Account).First(&member).Error; err != nil {
		return false, err
	}
	fmt.Println(member.Password, ps)
	ok, err := utils.PasswordVerify(member.Password, ps)
	if member.ID > 0 && ok {
		log.Println("Member login success:", member.Account)
		member.LastLoginTime = time.Now()
		if err := member.Update(); err != nil {
			log.Println("Update member login time error:", err)
		}
		return true, nil

	}
	return false, err
}

//新增用户
func (m *Member) Add() error {
	log.Println("Add member:", m.Account)
	if ok, err := regexp.MatchString(conf.RegexpAccount, m.Account); m.Account == "" || !ok || err != nil {
		return errors.New("用户名只能由英文字母数字组成，且在3-50个字符")
	}
	if m.Email == "" {
		log.Println("m.Email is empty")
		return errors.New("邮箱不能为空")
	}
	if ok, err := regexp.MatchString(conf.RegexpEmail, m.Email); !ok || err != nil || m.Email == "" {
		log.Println("邮箱格式不正确")
		return errors.New("邮箱格式不正确")
	}
	if m.AuthMethod == "" {
		m.AuthMethod = "local"
	}
	if m.AuthMethod == "local" {
		if l := strings.Count(m.Password, ""); l < 6 || l >= 50 {
			log.Println("密码长度错误")
			return errors.New("密码不能为空且必须在6-50个字符之间")
		}
	}
	if c, err := m.IsExistEmail(); err == nil && c {
		return errors.New("邮箱已被使用")
	}
	if c, err := m.IsExistAccount(m.Account); err == nil && c {
		return errors.New("账号已被使用")
	}
	hash, err := utils.PasswordHash(m.Password)
	if err != nil {
		log.Println("PasswordHash error:", err)
		return errors.New("密码加密失败")
	}
	m.Password = hash

	m.ResolveRoleName()
	log.Println("Add member:", m.Account, m.RoleName)
	if err := db.Create(m).Error; err != nil {
		log.Println("Create member error:", err)
		return errors.New("创建用户失败")
	}
	return nil
}

//更新用户
func (m *Member) Update(cols ...string) error {
	if m.Email == "" {
		return errors.New("邮箱不能为空")
	}
	if err := db.Where("id = ?", m.ID).Save(m).Error; err != nil {
		log.Println("Update member error:", err)
		return errors.New("更新用户失败")
	}
	return nil
}

//查找用户
func (m *Member) Find(id int) (*Member, error) {
	if err := db.Where("id = ?", id).First(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

//根据account查找用户ID
func (m *Member) FindIDByAccount(account string) (int, error) {
	if err := db.Where("account = ?", account).First(&m).Error; err != nil {
		return 0, err
	}
	return int(m.ID), nil
}

//根据账号查找用户
func (m *Member) FindByAccount(account string) (*Member, error) {
	if err := db.Where("account = ?", account).First(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

//根据账号批量查找用户
func (m *Member) FindByAccounts(accounts []string) ([]*Member, error) {
	var members []*Member
	if err := db.Where("account in (?)", accounts).Find(&members).Error; err != nil {
		return nil, err
	} else {
		for _, item := range members {
			item.ResolveRoleName()
		}
	}
	return members, nil
}

//分页查找用户
func (m *Member) FindToPager(pageIndex int, pageSize int) ([]*Member, int, error) {
	offset := (pageIndex - 1) * pageSize
	var totalCount int64
	err := db.Model(&Member{}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}
	var members []*Member
	err = db.Order("id desc").Offset(offset).Limit(pageSize).Find(&members).Error
	if err != nil {
		return members, 0, err
	}

	for _, m := range members {
		m.ResolveRoleName()
	}
	return members, int(totalCount), nil
}

//判断是否为管理员
func (m *Member) IsAdmin() bool {
	if m == nil || m.ID <= 0 {
		return false
	}
	return m.Role == 0 || m.Role == 1
}

// 根据指定字段查找用户
func (m *Member) FindByFieldFirst(field string, value interface{}) (*Member, error) {
	var member Member
	if err := db.Where(field+" = ?", value).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

//校验用户
func (m *Member) Valid(is_hash_password bool) error {
	//邮箱不能为空
	if m.Email == "" {
		return errors.New("邮箱不能为空")
	}
	//用户描述不能大于255个字符
	if len(m.Description) > 255 {
		return errors.New("用户描述不能大于255个字符")
	}
	//判断用户角色
	if m.Role != conf.MemberGeneralRole && m.Role != conf.MemberSuperRole && m.Role != conf.MemberAdminRole {
		return errors.New("用户角色不正确")
	}
	//判断用户状态
	if m.Status != 0 && m.Status != 1 {
		m.Status = 0
	}
	//邮箱格式校验
	if ok, err := regexp.MatchString(conf.RegexpEmail, m.Email); !ok || err != nil || m.Email == "" {
		return errors.New("邮箱格式不正确")
	}
	//校验密码格式
	if !is_hash_password {
		if l := strings.Count(m.Password, ""); m.Password == "" || l > 50 || l < 6 {
			return errors.New("密码不能为空且必须在6-50个字符之间")
		}
	}
	//校验邮箱是否被使用
	if c, err := m.IsExistEmail(); err != nil || c {
		return errors.New("邮箱已被使用")
	}
	//校验账号是否被使用
	if m.ID > 0 {
		//校验用户是否存在
		result := db.First(&Member{}, m.ID)
		if result.RowsAffected == 0 {
			return errors.New("用户不存在")
		}
		if result.Error != nil {
			return errors.New("查询用户失败")
		}
	} else {
		//校验账号格式是否正确
		if ok, err := regexp.MatchString(conf.RegexpAccount, m.Account); m.Account == "" || !ok || err != nil {
			return errors.New("账号格式不正确")
		}
		//校验账号是否被使用
		if c, err := m.IsExistAccount(m.Account); err == nil && c {
			return errors.New("账号已被使用")
		}
	}
	return nil
}

//根据ID删除用户
func (m *Member) Delete(id int) error {
	if err := db.Where("id = ?", id).Delete(&Member{}).Error; err != nil {
		return err
	}
	return nil
}

//查找用户是否存在
func (m *Member) IsExistAccount(account string) (bool, error) {
	var mb Member
	if err := db.Where("account = ?", account).First(&mb).Error; err != nil {
		return false, err
	}
	if mb.ID > 0 {
		return true, nil
	}
	return false, nil
}

//查找邮箱是否存在
func (m *Member) IsExistEmail() (bool, error) {
	var member Member
	if err := db.Where("email = ?", m.Email).First(&member).Error; err != nil {
		return false, err
	}
	if member.ID > 0 {
		return true, nil
	}
	return false, nil
}

//检查用户角色
func (m *Member) ResolveRoleName() {
	if m.Role == conf.MemberSuperRole {
		m.RoleName = i18n.Tr(m.Lang, "common.administrator")
	} else if m.Role == conf.MemberAdminRole {
		m.RoleName = i18n.Tr(m.Lang, "common.editor")
	} else if m.Role == conf.MemberGeneralRole {
		m.RoleName = i18n.Tr(m.Lang, "common.obverser")
	}
}
