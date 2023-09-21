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

type XllJobApp struct {
	port       int
	engine     *gin.Engine
	grpcServer *grpc.Server
}

func NewXllJobApp(port int) *XllJobApp {
	return &XllJobApp{
		port:   port,
		engine: gin.Default(),
	}
}
func (app *XllJobApp) Start() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.port),
		Handler: app.engine, // 使用 Gin 引擎作为处理程序
	}

	//初始化任务
	job := handle.NewXllJobHandle()
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
	job.Stop()
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
	dispatch.RegisterNodeServer(server, handle.NewRegisterHandle())
	dispatch.RegisterJobServer(server, handle.NewJobMonitorHandle())
	go func() {
		err := server.Serve(lis)
		if err != nil {
			log.Fatalf("Grpc Service startup failed:%s", err.Error())
		}
	}()
	return server
}
