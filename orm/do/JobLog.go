package do

import (
	"time"
	"xll-job/orm"
)

type JobLogDo struct {
	orm.BaseModel
	ManageId             int64      `gorm:"column:manage_id;comment:管理器id"`
	JobId                int64      `gorm:"column:job_id;comment:任务id"`
	Dispatch             bool       `gorm:"column:is_dispatch;comment:是否调度成功"`
	DispatchTime         *time.Time `gorm:"column:dispatch_time;comment:调度时间"`
	DispatchAddress      string     `gorm:"column:dispatch_address;comment:调度地址;type:varchar(64)"`
	DispatchHandler      string     `gorm:"type:varchar(64);comment:调度handler"`
	DispatchStatus       int64      `gorm:"comment:调度状态 1:调度 2:调度失败"`
	retry                int64      `gorm:"default:0;comment:重试次数"`
	ExecuteStartTime     *time.Time `gorm:"comment:执行开始时间"`
	ExecuteEndTime       *time.Time `gorm:"comment:执行结束时间"`
	ExecuteConsumingTime int64      `gorm:"comment:执行耗时"`
	ExecuteType          int64      `gorm:"comment:执行类型"`
	ExecuteStatus        int64      `gorm:"comment:执行状态 0:进行中 1:执行成功 2:执行出现异常"`
	ExecuteLogs          string     `gorm:"comment:执行日志"`
}

func (JobLogDo) TableName() string {
	return "tb_job_log"
}
