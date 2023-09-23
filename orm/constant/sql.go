package constant

const (
	DeleteLock          = "DELETE FROM tb_job_lock"
	EffectiveTimeoutJob = "SELECT  t2.*,t1.timeout FROM tb_job_info t1 INNER JOIN tb_job_log t2 ON t1.id = t2.job_id WHERE t2.execute_status=1 AND  DATE_ADD(t2.dispatch_time, INTERVAL t1.timeout SECOND)< NOW() AND t1.deleted_time IS NULL"
	LapseTimeoutJob     = "SELECT t1.id from tb_job_log t1 LEFT JOIN (SELECT * FROM tb_job_info WHERE deleted_time IS NULL) t2 ON t1.job_id=t2.id WHERE t2.id is null"
)
