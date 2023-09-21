package do

type jobLock struct {
	Id int64 `gorm:"primary_key;auto_increment:false"`
}

func (jobLock) TableName() string {
	return "tb_job_lock"
}
