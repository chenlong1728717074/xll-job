package core

type JobManager struct {
	Id         int64
	AppName    string
	Name       string
	ServerAddr []*ServiceNode
	Schedulers map[int64]*Scheduler
}

func NewJobManager(id int64, appName string, name string) *JobManager {
	manager := JobManager{
		Id:         id,
		AppName:    appName,
		Name:       name,
		ServerAddr: make([]*ServiceNode, 0), // 使用 make 创建新的切片
		Schedulers: make(map[int64]*Scheduler),
	}
	return &manager
}
