package core

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	"xll-job/orm"
	"xll-job/orm/do"
	"xll-job/scheduler/grpc/dispatch"
)

func Execute(s *Scheduler) {
	//Is there a service node present
	addr := s.JobManager.ServerAddr
	if len(addr) == 0 {
		return
	}
	//service lock;Prevent parallel processing of tasks
	lock := &do.JobLockDo{
		Id: s.Id,
	}
	tx := orm.DB.Create(lock)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return
	}
	fmt.Printf("Start scheduling:%s\n", s.JobHandler)
	//router
	/*	if len(s.jobManager.ServerAddr) == 0 {
			return
		}
		dial, _ := grpc.Dial(s.jobManager.ServerAddr[0])*/
	//add Log
	logDo := &do.JobLogDo{
		JobId:                s.Id,
		ManageId:             s.JobManager.Id,
		DispatchHandler:      s.JobHandler,
		Retry:                s.Retry,
		ExecuteConsumingTime: -1,
		ExecuteStatus:        -1,
	}
	orm.DB.Create(logDo)
	//router
	node := addr[0]
	now := time.Now()
	logDo.DispatchTime = &now
	// 0失败 1成功
	logDo.DispatchStatus = 1
	logDo.DispatchType = 1
	logDo.DispatchAddress = node.Addr
	logDo.ExecuteStatus = 1
	if dispatchService(s, node.Addr, logDo.Id) != nil {
		//调度失败
		logDo.DispatchStatus = 2
		logDo.ExecuteStatus = -1
		orm.DB.Delete(lock)
	}
	orm.DB.Updates(logDo)
	log.Printf("任务调度成功:[%s][%d]", s.JobHandler, s.Id)
}

func dispatchService(s *Scheduler, addr string, logId int64) error {
	dial, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer dial.Close()
	con := dispatch.NewServiceClient(dial)
	//dispatch
	if _, err := con.Call(context.Background(), &dispatch.Request{
		ServiceId:  s.JobHandler,
		Retry:      s.Retry,
		CallbackId: logId,
	}); err != nil {
		return err
	}
	return nil
}
