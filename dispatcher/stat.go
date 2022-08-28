package dispatcher

//统计结构

// 当前状态
type DispatcherStat struct {
	PullTaskCount   int64
	SuccHandleCount int64
	FailHandleCount int64
}
