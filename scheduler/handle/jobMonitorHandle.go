package handle

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"xll-job/orm"
	"xll-job/orm/do"
	"xll-job/scheduler/grpc/dispatch"
)

// JobMonitorHandle  这个结构体用于任务失败监听/**/
type JobMonitorHandle struct {
	retryDone   chan struct{}
	timeoutDone chan struct{}
	dispatch.UnimplementedJobServer
}

func NewJobMonitorHandle() *JobMonitorHandle {
	return &JobMonitorHandle{
		retryDone:   make(chan struct{}),
		timeoutDone: make(chan struct{}),
	}
}
func (*JobMonitorHandle) Callback(ctx context.Context, resp *dispatch.CallbackResponse) (*emptypb.Empty, error) {
	fmt.Println("接收到回调消息", resp)
	var jobLog do.JobLogDo
	orm.DB.First(&jobLog, resp.GetId())
	if jobLog.Id == 0 {
		return nil, errors.New("entry does not exist")
	}
	var job do.JobInfoDo
	orm.DB.First(&job, jobLog.JobId)
	if job.Id == 0 || !job.Enable {
		jobLog.ExecuteStatus = 3
		jobLog.Retry = 0
	} else {
		jobLog.ExecuteStatus = resp.Status
		startTime := resp.StartTime.AsTime()
		endTime := resp.EndTime.AsTime()
		jobLog.ExecuteStartTime = &startTime
		jobLog.ExecuteEndTime = &endTime
		consumingTime := (endTime.UnixNano() / 1000000) - (startTime.UnixNano() / 1000000)
		if consumingTime == 0 {
			consumingTime = 1
		}
		jobLog.ExecuteConsumingTime = consumingTime
		logs := ""
		for _, s := range resp.Logs {
			logs += s + "\n"
		}
		jobLog.ExecuteLogs = logs
		lock := do.JobLockDo{
			Id: job.Id,
		}
		orm.DB.Delete(lock)
	}
	orm.DB.Updates(&jobLog)
	return &emptypb.Empty{}, nil
}

func (jobMonitor *JobMonitorHandle) Start() {
	jobMonitor.EnableFailScan()
	jobMonitor.EnableTimeoutScan()
}
func (jobMonitor *JobMonitorHandle) EnableFailScan() {

	go func() {
		log.Println("开启失败任务处理")
		select {}
	}()
}

func (jobMonitor *JobMonitorHandle) EnableTimeoutScan() {
	go func() {
		log.Println("开启超时任务处理")
		select {}
	}()
}
