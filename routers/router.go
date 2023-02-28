package routers

import (
	"leedoc/controllers"
	"leedoc/controllers/auth"
	"leedoc/funcmap"
	"leedoc/middleware"

	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

func Init_router() *gin.Engine {
	router := gin.Default()
	router.SetFuncMap(funcmap.FuncMap)
	router.LoadHTMLGlob("views/**/*")
	router.Static("/static", "static")
	router.Use(favicon.New("static/favicon.ico"))
	router.GET("/login", controllers.Login)
	router.GET("/register", controllers.Register)
	router.GET("/book", controllers.Book_Index)

	router.POST("/register", controllers.Register_post)
	router.POST("/auth", controllers.Auth)

	api := router.Group("/")
	api.Use(middleware.CookieAuth())
	{
		api.GET("/", controllers.Home)
		api.GET("/ping", Ping)
		api.GET("/test", auth.Test)
	}

	routers := router.Routes()
	funcmap.URL(routers)
	return router
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong" + c.Param("name"),
	})
}
