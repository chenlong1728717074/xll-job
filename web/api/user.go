package api

import (
	"github.com/gin-gonic/gin"
)

type UserApi struct {
	router *gin.RouterGroup
}

func NewUserApi(router *gin.RouterGroup) *UserApi {
	return &UserApi{
		router: router,
	}
}

func (api *UserApi) Router() {
	api.router.POST("/add", api.add)
	api.router.GET("/delete", api.delete)
	api.router.POST("/update", api.update)
	api.router.GET("/getById", api.getById)
	api.router.GET("/page", api.page)
	api.router.GET("/restPassword", api.restPassword)
}

func (api *UserApi) add(context *gin.Context) {

}

func (api *UserApi) delete(context *gin.Context) {

}

func (api *UserApi) update(context *gin.Context) {

}

func (api *UserApi) getById(context *gin.Context) {

}

func (api *UserApi) page(context *gin.Context) {

}

func (api *UserApi) restPassword(context *gin.Context) {

}
