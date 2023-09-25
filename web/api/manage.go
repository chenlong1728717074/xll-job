package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
	"xll-job/orm"
	"xll-job/orm/do"
	"xll-job/scheduler/core"
	"xll-job/scheduler/handle"
	"xll-job/web/dto"
	"xll-job/web/vo"
)

type JobManagementApi struct {
	router *gin.RouterGroup
}

func NewJobManagementApi(router *gin.RouterGroup) *JobManagementApi {
	return &JobManagementApi{
		router: router,
	}
}

func (managementApi *JobManagementApi) Router() {
	managementApi.router.POST("/add", managementApi.Add)
	managementApi.router.GET("/getById", managementApi.GetById)
	managementApi.router.GET("/delete", managementApi.delete)
}

func (managementApi *JobManagementApi) GetById(ctx *gin.Context) {
	i := ctx.Query("id")
	//获取
	id, _ := strconv.ParseInt(i, 10, 64)
	manager := handle.JobManagerMap[id]
	addrs := make([]string, 0)
	for _, node := range manager.ServerAddr {
		addrs = append(addrs, node.Addr)
	}

	ctx.JSON(200, dto.NewOkResponse(vo.JobManagementVo{
		ID:            manager.Id,
		Name:          manager.Name,
		AppName:       manager.AppName,
		OnlineMachine: addrs,
	}))
}

func (managementApi *JobManagementApi) Add(ctx *gin.Context) {
	var jobManagementDto dto.JobManagementDto
	err := ctx.ShouldBindJSON(&jobManagementDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{
			Err: errors.New("校验失败"),
		})
		ctx.Done()
	}
	//保存
	var jobManagementDo do.JobManagementDo
	copier.Copy(&jobManagementDo, &jobManagementDto)
	orm.DB.Create(&jobManagementDo)
	//调用任务模块添加缓存
	manager := core.NewJobManager(jobManagementDo.Id, jobManagementDo.Name, jobManagementDo.AppName)
	// 存储新的 JobManager 实例到 缓存 中
	handle.JobManagerMap[manager.Id] = manager
	ctx.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (managementApi *JobManagementApi) delete(context *gin.Context) {
	i := context.Query("id")
	//获取
	id, _ := strconv.ParseInt(i, 10, 64)
	var manager do.JobManagementDo
	orm.DB.First(&manager, id)
	if manager.Id == 0 {
		context.JSON(500, gin.Error{
			Meta: "任务管理器不存在",
		})
		context.Done()
		return
	}
	var count int64
	orm.DB.Model(&do.JobInfoDo{}).Where("manage_id=? and is_enable=?", id, 1).Count(&count)
	if count != 0 {
		context.JSON(500, gin.Error{
			Meta: "该任务管理器下有已开启任务",
		})
		context.Done()
		return
	}
	orm.DB.Delete(&manager)
	managerMap := handle.JobManagerMap
	fmt.Println(managerMap)
	job := handle.Xll_Job
	fmt.Println(job)
	delete(handle.JobManagerMap, id)
	context.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}
