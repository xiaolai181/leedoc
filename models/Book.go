package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"leedoc/conf"
	"leedoc/utils"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/beego/i18n"
	"gorm.io/gorm"
)

var releaseQueue = make(chan int, 100)

var once = sync.Once{}

type Book struct {
	gorm.Model
	BookName string `gorm:"column(book_name);size(500)" json:"book_name" form:"book_name"`
	//所属项目空间
	ItemId int `gorm:"column(item_id);type(int);default(1)" json:"item_id" form:"item_id"`
	// Identify 项目唯一标识.
	Identify string `gorm:"column(identify);size(100);unique" json:"identify" form:"identify"`
	//是否是自动发布 0 否/1 是
	AutoRelease int `gorm:"column(auto_release);type(int);default(0)" json:"auto_release" form:"auto_release"`
	//是否开启下载功能 0 是/1 否
	IsDownload int `gorm:"column(is_download);type(int);default(0)"json:"is_download"form:"is_download"`
	OrderIndex int `gorm:"column(order_index);type(int);default(0)" json:"order_index" form:"order_index"`
	// Description 项目描述.
	Description string `gorm:"column(description);size(2000)"json:"description"form:"description"`
	//发行公司
	Publisher string `gorm:"column(publisher);size(500)" json:"publisher" form:"publisher"`
	Label     string `gorm:"column(label);size(500)"json:"label"form:"label"`
	// PrivatelyOwned 项目私有： 0 公开/ 1 私有
	PrivatelyOwned int `gorm:"column(privately_owned);type(int);default(0)" json:"privately_owned" form:"privately_owned"`
	// 当项目是私有时的访问Token.
	PrivateToken string `gorm:"column(private_token);size(500);null" json:"private_token" form:"private_token"`
	//访问密码.
	BookPassword string `gorm:"column(book_password);size(500);null" json:"book_password" form:"book_password"`
	//状态：0 正常/1 已删除
	Status int `gorm:"column(status);type(int);default(0)" json:"status" form:"status"`
	//默认的编辑器.
	Editor string `gorm:"column(editor);size(50)" json:"editor" form:"editor"`
	// DocCount 包含文档数量.
	DocCount int `gorm:"column(doc_count);type(int)" json:"doc_count" form:"doc_count"`
	// CommentStatus 评论设置的状态:open 为允许所有人评论，closed 为不允许评论, group_only 仅允许参与者评论 ,registered_only 仅允许注册者评论.
	CommentStatus string `gorm:"column(comment_status);size(20);default(open)" json:"comment_status" form:"comment_status"`
	CommentCount  int    `gorm:"column(comment_count);type(int)" json:"comment_count" form:"comment_count"`
	//封面地址
	Cover string `gorm:"column(cover);size(1000)" json:"cover" form:"cover"`
	//主题风格
	Theme string `gorm:"column(theme);size(255);default(default)" json:"theme" form:"theme"`
	//每个文档保存的历史记录数量，0 为不限制
	HistoryCount int `gorm:"column(history_count);type(int);default(0)" json:"history_count" form:"history_count"`
	//是否启用分享，0启用/1不启用
	IsEnableShare int   `gorm:"column(is_enable_share);type(int);default(0)" json:"is_enable_share" form:"is_enable_share"`
	MemberId      int   `gorm:"column(member_id);size(100)" json:"member_id" form:"member_id"`
	Version       int64 `gorm:"type(bigint);column(version)" json:"version"`
	//是否使用第一篇文章项目为默认首页,0 否/1 是
	IsUseFirstDocument int `gorm:"column(is_use_first_document);type(int);default(0)" json:"is_use_first_document"`
	//是否开启自动保存：0 否/1 是
	AutoSave int `gorm:"column(auto_save);type(tinyint);default(0)" json:"auto_save"`
}

func RegisterDefaultBook(memberId int) error {
	var m Book
	m.MemberId = memberId
	m.BookName = "默认项目"
	m.Identify = "default"
	m.Description = "默认项目"
	m.AutoRelease = 1
	m.IsDownload = 1
	m.OrderIndex = 0
	m.PrivatelyOwned = 0
	m.PrivateToken = ""
	m.BookPassword = ""
	m.Status = 0
	m.Editor = "markdown"
	m.DocCount = 0
	m.CommentStatus = "open"
	m.CommentCount = 0
	m.Cover = ""
	m.Theme = "default"
	m.HistoryCount = 0
	m.IsEnableShare = 0
	m.IsUseFirstDocument = 0
	m.AutoSave = 0
	m.Version = time.Now().Unix()

	if err := db.Create(&m).Error; err != nil {
		return err
	}
	return nil
}

