package core

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type Scheduler struct {
	Id            int64
	TriggerId     cron.EntryID
	lock          sync.RWMutex
	JobManager    *JobManager
	JobHandler    string
	cron          string
	RoutingPolicy int32
	Retry         int32
}

func NewScheduler(retry int32, expression string, jobHandler string, jobManager *JobManager, enable bool, routingPolicy int32) (*Scheduler, error) {
	s := &Scheduler{
		Retry:         retry,
		cron:          expression,
		JobHandler:    jobHandler,
		JobManager:    jobManager,
		RoutingPolicy: routingPolicy,
	}
	s.lock = sync.RWMutex{}
	return s, nil
}

// Execute At present, we can only implement single machine service call cluster deployment, and we will support cluster deployment in the future
func (s *Scheduler) Execute() {
	TriggerExecute(s)
}
