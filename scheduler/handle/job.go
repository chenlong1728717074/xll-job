package handle

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"xll-job/orm"
	"xll-job/orm/do"
	"xll-job/scheduler/core"
)

type XllJob struct {
	lock     sync.RWMutex
	Trigger  *cron.Cron
	Manager  map[int64]*core.JobManager
	register *Register
}

func NewXllJob() *XllJob {
	job := XllJob{
		lock:     sync.RWMutex{},
		register: NewRegister(),
	}
	return &job
}

// InitXllJob  不确定之后会不会有其他操作,所以先把初始化操作单独抽出/** **/
func (job *XllJob) InitXllJob() {
	//先初始化
	if job.Trigger == nil {
		job.Trigger = cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		)))
		//job.trigger = cron.New()
	}
	if job.Manager == nil {
		job.Manager = make(map[int64]*core.JobManager)
		JobManagerMap = job.Manager
	}
}
func (job *XllJob) Start() {
	job.Trigger.Start()
	job.register.Start()
}

func (job *XllJob) LoadJob() {
	/**不确定之后需不需要加锁,先预留代码*/
	/*job.lock.Lock()
	defer job.lock.Unlock()*/
	//先初始化
	job.InitXllJob()
	//加载所有任务管理器
	log.Println("开始加载任务管理器")
	var managers []do.JobManagementDo
	orm.DB.Model(&do.JobManagementDo{}).Find(&managers)
	manager := job.Manager
	for _, managementDo := range managers {
		jobManager := core.NewJobManager(managementDo.Id, managementDo.Name, managementDo.AppName)
		manager[managementDo.Id] = jobManager
		//加载任务
		var jobs []do.JobInfoDo
		orm.DB.Model(&do.JobInfoDo{}).Where(&do.JobInfoDo{ManageId: jobManager.Id}).Find(&jobs)
		if len(jobs) == 0 {
			continue
		}
		for _, infoDo := range jobs {
			if infoDo.Enable {
				scheduler, _ := core.NewScheduler(infoDo.Cron,
					infoDo.JobHandler, jobManager, true)
				jobManager.Schedulers[infoDo.Id] = scheduler
				//任务逻辑
				enId, _ := Xll_Job.Trigger.AddFunc(infoDo.Cron, scheduler.Execute)
				scheduler.TriggerId = enId
			}
		}
	}
	for _, v := range manager {
		fmt.Println(v.Schedulers)
	}
	log.Printf("任务管理器加载成功,size=%d\n", len(manager))
}
