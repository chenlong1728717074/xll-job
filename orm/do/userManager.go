package do

import "xll-job/orm"

type UserManager struct {
	orm.BaseModel
	UserId    int64
	ManagerId int64
}

func (UserManager) TableName() string {
	return "tb_user_manager"
}
