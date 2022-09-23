package dispatcher

import (
	"context"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pandaychen/go-workerqueue/broker"
	"github.com/pandaychen/go-workerqueue/broker/common"
	channellimit "github.com/pandaychen/go-workerqueue/concurrency/channel"
	"github.com/pandaychen/go-workerqueue/handler"
	"github.com/pandaychen/go-workerqueue/metrics"
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

	//生产者限速
	producersLimiter     map[string]channellimit.ConcurrencyChanLimiter
	producersLimiterLock sync.Mutex

	// key=topic/value=async task channel
	tasksChan map[string]chan task.TaskElement
	tasksLock sync.RWMutex

	logger     *zap.Logger
	metricsCli *metrics.Metrics
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

// initTopicHandlers：初始化任务生产者及限速桶
func (d *Dispatcher) initTopicHandlersProducer(ctx context.Context) {
	d.handlersStoreLock.Lock()
	defer d.handlersStoreLock.Unlock()
	for topicName, handler := range d.handlersStore {
		d.producersLimiterLock.Lock()
		d.producersLimiter[topicName] = channellimit.NewConcurrencyChanLimiter(handler.Concurrency)
		d.producersLimiterLock.Unlock()

		d.tasksLock.Lock()
		d.tasksChan[topicName] = make(chan task.TaskElement, 0)
		d.tasksLock.Unlock()
	}
}

// fetcherTasks：通过令牌桶限制获取任务，并下发到异步handler channel
func (d *Dispatcher) fetcherTasks(ctx context.Context, topic_name string) error {
	var (
		err      error
		taskinfo task.TaskElement
	)

	//get tasks from broker dequeue...
	if taskinfo, err = d.taskBroker.Dequeue(ctx, topic_name); err != nil {
		//
		return err
	}

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
