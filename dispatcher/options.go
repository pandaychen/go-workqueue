package dispatcher

//dispatcher配置

type DispatcherOptions struct {
	//retry：指数退避的等待超时
	retryTimes         int64
	fetcherConcurrency int64 //并发度
	workerConcurrency  int64

	slowLog bool //慢操作记录
}

func (o *DispatcherOptions) Check() {
	if o.fetcherConcurrency <= 0 {
		o.fetcherConcurrency = 1
	}

	if o.workerConcurrency <= 0 {
		o.workerConcurrency = 1
	}
}
