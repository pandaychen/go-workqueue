package handler

import (
	"github.com/pandaychen/goes-wrapper/go-workerqueue/task"
)

// 单个处理方法
type HandlerFunc func(task task.TaskElement) task.TaskRunResult

func (h HandlerFunc) Run(task task.TaskElement) task.TaskRunResult {
	return h(task)
}
