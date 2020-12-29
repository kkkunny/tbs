package tbs

import (
	"strconv"
)

// 命令处理函数
type Handler func(sess *Session)
// 任务处理函数
type TaskHandler func(sess *TaskSession)

// 参数
type params map[string]string
// 获取参数
func (this params)GetParam(key string)string{
	if _, ok := this[key]; ok{
		return this[key]
	}
	return ""
}
// 获取int类型参数
func (this params)GetParamInt(key string)int{
	value := this.GetParam(key)
	if value != ""{
		result, err := strconv.Atoi(value)
		if err == nil{
			return result
		}
	}
	return 0
}

// 会话
type Session struct {
	Api
	Update JsonUpdate // 消息更新
	Params params      // 参数
}

// 任务会话
type TaskSession struct {
	Api
}