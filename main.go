package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"xll-job/app"
)

/*todo 前期采用robfig的cron,后期将会利用cronexpr实现支持年的定时任务*/
func main() {
	server := app.NewApp(8081)
	server.Start()
}

func test() {
	/*parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	c := cron.New(cron.WithParser(parser))*/
	c := cron.New()
	c.AddFunc("0/2 * * * * ? *", func() {
		fmt.Println("xx")
	})
	c.Start()
	select {}
}
