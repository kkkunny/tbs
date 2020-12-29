package tbs

import (
	"errors"
	"fmt"
)

// 命令
type command struct {
	name       string   // 名字 eg./help
	Introd     string   // 简介
	paramsList []string // 参数名称列表
	params     params   // 参数
	handler    Handler  // 处理函数
}

// 运行处理函数
func (this *command) run(sess *Session) {
	if this.handler != nil {
		this.handler(sess)
	}
}

// 载入参数
func (this *command) loadParams(ps ...string) error {
	if len(ps) == len(this.paramsList) {
		for i := 0; i < len(ps); i++ {
			this.params[this.paramsList[i]] = ps[i]
		}
		return nil
	}
	return errors.New("params is error")
}

// 内置命令
func (this *TBServer) builtInCommand() {
	// 帮助
	help := func(sess *Session) {
		// 获取命令
		var msg = "{内置命令}"
		for name, c := range this.commands {
			msg += fmt.Sprintf("\n    [%s]%s", "/"+name, c.Introd)
		}
		for mname, m := range this.getRunningModels() {
			msg += fmt.Sprintf("\n{%s}", "/"+mname)
			for cname, c := range m.commands {
				msg += fmt.Sprintf("\n    [%s]%s", "/"+cname, c.Introd)
			}
		}
		if err := sess.SendMessage(sess.Update.Message.From.Id, msg); err != nil {
			_ = Logger.WriteError(err)
		}
	}
	this.AddCommand("help", "获取帮助", help)
	// 查看模块
	model := func(sess *Session) {
		// 获取模块
		var msg = "模块："
		for name, m := range this.getRunningModels() {
			msg += fmt.Sprintf("\n  [%s]%s", name, m.introd)
		}
		if err := sess.SendMessage(sess.Update.Message.From.Id, msg); err != nil {
			_ = Logger.WriteError(err)
		}
	}
	this.AddCommand("Model", "查看模块", model)
}
