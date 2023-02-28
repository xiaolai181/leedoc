package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"leedoc/conf"
	"leedoc/utils"
	"leedoc/utils/pagination"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
)

// book主页
func Book_Index(c *gin.Context) {
	Data := Base(c)
	TplName := "book/index.tpl"
	pageIndex, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	Member, ok := c.Get("user")
	if !ok {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "请求参数错误,用户不存在",
		})

		return
	}
	member_id, _ := models.NewMember().FindIDByAccount(fmt.Sprint(Member))
	books, totalcount, err := models.NewBook().FindToPager(pageIndex, conf.PageSize, member_id, "zh-cn")
	if err != nil {
		log.Println(err)
		c.Abort()
	}
	for i, book := range books {
		books[i].Description = utils.StripTags(string(blackfriday.Run([]byte(book.Description))))
		books[i].ModifyTime = book.ModifyTime.Local()
		books[i].CreateTime = book.CreateTime.Local()
	}
	if totalcount > 0 {
		pager := pagination.NewPagination(c.Request, totalcount, conf.PageSize)
		Data["PageHtml"] = pager.HtmlPages()
	} else {
		Data["PageHtml"] = ""
	}
	b, err := json.Marshal(books)

	if err != nil || len(books) <= 0 {
		Data["Result"] = template.JS("[]")
	} else {
		Data["Result"] = template.JS(string(b))
	}
	if itemsets, err := models.NewItemsets().First(1); err == nil {
		Data["Item"] = itemsets
	}
	// a := map[string]interface{}{"book_id": 2, "book_name": "qwe", "item_id": 1, "item_name": "", "identify": "qwe", "order_index": 0, "description": "", "publisher": "", "privately_owned": 0, "private_token": "", "book_password": "", "doc_count": 1, "comment_status": "closed", "comment_count": 0, "create_time": "2022-07-25T19:58:26.9650603+08:00", "create_name": "admin", "real_name": "", "modify_time": "2022-07-25T19:58:26.9650603+08:00", "cover": "/uploads/202207/cover_1705109da1a9a43c.png", "theme": "default", "label": "", "member_id": 1, "editor": "markdown", "auto_release": false, "history_count": 0, "role_id": 0, "role_name": "创始人", "status": 0, "is_enable_share": false, "is_use_first_document": true, "last_modify_text": "", "is_display_comment": false, "is_download": true, "auto_save": false, "Lang": ""}
	// b := Base(c)
	// s := []map[string]interface{}{a}
	// b["Result"] = s
	// b["Member"] = Member{1, 1}
	// b["ControllerName"] = "BookController"

	// c.HTML(200, "book/index.tpl", b)
}
