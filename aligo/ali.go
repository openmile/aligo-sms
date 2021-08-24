package aligo

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type SendData struct {
	Key      string `json:"key"`
	UserId   string `json:"user_id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Msg      string `json:"msg"`
	MsgType  string `json:"msg_type"` // default "SMS"
	Title    string `json:"title"`
}

type ReceiveData struct {
	ResultCode float64 `json:"result_code"`
	Message    string  `json:"message"`
	MsgId      float64 `json:"msg_id"`
	SuccessCnt float64 `json:"success_cnt"`
	ErrorCnt   float64 `json:"error_cnt"`
	MsgType    string  `json:"msg_type"`
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

	var aligoRes ReceiveData

	client := &http.Client{}
	resp, err := client.PostForm("https://apis.aligo.in/send/", formData)
	if err != nil {
		aligoRes.ResultCode = -1
		aligoRes.Message = "Aligo connection error: 잠시 후에 다시 시도해주세요."
		//log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			aligoRes.ResultCode = -1
			aligoRes.Message = "IO close error: 잠시 후에 다시 시도해주세요."
			//log.Fatal(err)
		}
	}(resp.Body)

	func(Body io.ReadCloser) {
		err := json.NewDecoder(resp.Body).Decode(&aligoRes)
		if err != nil {
			aligoRes.ResultCode = -1
			aligoRes.Message = "JSON decode error: 잠시 후에 다시 시도해주세요."
			//log.Fatal(err)
		}
	}(resp.Body)
	return aligoRes
}
