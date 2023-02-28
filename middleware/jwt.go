package middleware

import (
	"leedoc/utils"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

//jwt中间件
func JWTAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "请求头中无token",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "请求头中authtoken格式错误",
			})
			c.Abort()
			return
		}
		mc, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "token错误",
			})
			c.Abort()
			return
		}
		c.Set("user", mc.UserName)
		c.Next()
	}
}

//cookie中间件
func CookieAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		token, err := c.Cookie("gin_cookie")
		if err != nil || token == "" {
			c.Redirect(302, "/login")
			return
		}
		mc, err := utils.ParseToken(token)
		if err != nil {
			c.Redirect(302, "/login")
			return
		}
		c.Set("user", mc.UserName)
		log.Println("登陆", mc.UserName)
		c.Next()
	}
}
