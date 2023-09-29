package test

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/gorhill/cronexpr"
	"github.com/robfig/cron/v3"
	"os"
	"runtime"
	"testing"
	"time"
	"xll-job/orm"
	"xll-job/orm/do"
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
	orm.DB.Model(&do.JobInfoDo{}).Where("id = ?", 1).Update("is_enable", 0)
}
func TestPassword(t *testing.T) {
	Options := &password.Options{16, 100, 32, sha512.New}
	salt, encode := password.Encode("czcl123de", Options)
	verify := password.Verify("czcl123de", salt, encode, Options)
	fmt.Println(verify)
}

func TestJwt(t *testing.T) {

}

func TestRunTime(t *testing.T) {
	fmt.Println(runtime.GOOS)   // 输出操作系统架构
	fmt.Println(runtime.GOARCH) // 输出操作系统名称
	fmt.Println(runtime.GOMAXPROCS(1))
	fmt.Println(runtime.Version())
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Println(os.Kill.String())
}

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/utsname.h>

	void getOSInfo(char* buf, size_t bufLen) {
	    struct utsname unameData;
	    if (uname(&unameData) == 0) {
	        snprintf(buf, bufLen, "System name: %s\nNode name: %s\nRelease: %s\nVersion: %s\nMachine: %s\n",
	            unameData.sysname, unameData.nodename, unameData.release, unameData.version, unameData.machine);
	    } else {
	        snprintf(buf, bufLen, "Failed to get operating system information.");
	    }
	}
*/
