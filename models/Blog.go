package models

import (
	"bytes"
	"leedoc/conf"
	"leedoc/utils"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"leedoc/cache"

	"github.com/unknwon/com"
	"gorm.io/gorm"
)

//博文表
type Blog struct {
	gorm.Model
	BlogId int `gorm:"pk;auto;unique;column(blog_id)" json:"blog_id"`
	//文章标题
	BlogTitle string `gorm:"column(blog_title);size(500)" json:"blog_title"`
	//文章标识
	BlogIdentify string `gorm:"column(blog_identify);size(100);unique" json:"blog_identify"`
	//排序序号
	OrderIndex int `gorm:"column(order_index);type(int);default(0)" json:"order_index"`
	//所属用户
	MemberId int `gorm:"column(member_id);type(int);default(0);index" json:"member_id"`
	//用户头像
	MemberAvatar string `gorm:"-" json:"member_avatar"`
	//文章类型:0 普通文章/1 链接文章
	BlogType int `gorm:"column(blog_type);type(int);default(0)" json:"blog_type"`
	//链接到的项目中的文档ID
	DocumentId int `gorm:"column(document_id);type(int);default(0)" json:"document_id"`
	//文章的标识
	DocumentIdentify string `gorm:"-" json:"document_identify"`
	//关联文档的项目标识
	BookIdentify string `gorm:"-" json:"book_identify"`
	//关联文档的项目ID
	BookId int `gorm:"-" json:"book_id"`
	//文章摘要
	BlogExcerpt string `gorm:"column(blog_excerpt);size(1500)" json:"blog_excerpt"`
	//文章内容
	BlogContent string `gorm:"column(blog_content);type(text);null" json:"blog_content"`
	//发布后的文章内容
	BlogRelease string `gorm:"column(blog_release);type(text);null" json:"blog_release"`
	//文章当前的状态，枚举enum(’publish’,’draft’,’password’)值，publish为已 发表，draft为草稿，password 为私人内容(不会被公开) 。默认为publish。
	BlogStatus string `gorm:"column(blog_status);size(100);default(publish)" json:"blog_status"`
	//文章密码，varchar(100)值。文章编辑才可为文章设定一个密码，凭这个密码才能对文章进行重新强加或修改。
	Password string `gorm:"column(password);size(100)" json:"-"`
	//修改人id
	ModifyID       int    `gorm:"column(modify_at);type(int)" json:"-"`
	ModifyRealName string `gorm:"-" json:"modify_real_name"`
	//创建人
	CreateName string `gorm:"-" json:"create_name"`
	//版本号
	Version int64 `gorm:"type(bigint);column(version)" json:"version"`
	//附件列表
	AttachList []*Attachment `gorm:"-" json:"attach_list"`
}

// 多字段唯一键
func (b *Blog) TableUnique() [][]string {
	return [][]string{
		{"blog_id", "blog_identify"},
	}
}

// TableName 获取对应数据库表名.
func (b *Blog) TableName() string {
	return "blogs"
}

// TableEngine 获取数据使用的引擎.
func (b *Blog) TableEngine() string {
	return "INNODB"
}

func (b *Blog) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + b.TableName()
}

func NewBlog() *Blog {
	return &Blog{
		BlogStatus: "public",
	}
}

//根据文章ID获取文章
func (b *Blog) Find(blogId int) (*Blog, error) {
	err := db.Where("blog_id=?", blogId).First(&b).Error
	if err != nil {
		return nil, err
	}
	return b.Link()
}

//从缓存中获取文章
func (b *Blog) FindBlogCache(blogId int) (blog *Blog, err error) {
	key := "blog-id-" + com.ToStr(blogId)
	var temp Blog
	err = cache.Get(key, &temp)
	if err == nil {
		b = &temp
		b.Link()
		log.Println("get blog from cache sucess ->", key)
		return b, nil
	} else {
		log.Println("get blog from cache failed ->", key)
	}
	blog, err = b.Find(blogId)
	if err == nil {
		//缓存文章
		if err = cache.Put(key, blog, 0); err != nil {
			log.Println("cache blog failed ->", key)
		} else {
			log.Println("cache blog sucess ->", key)
		}
	}
	return
}

//查找用户的指定文章
func (b *Blog) FindUserBlogs(memberId int, blogId int) (blog *Blog, err error) {
	err = db.Where("member_id=? AND blog_id=?", memberId, blogId).First(&b).Error
	if err != nil {
		log.Println("get user blog failed ->", memberId, blogId)
		return nil, err
	}
	return b.Link()
}

//根据文章标识获取文章
func (b *Blog) FindByIdentify(identify string) (*Blog, error) {
	err := db.Where("blog_identify=?", identify).First(&b).Error
	if err != nil {
		log.Println("get blog by identify failed ->", identify)
		return nil, err
	}
	return b.Link()
}

