package tbs

import "strings"

// 回复JSON
type JsonResponse struct {
	Ok bool `json:"ok"`
	ErrorCode int `json:"error_code"`
	Description string `json:"description"`
}

// 请求JSON
type JsonRequest struct {
	Ok bool `json:"ok"`
	Result []JsonUpdate `json:"result"`
}
// 消息更新
type JsonUpdate struct {
	UpdateId int64 `json:"update_id"`
	Message JsonMessage `json:"message"`
}
// 获取命令信息
func (this *JsonUpdate)GetCommandInfos()(bool, []string, []string){
	var commands, params []string
	if this.Message.Text != "" && this.Message.Text[0] == '/'{
		infos := strings.Split(this.Message.Text, " ")
		if len(infos) > 0{
			for i:=1; i<len(infos); i++{
				params = append(params, infos[1])
			}
			commands := strings.Split(infos[0][1:], "/")
			return true, commands, params
		}
	}
	return false, commands, params
}
// 消息
type JsonMessage struct {
	MessageId int `json:"message_id"`
	From JsonUser `json:"from"`
	Chat JsonChat `json:"chat"`
	Date int64 `json:"date"`
	Text string `json:"text"`
	Entities []JsonMessageEntity `json:"entities"`  // 消息实体
}
// 用户
type JsonUser struct {
	Id int `json:"id"`
	IsBot bool `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username string `json:"username"`
	LanguageCode string `json:"language_code"`
}
// 对话
type JsonChat struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
	Username string `json:"username"`
	Type string `json:"type"`
}
// 实体
type JsonMessageEntity struct {
	Type string `json:"type"`  // 类型
	Offset int `json:"offset"`  // 开始位置
	Length int `json:"length"`  // 字符长度
}

// 上传媒体
type InputMedia struct {
	Type string `json:"type"`
	Media string `json:"media"`
}