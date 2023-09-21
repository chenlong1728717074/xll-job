package do

import "xll-job/orm"

type JobInfoDo struct {
	orm.BaseModel
	ManageId   int64
	JobName    string `gorm:"type:varchar(64)"`
	JobHandler string `gorm:"type:varchar(64)"`
	Cron       string `gorm:"type:varchar(64)"`
	Enable     bool   `gorm:"column:is_enable;default:false;not null"`
}

func (JobInfoDo) TableName() string {
	return "tb_job_info"
}
