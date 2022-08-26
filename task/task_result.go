package task

// 任务执行结果封装

type TaskRunResult struct {
	TaskId  string //回调ID
	RetCode int
	Result  string //
}
