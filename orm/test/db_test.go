package test

import (
	"fmt"
	"testing"
	"time"
	"xll-job/orm"
	"xll-job/orm/bo"
	"xll-job/orm/constant"
	"xll-job/orm/do"
)

func TestDb(t *testing.T) {
	orm.DB.AutoMigrate(&do.JobLogDo{})
	//orm.DB.AutoMigrate(&do.JobInfoDo{})
	//orm.DB.AutoMigrate(&do.JobManagementDo{})
	//orm.DB.AutoMigrate(&do.JobLockDo{})
	//orm.DB.AutoMigrate(&do.ExecutionLog{})
}

func TestAdd(t *testing.T) {
	log := &do.JobLockDo{
		Id: 1,
	}
	tx := orm.DB.Create(log)
	//row := tx.Row()
	fmt.Println(tx.RowsAffected)
}
func TestSelect(t *testing.T) {
	losedTime := time.Now().Add(-time.Minute * 10)
	var jobLogs []do.JobLogDo
	orm.DB.Model(&do.JobLogDo{}).
		Where("dispatch_time < ? and execute_status = ?",
			losedTime,
			constant.InProgress,
		).Find(&jobLogs)
	fmt.Println(jobLogs)
}

func TestSelectSql(t *testing.T) {
	var jobLogs []bo.RetryJobBo
	orm.DB.Raw(constant.RetryJob).Scan(&jobLogs)
	fmt.Println(jobLogs)
	fmt.Println(len(jobLogs))
	fmt.Println(jobLogs[0].Enable)

}

func TestDelete(t *testing.T) {
	lock := do.JobLockDo{
		Id: 0,
	}
	orm.DB.Delete(lock)
}
