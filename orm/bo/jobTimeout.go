package bo

import "xll-job/orm/do"

type JobTimeoutBo struct {
	do.JobLogDo
	Timeout int32
}
