package tbs

import "time"

// 任务
type task struct {
	interval time.Duration  // 间隔时间
	handler TaskHandler  // 任务函数
}
// 运行
func (this *task)Run(sess *TaskSession){
	if this.handler != nil{
		this.handler(sess)
	}
}