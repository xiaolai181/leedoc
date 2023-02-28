package models

import (
	"errors"
	"leedoc/conf"
	"leedoc/utils/filetil"
	"log"
	"os"
	"strings"
	"time"
)

// AttachmentModel struct .
type Attachment struct {
	AttachmentId int     `gorm:"column(attachment_id);primarykey;auto;unique" json:"attachment_id"` //附件ID
	BookId       int     `gorm:"column(book_id);type(int)" json:"book_id"`                          //项目ID
	DocumentId   int     `gorm:"column(document_id);type(int);null" json:"doc_id"`                  //文档ID
	FileName     string  `gorm:"column(file_name);size(255)" json:"file_name"`                      //文件名
	FilePath     string  `gorm:"column(file_path);size(2000)" json:"file_path"`                     //文件路径
	FileSize     float64 `gorm:"column(file_size);type(float)" json:"file_size"`                    //文件大小
	HttpPath     string  `gorm:"column(http_path);size(2000)" json:"http_path"`                     //http路径
	FileExt      string  `gorm:"column(file_ext);size(50)" json:"file_ext"`                         //文件后缀
	CreateID     int     `gorm:"column(create_id);type(int)" json:"create_id"`                      //创建人ID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

//获取创建Attachment表名
func (a *Attachment) TableName() string {
	return "attachment"
}

// TableEngine 获取数据使用的引擎.
func (m *Attachment) TableEngine() string {
	return "INNODB"
}

//获取TableNameWithPrefix
func (m *Attachment) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

//创建Attachment
func NewAttachment() *Attachment {
	return &Attachment{}
}
func (m *Attachment) Insert() error {
	return db.Create(m).Error
}
func (m *Attachment) Update() error {
	return db.Save(m).Error
}
func (m *Attachment) Delete() error {
	err := db.Delete(m).Error
	if err == nil {
		if err = os.Remove(m.FilePath); err != nil {
			log.Println("删除附件失败:", err)
			return err
		}
	}
	return err
}

//根据附件ID获取附件
func (m *Attachment) Find(id int) (*Attachment, error) {
	if id <= 0 {
		return m, errors.New("参数错误")
	}
	err := db.Where("attachment_id = ?", id).First(m).Error
	return m, err
}

//查询文档的附件列表
func (m *Attachment) FindAttachmentsByDocumentId(docId int) ([]*Attachment, error) {
	if docId <= 0 {
		return nil, errors.New("参数错误")
	}
	var list []*Attachment
	db.Where("document_id = ?", docId).Find(&list)
	return list, nil
}

//分页查询附件
func (m *Attachment) FindAttachments(page int, pageSize int) (attachList []*AttachmentResult, err error) {
	var total int64
	db.Find(&Attachment{}).Count(&total)
	offset := (page - 1) * pageSize
	var list []*Attachment
	err = db.Order("attachment_id desc").Limit(pageSize).Offset(offset).Find(&list).Error
	if err != nil {
		log.Println("分页查询附件列表失败:", err)
		return nil, err
	}
	for _, item := range list {
		attach := &AttachmentResult{}
		attach.Attachment = *item
		attach.FileShortSize = filetil.FormatBytes(int64(item.FileSize))
		//查询附件所属的项目,ID为0表示是文档的附件
		if item.BookId == 0 && item.DocumentId > 0 {
			blog := NewBlog()
			if err := db.Where("bolg_id = ?", item.DocumentId).First(blog).Error; err == nil {
				attach.BookName = blog.BlogTitle
			} else {
				attach.BookName = "[文章不存在]"
			}
		} else {
			book := NewBook()
			if err := db.Where("book_id = ?", item.BookId).First(book).Error; err == nil {
				attach.BookName = book.BookName
				doc := NewDocument()
				if err := db.Where("document_id = ?", item.DocumentId).First(doc).Error; err == nil {
					attach.DocumentName = doc.DocumentName
				} else {
					attach.DocumentName = "[文档不存在]"
				}
			} else {
				attach.BookName = "[项目不存在]"
			}

		}
		attach.LocalHttpPath = strings.Replace(item.FilePath, "\\", "/", -1)
		attachList = append(attachList, attach)
	}
	return
}
