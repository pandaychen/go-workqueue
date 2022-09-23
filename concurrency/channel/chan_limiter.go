package channel

type ConcurrencyChanLimiter struct {
	limit chan struct{}
	size  int64
}

func NewConcurrencyChanLimiter(concurrency int64) *ConcurrencyChanLimiter {
	var (
		i int64
	)
	if concurrency == 0 {
		concurrency = 1
	}
	l := &ConcurrencyChanLimiter{
		limit: make(chan struct{}, concurrency),
		size:  concurrency,
	}

	for i = 0; i < concurrency; i++ {
		// 初始化令牌
		l.limit <- struct{}{}
	}
	return l
}

func (l *ConcurrencyChanLimiter) Acquire() {
	<-l.limit
	return
}

func (l *ConcurrencyChanLimiter) Release() {
	l.limit <- struct{}{}
}
