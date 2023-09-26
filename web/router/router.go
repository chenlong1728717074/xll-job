package router

import (
	"github.com/gin-gonic/gin"
	"xll-job/web/api"
	"xll-job/web/middlewares"
)

func Router(engine *gin.Engine) {
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
