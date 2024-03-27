package route

import (
	"github.com/gin-gonic/gin"
	"github.com/hjrbill/quicker/demo/handler"
	"github.com/hjrbill/quicker/demo/middleware"
	"net/http"
)

func SetRoute(engine *gin.Engine) {
	engine.Static("js", "demo/views/js")
	engine.Static("css", "demo/views/css")
	engine.Static("img", "demo/views/img")
	engine.StaticFile("/favicon.ico", "img/dqq.png")                            //在 url 中访问文件/favicon.ico，相当于访问文件系统中的 views/img/dqq.png 文件
	engine.LoadHTMLFiles("demo/views/search.html", "demo/views/up_search.html") //使用这些.html 文件时就不需要加路径了

	engine.Use(middleware.GetUserInfo)                                                             //全局中间件
	classes := [...]string{"资讯", "社会", "热点", "生活", "知识", "环球", "游戏", "综合", "日常", "影视", "科技", "编程"} //数组，非切片
	engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "search.html", classes)
	})
	engine.GET("/up", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "up_search.html", classes)
	})

	// engine.POST("/search", handler.Search)
	engine.POST("/search", handler.SearchAll)
	engine.POST("/up_search", handler.SearchByAuthor)
}
