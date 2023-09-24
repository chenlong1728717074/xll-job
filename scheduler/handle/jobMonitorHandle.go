package handle

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
	"xll-job/orm"
	"xll-job/orm/bo"
	"xll-job/orm/constant"
	"xll-job/orm/do"
	"xll-job/scheduler/core"
	"xll-job/scheduler/grpc/dispatch"
	"xll-job/utils"
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

func (jobMonitor *JobMonitorHandle) Stop() {
	jobMonitor.failJobDone <- struct{}{}
	jobMonitor.timeoutDone <- struct{}{}
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
	if jobLog.ExecuteStatus != constant.Timeout {
		jobLog.ExecuteStatus = resp.Status
	}
	jobLog.ExecuteStatus = resp.Status
	startTime := resp.StartTime.AsTime()
	endTime := resp.EndTime.AsTime()
	jobLog.ExecuteStartTime = &startTime
	jobLog.ExecuteEndTime = &endTime
	jobLog.ExecuteConsumingTime = utils.ComputingTime(startTime, endTime)
	//The number of retries has reached the predetermined value
	if job.Retry-jobLog.Retry == 0 {
		jobLog.ProcessingStatus = constant.Processed
	}
	orm.DB.Updates(&jobLog)
	jobMonitor.handleExecuteLog(resp.GetId(), resp.Logs)
	jobMonitor.Unlock(job.Id)
}

func (jobMonitor *JobMonitorHandle) EnableFailScan() {
	go func() {
		log.Println("任务重试处理器已开启")
		//wait job call
		time.Sleep(time.Second)
		for {
			select {
			case <-jobMonitor.failJobDone:
				log.Println("失败重试处理器已关闭")
				return
			default:
				jobMonitor.retryJob()
				time.Sleep(time.Second * 10)
			}

		}
	}()
}

func (jobMonitor *JobMonitorHandle) retryJob() {
	var jobLogs []bo.RetryJobBo
	orm.DB.Raw(constant.RetryJob).Scan(&jobLogs)
	if len(jobLogs) == 0 {
		return
	}
	lapseJob := make([]int64, 0)
	retryJobs := make([]bo.RetryJobBo, 0)
	for _, jobLog := range jobLogs {
		if jobLog.Enable {
			retryJobs = append(retryJobs, jobLog)
			continue
		}
		lapseJob = append(lapseJob, jobLog.JobId)
	}
	if len(lapseJob) != 0 {
		orm.DB.Model(&do.JobLogDo{}).
			Where("id in (?)", lapseJob).Update("processing_status", constant.Processed)
	}
	for _, retryJob := range retryJobs {
		var joblog do.JobLogDo
		copier.Copy(&joblog, &retryJob)
		go core.RetryExecute(&joblog)
	}
}

func (jobMonitor *JobMonitorHandle) EnableTimeoutScan() {
	go func() {
		log.Println("超时任务处理器已开启")
		for {
			select {
			case <-jobMonitor.timeoutDone:
				log.Println("超时任务处理器已关闭")
				return
			default:
				jobMonitor.timeoutScan()
				time.Sleep(time.Minute)
			}
		}
	}()
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

// (now-dispatchTime)>10 min ->timeout
func (jobMonitor *JobMonitorHandle) timeoutScan() {
	//先扫描失效任务
	jobMonitor.lapseTimeoutJobScan()
	jobMonitor.effectiveTimeoutJobScan()
}

func (jobMonitor *JobMonitorHandle) lapseTimeoutJobScan() {
	var jobLogs []*do.JobLogDo
	orm.DB.Raw(constant.LapseTimeoutJob).Scan(&jobLogs)
	if len(jobLogs) == 0 {
		return
	}
	for _, jobLog := range jobLogs {
		//Tasks that do not exist do not need to be retried
		jobLog.ExecuteStatus = constant.Timeout
		jobLog.ProcessingStatus = constant.Processed
		orm.DB.Updates(&jobLog)
	}
}

func (jobMonitor *JobMonitorHandle) effectiveTimeoutJobScan() {
	var jobLogs []bo.JobTimeoutBo
	orm.DB.Raw(constant.EffectiveTimeoutJob).Scan(&jobLogs)
	if len(jobLogs) == 0 {
		return
	}
	now := time.Now()
	for _, jobLog := range jobLogs {
		dispatchTime := jobLog.DispatchTime
		orm.DB.Model(&do.JobLogDo{}).Where("id=?", jobLog.Id).Updates(map[string]interface{}{
			"execute_start_time":     dispatchTime,
			"execute_end_time":       &now,
			"execute_consuming_time": utils.ComputingTime(*dispatchTime, now),
			"execute_status":         constant.Timeout,
		})
	}
}
