package core

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
	"xll-job/orm"
	"xll-job/orm/constant"
	"xll-job/orm/do"
	"xll-job/scheduler/grpc/dispatch"
)

// TriggerExecute call on triggering
func TriggerExecute(s *Scheduler) {
	//I think scheduling should not be done without a service when triggered, so there is no need to save logs
	addrs := s.JobManager.ServerAddr
	if len(addrs) == 0 {
		return
	}
	//service lock;Prevent parallel processing of tasks
	lock := jobLock(s.Id)
	if lock == nil {
		return
	}
	//add Log
	logDo := &do.JobLogDo{
		JobId:                s.Id,
		ManageId:             s.JobManager.Id,
		DispatchHandler:      s.JobHandler,
		Retry:                0,
		ProcessingStatus:     constant.NoProcessingRequired,
		ExecuteConsumingTime: -1,
	}
	orm.DB.Create(logDo)
	//router
	now := time.Now()
	logDo.DispatchTime = &now
	// 0失败 1成功
	logDo.DispatchStatus = 1
	logDo.DispatchType = 1
	logDo.ExecuteStatus = 1
	addr := addrs[0].Addr
	logDo.DispatchAddress = addr
	dispatchService(logDo, lock)

}

func dispatchService(logDo *do.JobLogDo, lock *do.JobLockDo) {
	if err := callService(logDo.DispatchHandler, logDo.DispatchAddress, logDo.Id); err != nil {
		//调度失败
		logDo.DispatchStatus = 2
		logDo.ExecuteStatus = -1
		logDo.Remark = err.Error()
		orm.DB.Delete(lock)
	}
	orm.DB.Updates(logDo)
}

func jobLock(id int64) *do.JobLockDo {
	lock := &do.JobLockDo{
		Id: id,
	}
	tx := orm.DB.Create(lock)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil
	}
	return lock
}

// RetryExecute call on retry
func RetryExecute(jobLog *do.JobLogDo) {
	lock := jobLock(jobLog.JobId)
	if lock == nil {
		return
	}
	now := time.Now()
	jobLog.DispatchTime = &now
	jobLog.DispatchStatus = 1
	jobLog.ExecuteStatus = 1
	jobLog.ExecuteStartTime = nil
	jobLog.ExecuteEndTime = nil
	jobLog.ExecuteConsumingTime = -1
	jobLog.ExecuteStatus = 1
	jobLog.Retry = jobLog.Retry + 1
	dispatchService(jobLog, lock)
}

func callService(serviceId string, addr string, logId int64) error {
	dial, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer dial.Close()
	con := dispatch.NewServiceClient(dial)
	//dispatch
	_, callErr := con.Call(context.Background(), &dispatch.Request{
		ServiceId:  serviceId,
		CallbackId: logId,
	})
	return callErr
}
