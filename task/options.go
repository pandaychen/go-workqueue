package task

type TaskOptionFunc func(*TaskOption)

type TaskOption struct {
	Topic    string      `json:"topic"` //
	Message  string      `json:"message"`
	Args     interface{} `json:"args"` //任务参数
	TaskType int         `json:"task_type"`
	ExpireAt uint64      `json:"expire_at"` //任务执行具体时间戳
	//for crontab job
	CrontabSpec string `json:"cron_spec,omitempty"` //crontab任务表达式
}

func (o *TaskOption) Validator() error {
	return nil
}

func NewTaskOption(optFunctions ...TaskOptionFunc) *TaskOption {
	//init
	option := &TaskOption{}

	for _, o := range optFunctions {
		//apply options
		o(option)
	}
	return option
}

func SetTopicName(topic string) TaskOptionFunc {
	return func(options *TaskOption) {
		options.Topic = topic
	}
}

func SetCrontabSpec(crontab string) TaskOptionFunc {
	return func(options *TaskOption) {
		options.CrontabSpec = crontab
	}
}

func SetMessage(message string) TaskOptionFunc {
	return func(options *TaskOption) {
		options.Message = message
	}
}

func SetTaskType(task_type int) TaskOptionFunc {
	return func(options *TaskOption) {
		options.TaskType = task_type
	}
}

func SetArgs(args interface{}) TaskOptionFunc {
	return func(options *TaskOption) {
		options.Args = args
	}
}

func SetWithExpireAt(expire_at uint64) TaskOptionFunc {
	return func(options *TaskOption) {
		options.ExpireAt = expire_at
	}
}
