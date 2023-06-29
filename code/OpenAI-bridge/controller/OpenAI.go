package controller

import (
	"Learning_Record/code/OpenAI-bridge/internal"
	"Learning_Record/code/OpenAI-bridge/models"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
)

type OpenAI struct {
	BaseController
}

// GetModelList openAI开放接口 -- 获取可用模型列表
func (o *OpenAI) GetModelList(c *gin.Context) {
	//获取返回信息
	list, err := models.GetModelList()
	if err != nil {
		ApiReturn(c, internal.CodeExecError, internal.MsgDialerError, err.Error())
		return
	}
	ApiReturn(c, internal.CodeSuccess, internal.MsgSuccess, list)
}

// ChatGpt openAI开放接口 -- model: GPT-3.5 -- chat
func (o *OpenAI) ChatGpt(c *gin.Context) {
	//从body获取对话信息
	var params models.ChatReq
	if err := c.ShouldBind(&params); err != nil {
		internal.PrintLogRus("error", "controller/OpenAI.ChatGpt3", c.Request.Body, err.Error(), "get params err")
		ApiReturn(c, internal.CodeParamsError, internal.MsgParamsError, internal.MsgEmpty)
		return
	}
	//获取对话信息
	if len(params.MessageList) == 0 {
		ApiReturn(c, internal.CodeParamsError, internal.MsgParamsError, "对话输入不能为空")
		return
	}
	//获取返回信息
	resp, err := models.Chat(params)
	if err != nil {
		ApiReturn(c, internal.CodeExecError, internal.MsgDialerError, err.Error())
		return
	}
	ApiReturn(c, internal.CodeSuccess, internal.MsgSuccess, resp)
}

// ChatGptStream openAI开放接口流式传输 -- model: GPT-3.5 -- chat
func (o *OpenAI) ChatGptStream(c *gin.Context) {
	//从body获取对话信息
	var params models.ChatReq
	if err := c.ShouldBind(&params); err != nil {
		internal.PrintLogRus("error", "controller/OpenAI.ChatGpt3Stream", c.Request.Body, err.Error(), "get params err")
		ApiReturn(c, internal.CodeParamsError, internal.MsgParamsError, internal.MsgEmpty)
		return
	}
	//获取对话信息
	if len(params.MessageList) == 0 {
		ApiReturn(c, internal.CodeParamsError, internal.MsgParamsError, "对话输入不能为空")
		return
	}
	//获取返回信息
	stream, err := models.ChatStream(params)
	defer stream.Close()
	if err != nil {
		ApiReturn(c, internal.CodeExecError, internal.MsgDialerError, err.Error())
		return
	}
	//循环读取数据
	var respContent string
	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) { //读到数据流的结尾，退出循环
			break
		}
		if err != nil { //获取流数据报错，则返回错误信息
			internal.PrintLogRus("error", "controller/OpenAI.ChatGpt3Stream", "Stream error: ", err.Error())
			c.SSEvent("message", response) // 返回完整json
			c.Writer.Flush()
			respContent += response.Choices[0].Delta.Content
			return
		} else { //返回流数据信息
			c.SSEvent("message", response) // 返回完整json
			c.Writer.Flush()
			respContent += response.Choices[0].Delta.Content
		}
	}
	internal.PrintLogRus("debug", "controller/OpenAI.ChatGpt3Stream", "get answer is: ", respContent)
}

// DallEGenerateImages openAI开放接口 -- model: DALL·E -- Image generation
func (o *OpenAI) DallEGenerateImages(c *gin.Context) {
	//从body获取相关参数
	var params models.ImageReq
	if err := c.ShouldBind(&params); err != nil {
		internal.PrintLogRus("error", "controller/OpenAI.DallEGenerateImages", c.Request.Body, err.Error(), "get params err")
		ApiReturn(c, internal.CodeParamsError, internal.MsgParamsError, internal.MsgEmpty)
		return
	}
	//关键词信息不能为空
	if params.Prompt == "" {
		ApiReturn(c, internal.CodeParamsError, internal.MsgParamsError, "输入关键词不能为空")
		return
	}
	//获取返回信息
	resp, err := models.CreateImage(params)
	if err != nil {
		ApiReturn(c, internal.CodeExecError, internal.MsgDialerError, err.Error())
		return
	}
	ApiReturn(c, internal.CodeSuccess, internal.MsgSuccess, resp)
}

// Completions openAI开放接口 -- 文本补全 -- 支持ChatGPT3系列模型
func (o *OpenAI) Completions(c *gin.Context) {
	//从body获取对话信息
	var params models.CompReq
	if err := c.ShouldBind(&params); err != nil {
		internal.PrintLogRus("error", "controller/OpenAI.ChatGpt3", c.Request.Body, err.Error(), "get params err")
		ApiReturn(c, internal.CodeParamsError, internal.MsgParamsError, internal.MsgEmpty)
		return
	}
	//关键词信息不能为空
	if params.Prompt == "" {
		ApiReturn(c, internal.CodeParamsError, internal.MsgParamsError, "输入文本不能为空")
		return
	}
	//获取返回信息
	resp, err := models.Completions(params)
	if err != nil {
		ApiReturn(c, internal.CodeExecError, internal.MsgDialerError, err.Error())
		return
	}
	ApiReturn(c, internal.CodeSuccess, internal.MsgSuccess, resp)
}
