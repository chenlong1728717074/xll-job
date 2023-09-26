package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xll-job/orm"
	"xll-job/orm/do"
	"xll-job/utils"
	"xll-job/web/dto"
	"xll-job/web/middlewares"
)

type AuthApi struct {
	router *gin.RouterGroup
}

func NewAuthApi(router *gin.RouterGroup) *AuthApi {
	return &AuthApi{
		router: router,
	}
}

func (api *AuthApi) Router() {
	api.router.POST("/login", api.login)
	api.router.GET("/current", middlewares.JWTAuth(), api.current)
}

func (api *AuthApi) login(context *gin.Context) {
	var login dto.LoginDto
	if err := context.BindJSON(&login); err != nil {
		context.JSON(http.StatusOK, dto.NewResponse(http.StatusInternalServerError, "参数获取失败", ""))
		context.Abort()
		return
	}
	var user do.UserDo
	orm.DB.Where("user_name = ?", login.UserName).First(&user)
	if user.Id == 0 {
		context.JSON(http.StatusOK, dto.NewResponse(http.StatusInternalServerError, "用户不存在", ""))
		context.Abort()
		return
	}
	if !utils.CheckPassword(login.Password, user.Salt, user.Password) {
		context.JSON(http.StatusOK, dto.NewResponse(http.StatusInternalServerError, "密码错误", ""))
		context.Abort()
		return
	}
	claims := utils.NewCustomClaims(user.Id, user.UserName, user.Role)
	token, _ := utils.Jwt.CreateToken(claims)
	context.JSON(http.StatusOK, dto.NewOkResponse(token))
}

func (api *AuthApi) current(context *gin.Context) {

}
