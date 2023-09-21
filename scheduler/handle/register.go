package handle

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
	"xll-job/orm"
	"xll-job/orm/do"
	"xll-job/scheduler/core"
	"xll-job/scheduler/grpc/dispatch"
)

type Register struct {
	flushDone        chan struct{}
	registerDone     chan struct{}
	registerNodeChan chan *core.ServiceNode
	dispatch.UnimplementedNodeServer
}

func NewRegister() *Register {
	return &Register{
		flushDone:        make(chan struct{}),
		registerDone:     make(chan struct{}),
		registerNodeChan: make(chan *core.ServiceNode, 1000),
	}
}
func (register *Register) Start() {
	register.inspectServer()
	register.registerServer()
}

func (register *Register) Register(ctx context.Context, req *dispatch.RegisterRequest) (*emptypb.Empty, error) {
	var m do.JobManagementDo
	orm.DB.First(&m, req.GetJobManagerId())
	if m.Id == 0 {
		return nil, errors.New("JobManagement NOT FOUND")
	}
	//register.registerNodeChan <- core.NewServiceNode(req.ServiceAddr, m.Id, m.AppName)
	register.addNode(core.NewServiceNode(req.ServiceAddr, m.Id, m.AppName))
	return &emptypb.Empty{}, nil
}

func (register *Register) Stop(addr string) {
	register.flushDone <- struct{}{}
	register.registerDone <- struct{}{}
}

func (register *Register) registerServer() {
	//添加服务,刷新服务
	go func() {
		log.Printf("已开启服务注册服务")
		for {
			select {
			case <-register.registerDone:
				fmt.Println("注册服务关闭....")
				return
			case node := <-register.registerNodeChan:
				register.addNode(node)
			}
		}
	}()
}
func (register *Register) addNode(node *core.ServiceNode) {
	flag := true
	for index := range ServiceNodeList {
		if ServiceNodeList[index].Addr == node.Addr {
			ServiceNodeList[index].RegisterTime = time.Now()
			flag = false
			break
		}
	}
	if flag {
		ServiceNodeList = append(ServiceNodeList, node)
	}
	log.Printf("已刷新来自%s的注册\n", node.Addr)
}
func (register *Register) inspectServer() {
	go func() {
		//睡十秒,等待服务注册
		time.Sleep(time.Second * 10)
		for {
			select {
			case <-register.flushDone:
				fmt.Println("服务检查关闭服务....")
				return
			default:
				go flushServer()
				time.Sleep(time.Second * 30)
			}
		}
	}()
}

func flushServer() {
	startTime := time.Now().UnixNano() / 1000000
	log.Printf("开始进行服务检查:%d", startTime)

	now := time.Now().Add(-time.Second * 90)
	newServiceNodeList := make([]*core.ServiceNode, 0)
	//获取所有存活的服务
	for _, node := range ServiceNodeList {
		if now.Before(node.RegisterTime) {
			node.RegisterTime = time.Now()
			newServiceNodeList = append(newServiceNodeList, node)
		}
	}
	//刷新缓存中的服务
	ServiceNodeList = newServiceNodeList
	// 分组ServiceNodeList中的节点
	temp := make(map[int64][]*core.ServiceNode)
	for _, node := range ServiceNodeList {
		temp[node.JobManagerId] = append(temp[node.JobManagerId], node)
	}
	//重新分配服务
	for k := range JobManagerMap {
		manager := JobManagerMap[k]
		manager.ServerAddr = temp[k]
	}
	endTime := time.Now().UnixNano() / 1000000
	log.Printf("服务检查结束:%d,耗时%d", endTime, endTime-startTime)
}
