package controllers

import (
	"fmt"
	"leedoc/models"
	"log"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	var b models.BookResult
	member, ok := c.Get("user")
	if !ok {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "请求参数错误,用户不存在",
		})

		return
	}

	member_id, _ := models.NewMember().FindIDByAccount(fmt.Sprint(member))
	log.Println(member_id)
	test, err := b.FindBooksByMemberId(member_id)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}
	c.HTML(200, "index.tpl", gin.H{
		"SITE_NAME":        "leedoc",
		"site_description": "leedoc is a simple blog system",
		"Lists":            test,
		"Lang":             "zh-cn",
		"Member":           Member{1, 1},
		"TotalPages":       1,
	})
}

type Member struct {
	MemberId int `json:"member_id"`
	Role     int `json:"role"`
}

type book struct {
	Identify   string `json:"identify"`
	BookName   string `json:"book_name"`
	CreateName string `json:"create_name"`
	Cover      string `json:"cover"`
}
