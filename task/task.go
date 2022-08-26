package task

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pandaychen/goes-wrapper/generator/id_generator"
)

//  任务结构
type TaskElement struct {
	Ctx  context.Context
	Lock sync.Mutex

	TaskId string `json:"task_id"` //uuid

	opts     *TaskOption `json:"options"` //关联任务配置信息
	Topic    string      `json:"topic"`
	Message  string      `json:"message"`
	Args     interface{} `json:"args"` //任务参数
	TaskType int         `json:"task_type"`
	ExpireAt uint64      `json:"expiration_at"` //任务执行具体时间戳
	//for crontab job
	CrontabSpec string `json:"cron_spec,omitempty"` //crontab任务表达式
}

func NewTaskElement(optApply ...TaskOptionFunc) *TaskElement {
	var (
		taskinfo *TaskElement
	)
	taskinfo = &TaskElement{
		TaskId: id_generator.NextID(),
	}

	taskinfo.opts = NewTaskOption(optApply...)

	return taskinfo
}

func (t *TaskElement) Marshal() (string, error) {
	jstr, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(jstr), nil
}

func main() {
	task := NewTaskElement(SetTopicName("test"))
	fmt.Println(task.opts.Topic)
}
