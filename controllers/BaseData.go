package controllers

import (
	"github.com/gin-gonic/gin"
)

func Base(c *gin.Context) gin.H {
	return gin.H{
		"SITE_NAME":        "leedoc",
		"site_description": "leedoc is a simple blog system",
		"BaseUrl":          "/",
		"xsrfdata":         Xsrfdata(c.ClientIP()),
		"Lang":             "zh-cn",
	}
}
