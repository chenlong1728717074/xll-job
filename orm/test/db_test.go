package test

import (
	"fmt"
	"testing"
	"xll-job/orm"
	"xll-job/orm/do"
)

func TestDb(t *testing.T) {
	//orm.DB.AutoMigrate(&do.JobLogDo{})
	//orm.DB.AutoMigrate(&do.JobInfoDo{})
	//orm.DB.AutoMigrate(&do.JobManagementDo{})
	//orm.DB.AutoMigrate(&do.JobLockDo{})
	orm.DB.AutoMigrate(&do.ExecutionLog{})
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
	var m do.JobManagementDo
	orm.DB.First(&m, 8166072975360)
	fmt.Println(m)
}

func TestDelete(t *testing.T) {
	lock := do.JobLockDo{
		Id: 0,
	}
	orm.DB.Delete(lock)
}
