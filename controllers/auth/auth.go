package auth

import (
	"fmt"
	"leedoc/utils"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func Test(c *gin.Context) {
	cookie, err := c.Cookie("gin_cookie")
	if err != nil {
		c.JSON(200, gin.H{})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "hello",
		"cookie":  cookie,
	})
}
func Auth(c *gin.Context) {
	username := c.PostForm("username")
	fmt.Println(username)
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}
	fmt.Println(user.Username, user.Password)
	if user.Username == "admin" && user.Password == "admin" {
		tokenString, _ := utils.Gentoken(user.Username)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "登录成功 hello",
			"token":   tokenString,
		})
		return
	}
	c.JSON(401, gin.H{
		"code":    401,
		"message": "用户名或密码错误",
	})
	return
}
