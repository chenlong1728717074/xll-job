package test

import (
	"fmt"
	"testing"
	"xll-job/orm"
	"xll-job/orm/do"
)

func TestDb(t *testing.T) {
	orm.DB.AutoMigrate(&do.JobLogDo{})
	orm.DB.AutoMigrate(&do.JobInfoDo{})
	orm.DB.AutoMigrate(&do.JobManagementDo{})
}

func TestAdd(t *testing.T) {
	log := &do.JobLogDo{
		ManageId: 1,
		JobId:    4,
	}
	orm.DB.Create(log)
}
func TestSelect(t *testing.T) {
	var m do.JobManagementDo
	orm.DB.First(&m, 8166072975360)
	fmt.Println(m)
}

func TestDelete(t *testing.T) {
}
