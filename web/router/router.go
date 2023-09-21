package router

import (
	"github.com/gin-gonic/gin"
	"xll-job/web/api"
)

func Router(engine *gin.Engine) {
	managementApi := api.NewJobManagementApi(engine.Group("/jobManagement"))
	managementApi.Router()
	infoApi := api.NewJobInfoApi(engine.Group("/jobInfo"))
	infoApi.Router()
}
