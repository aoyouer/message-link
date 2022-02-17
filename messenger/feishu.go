package messenger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	FeishuMessenger *FeishuClient
)

type FeishuMessage struct {
	Title   string                 `json:"title"`
	Content []FeishuMessageContent `json:"content"`
}

type FeishuMessageContentItem struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
	Href string `json:"href,omitempty"`
}

type FeishuMessageContent []FeishuMessageContentItem

func initFeishu() {
	api := viper.GetViper().GetString("feishu.webhook")
	// 增加是否为空的判断
	FeishuMessenger = &FeishuClient{
		webhook: api,
	}
}

type FeishuClient struct {
	webhook string
}

func GetFeishuMessenger() *FeishuClient {
	return FeishuMessenger
}

func (f *FeishuClient) SendMessage(msg string) {
	contentType := "application/json"
	sendData := fmt.Sprintf(textMsgTemplate, msg)

	res, err := http.Post(f.webhook, contentType, strings.NewReader(sendData))
	if err != nil {
		zap.S().Errorf("failed to send feishu FeishuMessage: %s\n", err.Error())
	} else {
		zap.L().Info("success to send feishu FeishuMessage")
	}
	res.Body.Close()
}

func (f *FeishuClient) SendHyperTextMessage(msg FeishuMessage) {
	contentType := "application/json"
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		zap.L().Error("failed to marshal: " + err.Error())
		return
	}

	res, err := http.Post(f.webhook, contentType, strings.NewReader(fmt.Sprintf(hyperMsgTemplate, string(jsonMsg))))
	if err != nil {
		zap.S().Errorf("failed to send feishu FeishuMessage: %s\n", err.Error())
	} else {
		zap.L().Info("success to send feishu FeishuMessage")
	}
	res.Body.Close()
}

var textMsgTemplate = `{
	"msg_type": "text",
	"content": {"text": "%s" }
}`

var hyperMsgTemplate = `
{
	"msg_type": "post",
	"content": {
	  "post": {
		"zh_cn": %s
	  }
	}
  }
`
