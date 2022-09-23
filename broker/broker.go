package broker

import (
	"context"
	"strings"

	"go-workqueue/broker/common"
)

// 任务队列通用定义（默认采用json进行序列化）
type Broker interface {
	//入队
	Enqueue(ctx context.Context, key string, message string, args ...interface{}) error

	//批量入队
	BatchEnqueue(ctx context.Context, key string, messages []string, args ...interface{}) error

	//从broker中获取元数据
	Dequeue(ctx context.Context, key string) (string, error)

	//ack消息
	AckMsg(ctx context.Context, key string) (bool, error)

	//关闭
	Close() (err error)

	//返回长度
	Len(ctx context.Context, key string) int
}

func NewBroker(broker_type string) Broker {
	switch strings.ToLower(broker_type) {
	case common.QUEUE_TYPE_REDIS:

	}
}
