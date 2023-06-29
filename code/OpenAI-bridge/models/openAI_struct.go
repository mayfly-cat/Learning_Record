package models

import "github.com/sashabaranov/go-openai"

// ChatReq 对话接口请求参数
type ChatReq struct {
	MessageList []openai.ChatCompletionMessage `json:"messageList"` //对话列表--包括role / content
	ChatId      string                         `json:"chatId"`      //每个对话的唯一id
	ModelName   string                         `json:"modelName"`   //选择的模型
}

// ImageReq 图片生成接口请求参数
type ImageReq struct {
	Prompt         string `json:"prompt,omitempty"`          //必填字段
	N              int    `json:"n,omitempty"`               //选填字段，默认为1
	Size           string `json:"size,omitempty"`            //选填字段，默认为1024x1024
	ResponseFormat string `json:"response_format,omitempty"` //选填字段，默认为url
	User           string `json:"user,omitempty"`            //选填字段，无默认值
}

// ImageResp 图片生成接口返回结构
type ImageResp struct {
	Url []string `json:"url"` //url数组
}

// CompReq 文本补全接口请求参数
type CompReq struct {
	Model     string `json:"model"`                //模型名称
	Prompt    any    `json:"prompt,omitempty"`     //输入文本
	MaxTokens int    `json:"max_tokens,omitempty"` //最大返回字数
}
