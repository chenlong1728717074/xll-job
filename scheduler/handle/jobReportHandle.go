package handle

// JobReportHandle 咱不实现 没意义
type JobReportHandle struct {
	reportJobDone chan struct{}
}

func NewJobReportHandle() *JobReportHandle {
	return &JobReportHandle{
		reportJobDone: make(chan struct{}),
	}
}

func (jobReport *JobReportHandle) Start() {
	jobReport.EnableReportTask()
}

func (jobReport *JobReportHandle) Stop() {

}

func (jobReport *JobReportHandle) EnableReportTask() {

}
