package tbs

import (
	"errors"
	"time"
)

// 新建一个模块
func NewModel(name string, introd string, updateHandler Handler)*Model {
	return &Model{
		name: name,
		introd: introd,
		update: updateHandler,
		commands: make(map[string]*command),
	}
}
// 模块
type Model struct{
	name string                  // 名字
	introd string  // 简介
	update Handler               // 消息更新处理函数
	commands map[string]*command // 命令
	tasks []*task  // 任务
}
// 运行
func (this *Model)run(sess *Session){
	if this.update != nil{
		this.update(sess)
	}
}
// 获取命令
func (this *Model)getCommand(names string)*command {
	if command, ok := this.commands[names]; ok{
		return command
	}
	return nil
}
// 增加命令
func (this *Model)AddCommand(name string, introd string, handler Handler, ps ...string){
	if this.getCommand(name) != nil{
		panic(errors.New("this command exist : " + name))
	}
	var paramsMap = make(params)
	for _, p := range ps{
		paramsMap[p] = ""
	}
	this.commands[name] = &command{
		name:       name,
		Introd:     introd,
		paramsList: ps,
		params:     paramsMap,
		handler:    handler,
	}
}
// 添加任务
func (this *Model)AddTask(interval time.Duration, handler TaskHandler){
	t := &task{
		interval: interval,
		handler: handler,
	}
	this.tasks = append(this.tasks, t)
}