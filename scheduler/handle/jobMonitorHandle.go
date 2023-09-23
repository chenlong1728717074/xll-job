package handle

import (
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"xll-job/orm"
	"xll-job/orm/do"
	"xll-job/scheduler/grpc/dispatch"
	"xll-job/scheduler/util"
)

// JobMonitorHandle  这个结构体用于任务失败监听/**/
type JobMonitorHandle struct {
	failJobDone chan struct{}
	timeoutDone chan struct{}

	dispatch.UnimplementedJobServer
}

func NewJobMonitorHandle() *JobMonitorHandle {
	return &JobMonitorHandle{
		failJobDone: make(chan struct{}),
		timeoutDone: make(chan struct{}),
	}
}

func (jobMonitor *JobMonitorHandle) Start() {
	jobMonitor.EnableFailScan()
	jobMonitor.EnableTimeoutScan()
}
func (jobMonitor *JobMonitorHandle) EnableFailScan() {

	go func() {
		log.Println("失败任务处理器已开启失败")
		select {
		case <-jobMonitor.failJobDone:
			log.Println("失败任务处理器已关闭")
			return
		}
	}()
}

func (jobMonitor *JobMonitorHandle) EnableTimeoutScan() {
	go func() {
		log.Println("超时任务处理器已开启")
		select {
		case <-jobMonitor.timeoutDone:
			log.Println("超时任务处理器已关闭")
			return
		}
	}()
}

func (jobMonitor *JobMonitorHandle) Callback(ctx context.Context, resp *dispatch.CallbackResponse) (*emptypb.Empty, error) {
	var jobLog do.JobLogDo
	orm.DB.First(&jobLog, resp.GetId())
	if jobLog.Id == 0 {
		return nil, errors.New("entry does not exist")
	}
	//async handle log
	go jobMonitor.handleLog(&jobLog, resp)
	return &emptypb.Empty{}, nil
}

func (jobMonitor *JobMonitorHandle) handleLog(jobLog *do.JobLogDo, resp *dispatch.CallbackResponse) {
	var job do.JobInfoDo
	orm.DB.First(&job, jobLog.JobId)
	if job.Id == 0 || !job.Enable {
		jobLog.Retry = 0
	}
	jobLog.ExecuteStatus = resp.Status
	startTime := resp.StartTime.AsTime()
	endTime := resp.EndTime.AsTime()
	jobLog.ExecuteStartTime = &startTime
	jobLog.ExecuteEndTime = &endTime
	jobLog.ExecuteConsumingTime = util.ComputingTime(startTime, endTime)
	orm.DB.Updates(&jobLog)
	jobMonitor.handleExecuteLog(resp.GetId(), resp.Logs)
	jobMonitor.Unlock(job.Id)
}

func (jobMonitor *JobMonitorHandle) Unlock(id int64) {
	if id == 0 {
		return
	}
	lock := &do.JobLockDo{
		Id: id,
	}
	orm.DB.Delete(lock)
}

func (jobMonitor *JobMonitorHandle) handleExecuteLog(id int64, logs []string) {
	//日志处理
	var executeLogDo do.ExecutionLog
	orm.DB.Select("id").First(&executeLogDo, id)
	executeLogs := "--------------以下是执行日志--------------\n"
	if len(logs) != 0 {
		for _, s := range logs {
			executeLogs += s + "\n"
		}
	}
	executeLogDo.ExecuteLogs = executeLogs
	if executeLogDo.Id == 0 {
		executeLogDo.Id = id
		orm.DB.Create(&executeLogDo)
	} else {
		orm.DB.Updates(&executeLogDo)
	}
}
