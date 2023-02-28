package models

import (
	"leedoc/utils/filetil"
	"strings"
)

//附件结果
type AttachmentResult struct {
	Attachment
	IsExist       bool   //是否存在
	BookName      string //项目名称
	DocumentName  string //文档名称
	FileShortSize string //文件短尺寸
	Account       string //上传者
	LocalHttpPath string //本地路径
}

//new一个附件结果
func NewAttachmentResult() *AttachmentResult {
	return &AttachmentResult{IsExist: false}
}

//根据附件ID获取附件结果
func (m *AttachmentResult) Find(id int) (*AttachmentResult, error) {
	attch := NewAttachment()
	err := db.Where("attachment_id = ?", id).First(attch).Error
	if err != nil {
		return m, err
	}
	m.Attachment = *attch
	if m.Attachment.DocumentId > 0 && attch.BookId == 0 {
		blog := NewBlog()
		if err := db.Where("bolg_id = ?", attch.DocumentId).First(blog).Error; err == nil {
			m.BookName = blog.BlogTitle
		} else {
			m.BookName = "[文章不存在]"
		}
	} else {
		book := NewBook()
		if err := db.Where("book_id = ?", attch.BookId).First(book).Error; err == nil {
			m.BookName = book.BookName
			doc := NewDocument()
			if err := db.Where("document_id = ?", attch.DocumentId).First(doc).Error; err == nil {
				m.DocumentName = doc.DocumentName
			} else {
				m.DocumentName = "[文档不存在]"
			}
		} else {
			m.BookName = "[项目不存在]"
		}
	}
	if attch.CreatedAt.IsZero() {
		member := NewMember()
		if err := db.Where("member_id = ?", attch.CreateID).First(member).Error; err == nil {
			m.Account = member.Account
		} else {
			m.Account = "[用户不存在]"
		}
	}
	m.FileShortSize = filetil.FormatBytes(int64(attch.FileSize))
	m.LocalHttpPath = strings.Replace(m.FilePath, "\\", "/", -1)
	return m, nil
}