//获取文章的链接内容
//BlogType 1:文档 2:博客 3:项目
func (b *Blog) Link() (*Blog, error) {
	if b.BlogType == 1 && b.DocumentId > 0 {
		doc := NewDocument()
		if err := db.Where("document_id=?", b.DocumentId).First(&doc).Error; err != nil {
			log.Println("get document failed ->", b.DocumentId)
		} else {
			b.DocumentId = doc.DocumentId
			b.BlogRelease = doc.Release
			//
			b.BlogContent = doc.Markdown
			book := NewBook()
			if err := db.Where("book_id=?", doc.BookId).First(&book).Error; err != nil {
				log.Println("get book failed ->", doc.BookId)
			} else {
				b.BookIdentify = book.Identify
				b.BookId = int(book.ID)
			}
			if content, err := goquery.NewDocumentFromReader(bytes.NewBufferString(b.BlogRelease)); err == nil {
				content.Find(".wiki-bottom").Remove()
				if html, err := content.Html(); err == nil {
					b.BlogContent = html
				} else {
					log.Println("get blog content failed ->", err)

				}
			} else {
				log.Println("get blog content failed ->", err)
			}
		}

	}
	if b.ModifyID > 0 {
		member := NewMember()
		if err := db.Where("member_id=?", b.ModifyID).First(&member).Error; err == nil {
			if member.RealName != "" {
				b.ModifyRealName = member.RealName
			} else {
				b.ModifyRealName = member.Account
			}
		}
	}
	if b.MemberId > 0 {
		member := NewMember()
		if err := db.Where("member_id=?", b.MemberId).First(&member).Error; err == nil {
			if member.RealName != "" {
				b.CreateName = member.RealName
			} else {
				b.CreateName = member.Account
			}
			b.MemberAvatar = member.Avatar
		}
	}
	return b, nil
}

//判断文章是否存在
func (b *Blog) IsExist(identify string) bool {
	var i int64
	db.Where("blog_identify=?", identify).Count(&i)
	return i > 0
}
func (b *Blog) Save(cols ...string) error {
	if b.OrderIndex <= 0 {
		blog := NewBlog()
		if err := db.Order("blog_id desc").First(&blog).Error; err == nil {
			b.OrderIndex = blog.OrderIndex + 1
		} else {
			var c int64
			db.Find(&Blog{}).Count(&c)
			b.OrderIndex = int(c) + 1
		}
	}
	var err error
	b.Processor().Version = time.Now().Unix()
	if b.BlogId > 0 {
		b.UpdatedAt = time.Now()
		err = db.Save(b).Error
		key := "blog-id-" + com.ToStr(b.BlogId)
		cache.Delete(key)
	} else {
		b.CreatedAt = time.Now()
		err = db.Create(b).Error
	}
	return err
}

//过滤文章内容,处理文章内容中的链接和图片
func (b *Blog) Processor() *Blog {
	b.BlogRelease = utils.SafetyProcessor(b.BlogRelease)
	if content, err := goquery.NewDocumentFromReader(bytes.NewBufferString(b.BlogRelease)); err == nil {

		content.Find("a").Each(func(i int, contentSelection *goquery.Selection) {
			if src, ok := contentSelection.Attr("href"); ok {
				if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
					//logs.Info(src,conf.BaseUrl,strings.HasPrefix(src,conf.BaseUrl))
					if conf.BaseUrl != "" && !strings.HasPrefix(src, conf.BaseUrl) {
						contentSelection.SetAttr("target", "_blank")
						if html, err := content.Html(); err == nil {
							b.BlogRelease = html
						}
					}
				}

			}
		})
		//设置图片为CDN地址
		if cdnimg := conf.GetConfig("cdnimg", "").(string); cdnimg != "" {
			content.Find("img").Each(func(i int, contentSelection *goquery.Selection) {
				if src, ok := contentSelection.Attr("src"); ok && strings.HasPrefix(src, "/uploads/") {
					contentSelection.SetAttr("src", utils.JoinURI(cdnimg, src))
				}
			})
		}
	}

	return b
}

//查询下一篇文章*
func (b *Blog) QueryNext(blogId int) (*Blog, error) {
	if err := db.Where("blog_id>?", blogId).Order("blog_id asc").First(&b).Error; err != nil {
		log.Println("get next blog failed ->", err)
		return nil, err
	}
	return b, nil
}

//分页查询文章列表
func (b *Blog) FindToPager(pageIndex, pageSize int, memberId int, status string) (blogList []*Blog, totalCount int, err error) {
	offet := (pageIndex - 1) * pageSize

	query := db.Order("order_index,blog_id")
	if memberId > 0 {
		query = query.Where("member_id=?", memberId)
	}
	if status != "" {
		query = query.Where("status=?", status)
	}
	var count int64
	err = query.Offset(offet).Limit(pageSize).Find(&blogList).Count(&count).Error
	totalCount = int(count)
	if err != nil {
		if count == 0 {
			err = nil
		}
		log.Println("get blog list failed ->", err)
		return
	}
	if count <= 0 {
		log.Println("获取文章数量 ->", err)
		return nil, 0, err
	}
	for _, blog := range blogList {
		if blog.BlogType == 1 {
			blog.Link()
		}
	}

	return
}

//关联文章附件
func (b *Blog) LinkAttach() (err error) {
	var attachList []*Attachment
	if b.BlogType != 1 || b.BlogId <= 0 {
		var count int64
		err := db.Where("blog_id=?", b.BlogId).Find(&attachList).Count(&count).Error
		if err != nil && count <= 0 {
			log.Println("未查询文章附件 ->", err)
		} else {
			err := db.Where("document_id=?", b.DocumentId).Where("book_id=?", b.BookId).Find(&attachList).Count(&count).Error
			if err != nil && count <= 0 {
				log.Println("未查询文章附件 ->", err)
			}
		}
	}
	b.AttachList = attachList
	return
}
