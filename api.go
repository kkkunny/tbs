package tbs

import (
	"encoding/json"
	"errors"
	"github.com/kkkunny/GoMy/crypto"
	"github.com/kkkunny/GoMy/http/requests"
)

// 与服务器之间的交流
type api struct {
	botApi string
	request requests.Request
	lastUpdateHash string  // 上一次更新的哈希
}
// 处理回复
func (this *api)handleResponse(response *requests.Response)error{
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
func (this *api)getUpdates()(*JsonRequest, error){
	var url = this.botApi + "/getUpdates"
	response, err := this.request.Get(url, nil)
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
func (this *api)getUpdate()(*JsonUpdate, error){
	request, err := this.getUpdates()
	if err == nil && len(request.Result) > 0{
		result := request.Result[len(request.Result)-1]
		data, err := json.Marshal(&result)
		if err != nil{
			return &result, err
		}
		hash := string(crypto.EncodeSha1(data))
		if this.lastUpdateHash == ""{
			this.lastUpdateHash = hash
		}else if this.lastUpdateHash != hash{
			this.lastUpdateHash = hash
			return &result, nil
		}else{
			return nil, nil
		}
	}
	return nil, err
}
// 发送信息(上限4096个字符)
func (this *api)SendMessage(id int, msg string)error{
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
func (this *api)SendPhoto(id int, photoUrl string, title string)error{
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
func (this *api)SendVideo(id int, videoUrl string, title string)error{
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
func (this *api)SendMediaGroup(id int, mediaType string, medias []string)error{
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