package core

import (
	"google.golang.org/grpc"
)

type JobInfo struct {
	len        int
	jobHandler string
	cron       string
	retry      int
	con        *grpc.ClientConn
}