func (book *Book) String() string {
	ret, err := json.Marshal(*book)

	if err != nil {
		return ""
	}
	return string(ret)
}

// TableName 获取对应数据库表名.
func (book *Book) TableName() string {
	return "books"
}

// TableEngine 获取数据使用的引擎.
func (book *Book) TableEngine() string {
	return "INNODB"
}
func (book *Book) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + book.TableName()
}

//new book
func NewBook() *Book {
	return &Book{}
}

//添加项目
func (b *Book) Insert(lang string) error {
	b.BookName = utils.StripTags(b.BookName)
	if b.ItemId <= 0 {
		b.ItemId = 1
	}
	err := db.Create(b).Error
	if err != nil {
		if b.Label != "" {
			//插入label,待实现
		}

		relationship := NewRelationship()
		relationship.BookId = int(b.ID)
		relationship.RoleId = 0
		relationship.MemberId = b.MemberId
		err = relationship.Insert()
		if err != nil {
			return err
		}
	}
	return nil
}

//查找指定字段的项目
func (book *Book) Find(id int, cols ...string) (*Book, error) {
	if id <= 0 {
		return nil, errors.New("id is empty")
	}
	err := db.Select(cols).First(book, id).Error
	if err != nil {
		return nil, err
	}
	return book, nil
}

//更新项目
func (b *Book) Update(cols ...string) error {
	b.BookName = utils.StripTags(b.BookName)
	temp := NewBook()
	temp.ID = b.ID
	if err := db.First(temp).Error; err != nil {
		return err
	}
	if b.Label != "" || temp.Label != "" {
		//更新label,待实现
	}
	err := db.Save(b).Error
	if err != nil {
		return err
	}
	return nil
}

//复制项目
func (b *Book) Copy(identify string) error {
	if err := db.Where("identify = ?", identify).First(b).Error; err != nil {
		log.Println("查询项目出错")
		return err
	}
	bookId := b.ID
	b.Identify = b.Identify + fmt.Sprintf("%s-%s", identify, strconv.FormatInt(time.Now().UnixNano(), 32))
	b.BookName = b.BookName + "[副本]"
	b.CommentCount = 0
	b.HistoryCount = 0
	if err := db.Create(b).Error; err != nil {
		log.Println("复制项目出错")
		return err
	}
	//复制项目成员
	relationship := []*Relationship{}
	if err := db.Where("book_id = ?", bookId).Find(&relationship).Error; err != nil {
		log.Println("复制项目关系出错")
		return err
	}
	for _, rel := range relationship {
		rel.BookId = int(b.ID)
		rel.RelationshipId = 0
		if err := db.Create(&rel).Error; err != nil {
			log.Println("复制项目关系出错")
			return err
		}
	}
	docs := []*Document{}
	if err := db.Where("book_id = ?", bookId).Find(&docs).Error; err != nil {
		log.Println("复制项目文档出错")
		return err
	}
	if len(docs) > 0 {
		if err := recursiveInsertDocument(docs, int(b.ID), 0); err != nil {
			log.Println("复制项目文档出错")
			return err
		}
	}
	return nil
}

//递归的复制文档
func recursiveInsertDocument(docs []*Document, bookId int, parentId int) error {
	for _, doc := range docs {
		docId := doc.DocumentId
		doc.DocumentId = 0
		doc.ParentId = parentId
		doc.BookId = bookId
		doc.Version = time.Now().Unix()
		if err := db.Create(&doc).Error; err != nil {
			log.Println("插入项目文档出错")
			return err
		}
		var attachList []*Attachment
		if err := db.Where("document_id = ?", docId).Find(&attachList).Error; err != nil {
			log.Println("读取项目文档附件出错")
			return err
		} else {
			for _, attach := range attachList {
				attach.BookId = bookId
				attach.DocumentId = doc.DocumentId
				attach.AttachmentId = 0
				if err := db.Create(&attach).Error; err != nil {
					log.Println("插入项目文档附件出错")
					return err
				}
			}
		}

		var subdocs []*Document
		if err := db.Where("parent_id = ?", docId).Find(&subdocs).Error; err != nil {
			log.Println("读取项目文档子文档出错")
			return err
		}
		if len(subdocs) > 0 {
			if err := recursiveInsertDocument(subdocs, bookId, doc.DocumentId); err != nil {
				log.Println("递归插入项目文档子文档出错")
				return err
			}
		}
	}
	return nil
}

