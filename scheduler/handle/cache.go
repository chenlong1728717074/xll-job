package handle

import (
	"xll-job/scheduler/core"
)

var Xll_Job *XllJobHandle

var JobManagerMap map[int64]*core.JobManager

var SchedulerMap map[int64]*core.Scheduler

var ServiceNodeList []*core.ServiceNode
