package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xll-job/scheduler/grpc/dispatch"
	"xll-job/scheduler/handle"
	"xll-job/web/router"
)

type XllJob struct {
	port       int
	engine     *gin.Engine
	grpcServer *grpc.Server
}

func NewApp(port int) *XllJob {
	return &XllJob{
		port:   port,
		engine: gin.Default(),
	}
}

func LoadCache() {
	//加载所有的管理器

	//然后再加载所有的任务
}

func (app *XllJob) Start() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.port),
		Handler: app.engine, // 使用 Gin 引擎作为处理程序
	}

	//初始化任务
	job := handle.NewXllJob()
	handle.Xll_Job = job
	job.LoadJob()
	job.Start()
	//grpc监听服务
	app.grpcServer = grpcLis()
	//路由
	router.Router(app.engine)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("SERVER START SUCCESS")
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown err:", err)
	}
	log.Println("Server exiting")
}
func grpcLis() *grpc.Server {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 8082))
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	dispatch.RegisterNodeServer(server, handle.NewRegister())
	go func() {
		server.Serve(lis)
	}()
	return server
}
