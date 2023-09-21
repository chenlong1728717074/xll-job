package core

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"xll-job/scheduler/grpc/dispatch"
)

type Scheduler struct {
	Id         int64
	TriggerId  cron.EntryID
	lock       sync.RWMutex
	jobManager *JobManager
	jobHandler string
	cron       string
	retry      int
	con        *grpc.ClientConn
}

func NewScheduler(expression string, jobHandler string, jobManager *JobManager, enable bool) (*Scheduler, error) {
	s := &Scheduler{
		retry:      3,
		cron:       expression,
		jobHandler: jobHandler,
		jobManager: jobManager,
	}
	s.lock = sync.RWMutex{}
	return s, nil
}
func (s *Scheduler) Execute() {
	//目前只能实现单机服务调用多台服务 后续移除robfig再实现集训调用
	//
	addr := s.jobManager.ServerAddr
	if len(addr) == 0 {
		return
	}
	fmt.Printf("开始调度:%s\n", s.jobHandler)
	//路由
	/*	if len(s.jobManager.ServerAddr) == 0 {
			return
		}
		dial, _ := grpc.Dial(s.jobManager.ServerAddr[0])*/
	//调用前插入日志记录

	//连接  后期这个连接将会在路由实现
	dial, err := grpc.Dial(addr[0].Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dial.Close()
	con := dispatch.NewServiceClient(dial)
	//调度
	con.Call(context.Background(), &dispatch.Request{
		ServiceId: s.jobHandler,
		Retry:     3,
	})
	//更改调度记录 0失败 1成功
}
