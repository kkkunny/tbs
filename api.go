package tbs

import (
	"errors"
	"github.com/kkkunny/GoMy/http/requests"
)

// 与服务器之间的交流
type Api struct {
	botApi string
	request requests.Request
	lastUpdateId int64  // 上一次更新的id
}
// 处理回复
func (this *Api)handleResponse(response *requests.Response)error{
	var result JsonResponse
	if err := response.Json(&result); err != nil{
		return err
	}
	if !result.Ok{
		return errors.New(result.Description)
	}
	return nil
}
// 获取请求
func (this *Api)getUpdates(offset int, limit int, timeout int, allowedUpdates []string)(*JsonRequest, error){
	var url = this.botApi + "/getUpdates"
	data := map[string]interface{}{
		"offset": offset,
		"limit": limit,
		"timeout": timeout,
		"allowed_updates": allowedUpdates,
	}
	response, err := this.request.Post(url, data, true)
	if err != nil{
		return nil, err
	}
	var result JsonRequest
	err = response.Json(&result)
	if err != nil{
		return &result, err
	}
	return &result, nil
}
// 获取新消息更新
func (this *Api)getUpdate()(*JsonUpdate, error){
	request, err := this.getUpdates(-1, 1, 0, []string{})
	if err == nil && len(request.Result) > 0{
		result := request.Result[len(request.Result)-1]
		if this.lastUpdateId == 0{
			this.lastUpdateId = result.UpdateId
		}else if this.lastUpdateId != result.UpdateId{
			this.lastUpdateId = result.UpdateId
			return &result, nil
		}else{
			return nil, nil
		}
	}
	return nil, err
}
// 发送信息(上限4096个字符)
func (this *Api)SendMessage(id int, msg string)error{
	var url = this.botApi + "/sendMessage"
	data := map[string]interface{}{
		"chat_id": id,
		"text": msg,
	}
	response, err := this.request.Post(url, data, false)
	if err != nil{
		return err
	}
	if err := this.handleResponse(response); err != nil{
		return err
	}
	return nil
}
// 发送图片
func (this *Api)SendPhoto(id int, photoUrl string, title string)error{
	var url = this.botApi + "/sendPhoto"
	data := map[string]interface{}{
		"chat_id": id,
		"photo": photoUrl,
		"caption": title,
	}
	response, err := this.request.Post(url, data, false)
	if err != nil{
		return err
	}
	if err := this.handleResponse(response); err != nil{
		return err
	}
	return nil
}
// 发送视频（只支持mp4）
func (this *Api)SendVideo(id int, videoUrl string, title string)error{
	var url = this.botApi + "/sendPhoto"
	data := map[string]interface{}{
		"chat_id": id,
		"video": videoUrl,
		"caption": title,
	}
	response, err := this.request.Post(url, data, false)
	if err != nil{
		return err
	}
	if err := this.handleResponse(response); err != nil{
		return err
	}
	return nil
}
// 发送一组媒体[2, 9](photo, video)
func (this *Api)SendMediaGroup(id int, mediaType string, medias []string)error{
	var url = this.botApi + "/sendMediaGroup"
	var inputMedias []InputMedia
	for k, mediaUrl := range medias{
		inputMedia := InputMedia{Type: mediaType, Media: mediaUrl}
		inputMedias = append(inputMedias, inputMedia)
		if k >= 8{
			break
		}
	}
	data := map[string]interface{}{
		"chat_id": id,
		"media": inputMedias,
	}
	response, err := this.request.Post(url, data, true)
	if err != nil{
		return err
	}
	if err := this.handleResponse(response); err != nil{
		return err
	}
	return nil
}