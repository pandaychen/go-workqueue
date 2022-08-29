package handler

import (
	"github.com/pandaychen/goes-wrapper/go-workerqueue/task"
)

// 单个处理方法
type HandlerFunc func(task task.TaskElement) task.TaskRunResult

func (h HandlerFunc) Run(task task.TaskElement) task.TaskRunResult {
	return h(task)
}

//注册参数
type HandlerBindParams struct {
	funcCall    HandlerFunc
	topic       string
	concurrency int
}
