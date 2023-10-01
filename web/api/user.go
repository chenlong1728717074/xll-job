package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xll-job/orm"
	"xll-job/orm/do"
	"xll-job/utils"
	"xll-job/web/dto"
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

func (api *UserApi) add(ctx *gin.Context) {
	var userDto dto.UserDto
	ctx.ShouldBindJSON(&userDto)
	var user do.UserDo
	orm.DB.Model(&user).Where("user_name = ?", userDto.UserName).First(&user)
	if user.UserName != "" {
		ctx.JSON(200, dto.NewErrResponse("用户已存在", ""))
		ctx.Done()
		return
	}
	salt, code := utils.GeneratePassword(userDto.Password)
	user.Password = code
	user.Salt = salt
	user.UserName = userDto.UserName
	orm.DB.Create(&user)
	//先删除后bind
	orm.DB.Where("user_id=?", 1).Delete(&do.UserManager{})

	var ums []*do.UserManager
	for _, managerId := range userDto.Manager {
		ums = append(ums, &do.UserManager{
			UserId:    user.Id,
			ManagerId: managerId,
		})
	}
	orm.DB.Create(ums)
	ctx.JSON(200, dto.NewOkResponse(""))
}

func (api *UserApi) delete(ctx *gin.Context) {
	i := ctx.Query("id")
	//获取
	id, _ := strconv.ParseInt(i, 10, 64)
	var user do.UserDo
	orm.DB.First(&user, id)
	if user.Id != 0 {
		ctx.JSON(200, dto.NewErrResponse("条目不存在", ""))
		ctx.Done()
		return
	}
	orm.DB.Delete(&user)
	ctx.JSON(200, dto.NewOkResponse(""))
}

func (api *UserApi) update(ctx *gin.Context) {
	var userDto dto.UserDto
	ctx.ShouldBindJSON(&userDto)
	var user do.UserDo
	orm.DB.Model(&user).Where("user_name = ?", userDto.UserName).First(&user)
	if user.UserName != "" {
		ctx.JSON(200, dto.NewErrResponse("用户已存在", ""))
		ctx.Done()
		return
	}
	if len(userDto.Manager) == 0 {
		//先删除后bind
		orm.DB.Where("user_id=?", 1).Delete(&do.UserManager{})
		var ums []*do.UserManager
		for _, managerId := range userDto.Manager {
			ums = append(ums, &do.UserManager{
				UserId:    user.Id,
				ManagerId: managerId,
			})
		}
		orm.DB.Create(ums)
	}
	ctx.JSON(200, dto.NewOkResponse(""))
}

func (api *UserApi) getById(ctx *gin.Context) {

}

func (api *UserApi) page(ctx *gin.Context) {

}

func (api *UserApi) restPassword(ctx *gin.Context) {

}
