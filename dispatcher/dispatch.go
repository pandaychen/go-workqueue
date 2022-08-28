package dispatcher

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pandaychen/go-workerqueue/broker"
	"github.com/pandaychen/go-workerqueue/handler"
	"github.com/pandaychen/go-workerqueue/task"
)

// 核心结构定义
type Dispatcher struct {
	ApiService gin.Engine //提供对外操作API

	waitGroup sync.WaitGroup

	ctx       context.Context
	ctxCancel context.CancelFunc

	//broker
	taskBroker broker.Broker
	brokerType string

	// 注册信息：topic-handler映射
	handlersStore     map[string]*handler.PoolHandler
	handlersStoreLock sync.RWMutex

	tasksChan map[string]chan task.TaskElement
	tasksLock sync.RWMutex
}

func NewDispatcher() *Dispatcher {
	wCtx, wCtxCancel := context.WithCancel(context.Background())
	dis := &Dispatcher{
		ctx:               wCtx,
		ctxCancel:         wCtxCancel,
		handlersStore:     make(map[string]*handler.PoolHandler),
		tasksChan:         make(map[string]chan task.TaskElement),
		handlersStoreLock: sync.RWMutex{},
		tasksLock:         sync.RWMutex{},
		waitGroup:         sync.WaitGroup{},
	}

	//opt

	return dis
}
