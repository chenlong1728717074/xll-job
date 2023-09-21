package test

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/robfig/cron/v3"
	"testing"
	"time"
	"xll-job/scheduler/handle"
)

func TestCron(t *testing.T) {
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	c := cron.New(cron.WithParser(parser))

	c.AddFunc("0/2 * * * * ? *", func() {
		fmt.Println("xx")
	})
	c.Start()
	select {}
}
func TestCronexpr(t *testing.T) {
	expr, err := cronexpr.Parse("22 * * * * ? *")
	if err != nil {
		// 处理解析错误
		panic(err)
	}
	nextTime := expr.Next(time.Now())
	fmt.Println(nextTime)
}
func TestChan(t *testing.T) {
	register := handle.NewRegisterHandle()
	register.Start()
	//register.RegisterNodeChan <- core.NewServiceNode("11", 1, "xx")
	select {}
}
