package models

import (
	"Learning_Record/code/OpenAI-bridge/config"
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
)

// GetModelList 获取可用模型列表
func GetModelList() (list openai.ModelsList, err error) {
	list, err = config.Client.ListModels(context.Background())
	if err != nil {
		if ParseOpenAiApiErr(err) { //如报错信息为429，则切换appKey，并重新请求
			//重新请求接口
			list, err = config.Client.ListModels(context.Background())
		}
	}
	return list, err
}

// Chat 对话接口
func Chat(param ChatReq) (resp openai.ChatCompletionResponse, err error) {
	resp, err = config.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    param.ModelName,
			Messages: param.MessageList,
			User:     param.ChatId,
		},
	)
	if err != nil {
		if ParseOpenAiApiErr(err) { //如报错信息为429，则切换appKey，并重新请求
			resp, err = config.Client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model:    param.ModelName,
					Messages: param.MessageList,
					User:     param.ChatId,
				},
			)
		}
	}
	return resp, err
}

// ChatStream 对话接口-流式传输
func ChatStream(param ChatReq) (stream *openai.ChatCompletionStream, err error) {
	stream, err = config.Client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    param.ModelName,
			Messages: param.MessageList,
			User:     param.ChatId,
		},
	)
	if err != nil {
		stream.Close()
		if ParseOpenAiApiErr(err) { //如报错信息为429，则切换appKey，并重新请求
			stream, err = config.Client.CreateChatCompletionStream(
				context.Background(),
				openai.ChatCompletionRequest{
					Model:    param.ModelName,
					Messages: param.MessageList,
					User:     param.ChatId,
				},
			)
		}
	}
	return stream, err
}

// CreateImage 生成图像接口
func CreateImage(param ImageReq) (resp openai.ImageResponse, err error) {
	//选择图片分辨率信息
	var size string
	if param.Size == "256" {
		size = openai.CreateImageSize256x256
	} else if param.Size == "512" {
		size = openai.CreateImageSize512x512
	} else if param.Size == "1024" || param.Size == "" {
		size = openai.CreateImageSize1024x1024
	} else {
		return resp, errors.New("Image size illegal. ")
	}
	//请求接口
	resp, err = config.Client.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Prompt:         param.Prompt,
			N:              param.N,
			Size:           size,
			ResponseFormat: param.ResponseFormat,
		},
	)
	if err != nil {
		if ParseOpenAiApiErr(err) { //如报错信息为429，则切换appKey，并重新请求
			//重新请求接口
			resp, err = config.Client.CreateImage(
				context.Background(),
				openai.ImageRequest{
					Prompt:         param.Prompt,
					N:              param.N,
					Size:           size,
					ResponseFormat: param.ResponseFormat,
				},
			)
		}
	}
	return resp, err
}

// Completions 文本补全接口
func Completions(param CompReq) (resp openai.CompletionResponse, err error) {
	//获取返回信息
	resp, err = config.Client.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			Model:     param.Model,
			Prompt:    param.Prompt,
			MaxTokens: param.MaxTokens,
		},
	)
	if err != nil {
		if ParseOpenAiApiErr(err) { //如报错信息为429，则切换appKey，并重新请求
			//重新请求接口
			resp, err = config.Client.CreateCompletion(
				context.Background(),
				openai.CompletionRequest{
					Model:     param.Model,
					Prompt:    param.Prompt,
					MaxTokens: param.MaxTokens,
				},
			)
		}
	}
	return resp, err
}

// ParseOpenAiApiErr 解析openAI接口返回错误信息
func ParseOpenAiApiErr(err error) bool {
	e := &openai.APIError{}
	if errors.As(err, &e) {
		switch e.HTTPStatusCode {
		case 429: // rate limiting or engine overload (wait and retry) //账号余额用完，切换appKey重试
			//切换appKey
			config.ChangeAppKey()
			//需要重新请求
			return true
		}
	}
	return false
}
