package models

import (
	"encoding/json"
	"errors"
	"leedoc/conf"
	"leedoc/utils/gopool"
	"log"
	"strings"
	"time"
)

var (
	exportLimitWorkerChannel = gopool.NewChannelPool(conf.GetExportLimitNum(), conf.GetExportQueueLimitNum())
)

type BookResult struct {
	BookId         int       `json:"book_id"`
	BookName       string    `json:"book_name"`
	ItemId         int       `json:"item_id"`
	ItemName       string    `json:"item_name"`
	Identify       string    `json:"identify"`
	OrderIndex     int       `json:"order_index"`
	Description    string    `json:"description"`
	Publisher      string    `json:"publisher"`
	PrivatelyOwned int       `json:"privately_owned"`
	PrivateToken   string    `json:"private_token"`
	BookPassword   string    `json:"book_password"`
	DocCount       int       `json:"doc_count"`
	CommentStatus  string    `json:"comment_status"`
	CommentCount   int       `json:"comment_count"`
	CreateTime     time.Time `json:"create_time"`
	CreateName     string    `json:"create_name"`
	RealName       string    `json:"real_name"`
	ModifyTime     time.Time `json:"modify_time"`
	Cover          string    `json:"cover"`
	Theme          string    `json:"theme"`
	Label          string    `json:"label"`
	MemberId       int       `json:"member_id"`
	Editor         string    `json:"editor"`
	AutoRelease    bool      `json:"auto_release"`
	HistoryCount   int       `json:"history_count"`

	//RelationshipId     int           `json:"relationship_id"`
	//TeamRelationshipId int           `json:"team_relationship_id"`
	RoleId             conf.BookRole `json:"role_id"`
	RoleName           string        `json:"role_name"`
	Status             int           `json:"status"`
	IsEnableShare      bool          `json:"is_enable_share"`
	IsUseFirstDocument bool          `json:"is_use_first_document"`

	LastModifyText   string `json:"last_modify_text"`
	IsDisplayComment bool   `json:"is_display_comment"`
	IsDownload       bool   `json:"is_download"`
	AutoSave         bool   `json:"auto_save"`
	Lang             string
}

func NewBookResult() *BookResult {
	return &BookResult{}
}

func (m *BookResult) String() string {
	ret, err := json.Marshal(*m)

	if err != nil {
		return ""
	}
	return string(ret)
}

//设置语言
func (m *BookResult) SetLang(lang string) *BookResult {
	m.Lang = lang
	return m
}

//生成默认的项目

// 根据项目标识查询项目以及指定用户权限的信息.
func (m *BookResult) FindByIdentify(identify string, memberId int) (*BookResult, error) {
	if identify == "" || memberId <= 0 {
		return nil, errors.New("参数错误")
	}
	var book Book
	if err := db.Where("identify = ?", identify).First(&book).Error; err != nil {
		log.Println("获取项目失败:", err.Error())
		return nil, err
	}
	roleId, err := NewBook().FindForRoleId(int(book.ID), memberId)
	if err != nil {
		return nil, err
	}
	var relationship2 Relationship
	//查找项目创始人
	if err := db.Where("book_id = ? AND role_id = ?", book.ID, roleId).First(&relationship2).Error; err != nil {
		log.Println("获取项目失败:", err.Error())
		return nil, err
	}
	member, err := NewMember().Find(relationship2.MemberId)
	if err != nil {
		return nil, err
	}
	m.ToBookResult(book)
	m.RoleId = conf.BookRole(roleId)
	m.MemberId = memberId
	m.CreateName = member.Account

	if member.RealName != "" {
		m.RealName = member.RealName
	}
	return m, nil
}

//实体转换
func (m *BookResult) ToBookResult(book Book) *BookResult {

	m.BookId = int(book.ID)
	m.BookName = book.BookName
	m.Identify = book.Identify
	m.OrderIndex = book.OrderIndex
	m.Description = strings.Replace(book.Description, "\r\n", "<br/>", -1)
	m.PrivatelyOwned = book.PrivatelyOwned
	m.PrivateToken = book.PrivateToken
	m.BookPassword = book.BookPassword
	m.DocCount = book.DocCount
	m.CommentStatus = book.CommentStatus
	m.CommentCount = book.CommentCount
	m.CreateTime = book.CreatedAt
	m.ModifyTime = book.UpdatedAt
	m.Cover = book.Cover
	m.Label = book.Label
	m.Status = book.Status
	m.Editor = book.Editor
	m.Theme = book.Theme
	m.AutoRelease = book.AutoRelease == 1
	m.IsEnableShare = book.IsEnableShare == 0
	m.IsUseFirstDocument = book.IsUseFirstDocument == 1
	m.Publisher = book.Publisher
	m.HistoryCount = book.HistoryCount
	m.IsDownload = book.IsDownload == 0
	m.AutoSave = book.AutoSave == 1
	m.ItemId = book.ItemId

	if book.Theme == "" {
		m.Theme = "default"
	}
	if book.Editor == "" {
		m.Editor = "markdown"
	}

	doc := NewDocument()

	err := db.Where("book_id = ?", book.ID).Order("modify_time desc").First(&doc).Error

	if err == nil {
		member2 := NewMember()
		member2.Find(doc.ModifyAt)

		m.LastModifyText = member2.Account + " 于 " + doc.ModifyTime.Local().Format("2006-01-02 15:04:05")
	}

	if m.ItemId > 0 {
		// if item, err := NewItemsets().First(m.ItemId); err == nil {
		// 	m.ItemName = item.ItemName
		// }
	}
	return m
}

//展示用户ID的项目列表
func (m *BookResult) FindBooksByMemberId(memberId int) ([]*BookResult, error) {
	var books []Book
	var bookResults []*BookResult
	if err := db.Where("member_id = ?", memberId).Find(&books).Error; err != nil {
		return nil, err
	}
	for _, book := range books {
		bookResults = append(bookResults, m.ToBookResult(book))
	}
	return bookResults, nil
}

//根据account查询项目
func (m *BookResult) FindBooksByAccount(account string) ([]*BookResult, error) {
	var books []Book
	var bookResults []*BookResult
	if err := db.Where("account = ?", account).Find(&books).Error; err != nil {
		return nil, err
	}
	for _, book := range books {
		bookResults = append(bookResults, m.ToBookResult(book))
	}
	return bookResults, nil
}

//分页查找项目
func (m *BookResult) FindToPager(page int, pageSize int) (books []*BookResult, totalCount int, err error) {
	var count int64
	err = db.Model(&Book{}).Count(&count).Error
	if err != nil {
		log.Println("查询项目总数失败:", err.Error())
		return
	}
	totalCount = int(count)
	sql := `Select 
	book.*,rel.relationship_id,rel.role_id,m.account as create_name,m.real_name as real_name
	from books AS book
	left join relationships AS rel on rel.book_id = book.book_id and rel.role_id = 0
	left join members AS m on m.member_id = rel.member_id
	order by book.order_index desc,book.book_id desc limit ?,?`
	offset := (page - 1) * pageSize
	err = db.Raw(sql, offset, pageSize).Scan(&books).Error
	if err != nil {
		return
	}
	return
}