//查找项目的RoleId
func (b *Book) FindForRoleId(bookId, memberId int) (int, error) {
	var relationship Relationship
	err := db.Where("book_id = ? AND member_id = ?", bookId, memberId).First(&relationship).Error
	if err != nil {
		return 0, err
	}

	return relationship.RoleId, nil

}

//分页查询指定用户的项目
func (book *Book) FindToPager(pageIndex, pageSize, memberId int, lang string) (books []*BookResult, totalCount int, err error) {

	//sql1 := "SELECT COUNT(book.book_id) AS total_count FROM " + book.TableNameWithPrefix() + " AS book LEFT JOIN " +
	//	relationship.TableNameWithPrefix() + " AS rel ON book.book_id=rel.book_id AND rel.member_id = ? WHERE rel.relationship_id > 0 "

	sql1 := `SELECT
count(*) AS total_count
FROM md_books AS book
  LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.member_id = ?
  left join (select book_id,min(role_id) as role_id
             from (select book_id,team_member_id,role_id
                   from md_team_relationship as mtr
                     left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )
					as t group by t.book_id)
			as team on team.book_id=book.book_id WHERE rel.role_id >= 0 or team.role_id >= 0`

	err = db.Raw(sql1, memberId, memberId).Scan(&totalCount).Error

	if err != nil {
		return
	}

	offset := (pageIndex - 1) * pageSize

	//sql2 := "SELECT book.*,rel.member_id,rel.role_id,m.account as create_name FROM " + book.TableNameWithPrefix() + " AS book" +
	//	" LEFT JOIN " + relationship.TableNameWithPrefix() + " AS rel ON book.book_id=rel.book_id AND rel.member_id = ?" +
	//	" LEFT JOIN " + relationship.TableNameWithPrefix() + " AS rel1 ON book.book_id=rel1.book_id  AND rel1.role_id=0" +
	//	" LEFT JOIN " + NewMember().TableNameWithPrefix() + " AS m ON rel1.member_id=m.member_id " +
	//	" WHERE rel.relationship_id > 0 ORDER BY book.order_index DESC,book.book_id DESC LIMIT " + fmt.Sprintf("%d,%d", offset, pageSize)

	sql2 := `SELECT
  book.*,
  case when rel.relationship_id  is null then team.role_id else rel.role_id end as role_id,
  m.account as create_name
FROM md_books AS book
  LEFT JOIN md_relationship AS rel ON book.book_id = rel.book_id AND rel.member_id = ?
  left join (select book_id,min(role_id) as role_id
             from (select book_id,team_member_id,role_id
                   from md_team_relationship as mtr
                     left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )
					as t group by book_id) as team 
			on team.book_id=book.book_id
  LEFT JOIN md_relationship AS rel1 ON book.book_id = rel1.book_id AND rel1.role_id = 0
  LEFT JOIN md_members AS m ON rel1.member_id = m.member_id
WHERE rel.role_id >= 0 or team.role_id >= 0
ORDER BY book.order_index, book.book_id DESC limit ?,?`

	_, err = db.Raw(sql2, memberId, memberId, offset, pageSize).Scan(&books).Rows()
	if err != nil {
		log.Println("分页查询项目列表 => ", err)
		return
	}
	sql := "SELECT m.account,doc.modify_time FROM md_documents AS doc LEFT JOIN md_members AS m ON doc.modify_at=m.member_id WHERE book_id = ? LIMIT 1 ORDER BY doc.modify_time DESC"

	if len(books) > 0 {
		for index, book := range books {
			var text struct {
				Account    string
				ModifyTime time.Time
			}

			err1 := db.Raw(sql, book.BookId).Scan(&text).Error
			if err1 == nil {
				books[index].LastModifyText = text.Account + " 于 " + text.ModifyTime.Format("2006-01-02 15:04:05")
			}
			if book.RoleId == 0 {
				book.RoleName = i18n.Tr(lang, "common.creator")
			} else if book.RoleId == 1 {
				book.RoleName = i18n.Tr(lang, "common.administrator")
			} else if book.RoleId == 2 {
				book.RoleName = i18n.Tr(lang, "common.editor")
			} else if book.RoleId == 3 {
				book.RoleName = i18n.Tr(lang, "common.observer")
			}
		}
	}
	return
}
