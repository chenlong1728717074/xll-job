package handle

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"sync"
	"time"
	"xll-job/orm"
	"xll-job/orm/do"
	"xll-job/scheduler/core"
	"xll-job/scheduler/grpc/dispatch"
)

// RegisterHandle 服务注册
type RegisterHandle struct {
	flushDone        chan struct{}
	registerDone     chan struct{}
	registerNodeChan chan *core.ServiceNode
	lock             sync.Mutex
	dispatch.UnimplementedNodeServer
}

func NewRegisterHandle() *RegisterHandle {
	return &RegisterHandle{
		flushDone:        make(chan struct{}),
		registerDone:     make(chan struct{}),
		registerNodeChan: make(chan *core.ServiceNode, 1000),
	}
}
func (register *RegisterHandle) Start() {
	register.inspectServer()
	register.registerServer()
}

func (register *RegisterHandle) Register(ctx context.Context, req *dispatch.RegisterRequest) (*emptypb.Empty, error) {
	var m do.JobManagementDo
	orm.DB.First(&m, req.GetJobManagerId())
	if m.Id == 0 {
		return nil, errors.New("JobManagement NOT FOUND")
	}
	register.registerNodeChan <- core.NewServiceNode(req.ServiceAddr, m.Id, m.AppName)
	return &emptypb.Empty{}, nil
}

func (register *RegisterHandle) Stop() {
	register.flushDone <- struct{}{}
	register.registerDone <- struct{}{}
}

func (register *RegisterHandle) registerServer() {
	//添加服务,刷新服务
	go func() {
		log.Printf("服务注册处理器已开启")
		for {
			select {
			case <-register.registerDone:
				fmt.Println("注服务注册处理器已关闭....")
				return
			case node := <-register.registerNodeChan:
				register.addNode(node)

			}
		}
	}()
}
func (register *RegisterHandle) addNode(node *core.ServiceNode) {
	flag := true
	for index := range ServiceNodeList {
		if ServiceNodeList[index].Addr == node.Addr {
			ServiceNodeList[index].RegisterTime = time.Now()
			flag = false
			break
		}
	}
	if flag {
		register.lock.Lock()
		ServiceNodeList = append(ServiceNodeList, node)
		register.lock.Unlock()
	}
	log.Printf("registration from[%s]has been refreshed\n", node.Addr)
}
func (register *RegisterHandle) inspectServer() {
	go func() {
		log.Println("服务检查处理器已开启")
		//睡十秒,等待服务注册
		time.Sleep(time.Second * 10)
		for {
			select {
			case <-register.flushDone:
				log.Println("服务检查处理器已关闭....")
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
	log.Printf("start scrubbing service node[刷新服务]:%d\n", startTime)
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
	fmt.Println(ServiceNodeList)
	log.Printf("service node refresh completed[刷新完成]:%d,time consuming:%d", endTime, endTime-startTime)
}
