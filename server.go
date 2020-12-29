package tbs

import (
	"errors"
	"fmt"
	"github.com/kkkunny/GoMy/http/requests"
	"github.com/kkkunny/GoMy/log"
	"os"
	"strings"
	"sync"
	"time"
)

// 日志管理器
var Logger *log.Logger

// 新建一个电报机器人服务器
func NewTBServer(bootApi string, proxy string) *TBServer {
	server := &TBServer{
		mutex:       sync.RWMutex{},
		api:         &Api{botApi: bootApi, request: *requests.NewRequestWithProxy(proxy)},
		proxy:       proxy,
		models:      make(map[string]*Model),
		closeModels: make(map[string]string),
		commands:    make(map[string]*command),
	}
	server.builtInCommand()
	return server
}

// 电报机器人服务器
type TBServer struct {
	mutex       sync.RWMutex        // 读写锁
	api         *Api                // 机器人api
	proxy       string              // 代理
	models      map[string]*Model   // 模块map
	closeModels map[string]string   // 关闭的模块
	commands    map[string]*command // 命令
}

// 添加模块
func (this *TBServer) AddModel(model *Model) {
	if _, ok := this.models[model.name]; !ok {
		this.models[model.name] = model
		return
	}
	panic(errors.New("Model is existed : " + model.name))
}

// 新增模块
func (this *TBServer) NewModel(name string, introd string, updateHandler Handler) *Model {
	if _, ok := this.models[name]; !ok {
		model := NewModel(name, introd, updateHandler)
		this.models[name] = model
		return model
	}
	panic(errors.New("Model is existed : " + name))
}

// 关闭模块
func (this *TBServer) CloseModel(name string) error {
	if _, ok := this.closeModels[name]; ok {
		return errors.New("this Model is closing : " + name)
	}
	if _, ok := this.models[name]; ok {
		this.closeModels[name] = ""
		return nil
	}
	return errors.New("don't exist this Model : " + name)
}

// 打开模块
func (this *TBServer) OpenModel(name string) error {
	if _, ok := this.models[name]; !ok {
		return errors.New("don't exist this Model : " + name)
	}
	if _, ok := this.closeModels[name]; ok {
		delete(this.closeModels, name)
		return nil
	}
	return errors.New("this Model is running : " + name)
}

// 获取运行中的模块
func (this *TBServer) getRunningModels() map[string]*Model {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	var result = make(map[string]*Model)
	for name, model := range this.models {
		if _, ok := this.closeModels[name]; !ok {
			result[name] = model
		}
	}
	return result
}

// 获取命令
func (this *TBServer) getCommand(names ...string) *command {
	if len(names) == 1 {
		if command, ok := this.commands[names[0]]; ok {
			return command
		}
	} else if len(names) > 1 {
		for _, model := range this.models {
			if command := model.getCommand(names[1]); command != nil {
				return command
			}
		}
	}
	return nil
}

// 增加命令
func (this *TBServer) AddCommand(name string, introd string, handler Handler, ps ...string) {
	if this.getCommand(name) != nil {
		panic(errors.New("this command exist : " + name))
	}
	var paramsMap = make(params)
	for _, p := range ps {
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

// 新消息接收运行函数
func (this *TBServer) runNewUpdate() {
	for {
		ud, err := (*this.api).getUpdate()
		if err == nil && ud != nil {
			// 命令处理
			if ok, commands, ps := ud.GetCommandInfos(); ok {
				if com := this.getCommand(commands...); com != nil && com.loadParams(ps...) == nil {
					if err = Logger.WriteInfoLog(fmt.Sprintf("来自%s的命令：%s\t参数：%s", ud.Message.From.Username, "/"+strings.Join(commands, "/"), strings.Join(ps, ","))); err != nil {
						_ = Logger.WriteError(err)
					}
					handle := func() {
						session := &Session{Api: *this.api, Update: *ud, Params: com.params}
						(*com).run(session)
					}
					go handle()
				}
			} else { // 转发给各个模块
				if err = Logger.WriteInfoLog(fmt.Sprintf("来自%s的消息：%s", ud.Message.From.Username, ud.Message.Text)); err != nil {
					_ = Logger.WriteError(err)
				}
				for _, m := range this.getRunningModels() {
					handle := func(m2 Model) {
						session := &Session{Api: *this.api, Update: *ud}
						m2.run(session)
					}
					go handle(*m)
				}
			}
		} else if err != nil {
			_ = Logger.WriteError(err)
		}
		time.Sleep(1 * time.Second)
	}
}

// 任务运行函数
func (this *TBServer) runTask() {
	for _, m := range this.models {
		for _, t := range m.tasks {
			handle := func(t2 task) {
				for {
					sess := &TaskSession{*this.api}
					t2.Run(sess)
					time.Sleep(t2.interval)
				}
			}
			go handle(*t)
		}
	}
}

// 运行
func (this *TBServer) Run() {
	_ = Logger.WriteLog("TelegramBotServer is running......\n")
	wait := sync.WaitGroup{}
	wait.Add(2)
	// 新消息接收
	go func() {
		defer wait.Done()
		this.runNewUpdate()
	}()
	// 任务
	go func() {
		defer wait.Done()
		this.runTask()
	}()
	wait.Wait()
}

func init() {
	Logger = log.New(os.Stdout)
}
