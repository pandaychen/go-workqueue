package dispatcher

import (
	"context"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pandaychen/go-workerqueue/broker"
	"github.com/pandaychen/go-workerqueue/broker/common"
	"github.com/pandaychen/go-workerqueue/handler"
	"github.com/pandaychen/go-workerqueue/task"
	"go.uber.org/zap"

	"github.com/pandaychen/goes-wrapper/zaplog"

	"github.com/pandaychen/goes-wrapper/pyerrors"
)

// 核心结构定义
type Dispatcher struct {
	ApiService gin.Engine //提供对外操作API

	swgp sync.WaitGroup

	ctx       context.Context
	ctxCancel context.CancelFunc

	//broker 通用结构
	taskBroker broker.Broker
	brokerType string

	// 注册信息：key=topic/value=handler处理方法
	handlersStore     map[string]*handler.PoolHandler
	handlersStoreLock sync.RWMutex

	// key=topic/value=async task channel
	tasksChan map[string]chan task.TaskElement
	tasksLock sync.RWMutex

	logger *zap.Logger
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
		swgp:              sync.WaitGroup{},
	}

	dis.logger = zaplog.ZapLoggerInit("go-workerqueue", "./queue.log")

	//opt

	return dis
}

// HandlerBindTopic：handler注入，绑定到topic
func (d *Dispatcher) HandlerBindTopic(bindParams *task.HandlerBindParams) error {
	var (
		err error
	)
	d.handlersStoreLock.Lock()
	defer d.handlersStoreLock.Unlock()

	if err = bindParams.Validator(); err != nil {
		d.logger.Error("bindParams.Validator error", zap.Any("errmsg", err))
		return err
	}

	if _, exists := d.handlersStore[bindParams.Topic]; exists {
		d.logger.Error("not found topic", zap.String(bindParams.Topic))
		return nil
	}

	//构造poolHandler
	handler := &handler.PoolHandler{
		Topic:       bindParams.Topic,
		FuncCaller:  bindParams.FuncCall,
		Concurrency: bindParams.Concurrency,
	}

	d.handlersStore[bindParams.Topic] = handler

	return nil
}

/* 注册broker
1. queueType：broker类型
2. args：broker配置
*/
func (d *Dispatcher) RegisteBrokerDriver(queueType string, args ...interface{}) error {
	if d.taskBroker != nil {
		return pyerrors.ErrWorkQueDriverExists
	}
	queueType = strings.ToLower(queueType)
	switch queueType {
	case common.QUEUE_TYPE_LOCAL:
		return pyerrors.ErrWorkQueBadDriver
	case common.QUEUE_TYPE_REDIS:
		return false, pyerrors.ErrWorkQueBadDriverConfig
	default:
		return pyerrors.ErrWorkQueBadDriver
	}

	return nil
}
