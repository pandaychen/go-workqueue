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
	FuncCall    HandlerFunc
	Topic       string `json:"topic" validate:"required"`
	Concurrency int    `json:"concurrency" validate:"required"`
}

// 校验参数合法性
func (p *HandlerBindParams) Validator() error {
	return g_validator.Struct(p)
}
