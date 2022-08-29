package handler

import "sync"

// 池化的工作方法
type PoolHandler struct {
	sync.RWMutex
	Topic       string
	Concurrency int //并发数
	FuncCaller  HandlerFunc
}

func NewPoolHandler(pool_size int, caller HandlerFunc) *PoolHandler {
	return &PoolHandler{
		Concurrency: pool_size,
		FuncCaller:  caller,
	}
}
