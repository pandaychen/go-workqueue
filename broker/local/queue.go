package local

//使用Map模拟本地任务队列

import (
	"sync"
)

type LocalTaskQueue struct {
	sync.RWMutex
	Queue map[string][]string //topic==>task队列
}
