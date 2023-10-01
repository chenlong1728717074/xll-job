package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xll-job/web/api"
	"xll-job/web/middlewares"
)

func Router(engine *gin.Engine) {
	engine.LoadHTMLGlob("./static/html/*")
	engine.Static("/static", "./static/")
	//fileServer := http.FileServer(http.Dir("./web/static"))
	engine.GET("/", func(c *gin.Context) {
		// 渲染 "./web/static/index.html" 页面
		c.HTML(http.StatusOK, "index.html", "首页")
	})
	managementApi := api.NewJobManagementApi(engine.Group("/jobManagement", middlewares.JWTAuth()))
	managementApi.Router()
	infoApi := api.NewJobInfoApi(engine.Group("/jobInfo", middlewares.JWTAuth()))
	infoApi.Router()
	userApi := api.NewUserApi(engine.Group("/user", middlewares.JWTAuth(), middlewares.AdminAuth()))
	userApi.Router()
	authApi := api.NewAuthApi(engine.Group("/auth"))
	authApi.Router()
	monitorApi := api.NewMonitorApi(engine.Group("/monitor"))
	monitorApi.Router()
}
