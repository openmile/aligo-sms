package aligo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type SendData struct {
	AKey     string `json:"a_key"` // request (accepted by this server) key name
	Key      string `json:"key"`   // aligo key name
	UserId   string `json:"user_id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Msg      string `json:"msg"`
	MsgType  string `json:"msg_type"` // default "SMS"
	Title    string `json:"title"`
}

type ReceiveData struct {
	ResultCode interface{} `json:"result_code"`
	Message    string      `json:"message"`
	MsgId      string      `json:"msg_id"`
	SuccessCnt float64     `json:"success_cnt"`
	ErrorCnt   float64     `json:"error_cnt"`
	MsgType    string      `json:"msg_type"`
}

func PostAligo(data *SendData) ReceiveData {
	formData := url.Values{}
	formData.Set("key", data.Key)
	formData.Set("user_id", data.UserId)
	formData.Set("sender", data.Sender)
	formData.Set("receiver", data.Receiver)
	formData.Set("msg", data.Msg)
	formData.Set("msg_type", data.MsgType)
	formData.Set("title", data.Title)

	aligoRes := new(ReceiveData)

	client := &http.Client{}
	resp, err := client.PostForm("https://apis.aligo.in/send/", formData)
	if err != nil {
		aligoRes.ResultCode = "-1"
		aligoRes.Message = "Aligo connection error: 잠시 후에 다시 시도해주세요."
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			aligoRes.ResultCode = "-1"
			aligoRes.Message = "IO close error: 잠시 후에 다시 시도해주세요."
		}
	}(resp.Body)

	func(Body io.ReadCloser) {
		err := json.NewDecoder(resp.Body).Decode(&aligoRes)
		if err != nil {
			aligoRes.ResultCode = "-1"
			aligoRes.Message = "JSON decode error: 잠시 후에 다시 시도해주세요."
		}
	}(resp.Body)

	return *aligoRes
}
