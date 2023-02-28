package controllers

import (
	"fmt"
	"leedoc/models"
	"leedoc/utils"

	"github.com/gin-gonic/gin"
)

//生成xsrfdata
func Xsrfdata(Ip string) string {
	xsrfdata, _ := utils.GenXsrf(Ip)
	return xsrfdata
}

//用户登录
func Login(c *gin.Context) {
	a := Base(c)
	template := "login.tpl"
	c.HTML(200, template, a)
}

func Auth(c *gin.Context) {

	var user models.Member
	if err := c.Bind(&user); err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}
	fmt.Println(user.Account, user.Password)

	if ok, err := user.Login(); ok && err == nil {
		tokenString, _ := utils.Gentoken(user.Account)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "登录成功 hello",
			"token":   tokenString,
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    401,
		"message": "用户名或密码错误",
	})
	return
}

// 用户注册页面
func Register(c *gin.Context) {
	a := Base(c)
	c.HTML(200, "register.tpl", a)
}

//注册,由管理员添加
func Register_post(c *gin.Context) {
	var member models.Member
	password1 := c.PostForm("password1")
	password2 := c.PostForm("password2")
	if err := c.Bind(&member); err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}
	if password1 != password2 {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "两次密码不一致",
		})
		return
	}
	member.Password = (password1)
	fmt.Println(member.Account, member.Password)
	if member.Account == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "账号不能为空",
		})
		return
	}
	if member.Password == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "密码不能为空",
		})
		return
	}
	if member.Email == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "邮箱不能为空",
		})
		return
	}
	err := member.Add()

	if err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": fmt.Sprintf("注册失败:%s", err.Error()),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "注册成功",
	})

}
