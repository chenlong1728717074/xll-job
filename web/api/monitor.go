package api

import (
	"github.com/gin-gonic/gin"
	"runtime"
)

type MonitorApi struct {
	router *gin.RouterGroup
}

func NewMonitorApi(router *gin.RouterGroup) *MonitorApi {
	return &MonitorApi{
		router: router,
	}
}

func (api *MonitorApi) Router() {
	api.router.POST("/basicInformation", api.basicInformation)
	api.router.POST("/goroutine", api.goroutine)
}

func (api *MonitorApi) basicInformation(context *gin.Context) {

}

func (api *MonitorApi) goroutine(context *gin.Context) {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	context.JSON(200, map[string]interface{}{
		"msg": string(buf),
	})
}
