package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
	"xll-job/web/api"
)

func Router(engine *gin.Engine) {
	engine.GET("/moniter", func(context *gin.Context) {
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, true)
		fmt.Printf("%s", buf)
		context.JSON(200, map[string]interface{}{
			"msg": string(buf),
		})
	})
	managementApi := api.NewJobManagementApi(engine.Group("/jobManagement"))
	managementApi.Router()
	infoApi := api.NewJobInfoApi(engine.Group("/jobInfo"))
	infoApi.Router()
}
