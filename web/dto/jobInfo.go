package dto

type JobInfoDto struct {
	Id         int64  `json:"id" form:"id" json:"id" uri:"id" xml:"id" yaml:"id" `
	ManageId   int64  `json:"manageId" form:"manageId" json:"manageId" uri:"manageId" xml:"manageId" yaml:"manageId" binding:"required"`
	JobName    string `json:"jobName" form:"jobName" json:"jobName" uri:"jobName" xml:"jobName" yaml:"jobName" binding:"required"`
	JobHandler string `json:"jobHandler" form:"jobHandler" json:"jobHandler" uri:"jobHandler" xml:"jobHandler" yaml:"jobHandler" binding:"required"`
	Cron       string `json:"core" form:"core" json:"core" uri:"core" xml:"core" yaml:"core" binding:"required"`
}
