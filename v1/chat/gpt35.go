package chat

import (
	"encoding/json"
	"fmt"
	"free-gpt3.5-2api/chatgpt"
	v1 "free-gpt3.5-2api/v1"
	"free-gpt3.5-2api/v1/chat/reqmodel"
	"free-gpt3.5-2api/v1/chat/respmodel"
	"github.com/aurorax-neo/go-logger"
	"github.com/gin-gonic/gin"
	rv2 "github.com/go-resty/resty/v2"
	"github.com/launchdarkly/eventsource"
	"net/http"
	"strings"
	"time"
)

func gpt35(c *gin.Context, apiReq *reqmodel.ApiReq) {
	// 获取 chatgpt 实例
	instance := chatgpt.GetGpt35Instance()
	// 转换请求
	ChatReq35 := reqmodel.ApiReq2ChatReq35(apiReq)
	// 发送请求
	resp, err := instance.Client.R().
		SetHeader("oai-device-id", instance.Session.OaiDeviceId).
		SetHeader("openai-sentinel-chat-requirements-token", instance.Session.Token).
		SetBody(ChatReq35).
		SetDoNotParseResponse(true).
		Post(chatgpt.ApiUrl)
	if err != nil || resp.StatusCode() != http.StatusOK {
		v1.ErrorResponse(c, http.StatusInternalServerError, "", err)
		logger.Logger.Error(err.Error())
		return
	}

	// 流式返回
	if apiReq.Stream {
		__CompletionsStream(c, apiReq, resp)
	} else { // 非流式回应
		__CompletionsNoStream(c, apiReq, resp)
	}
}

func __CompletionsStream(c *gin.Context, apiReq *reqmodel.ApiReq, resp *rv2.Response) {
	messageTemp := ""
	decoder := eventsource.NewDecoder(resp.RawResponse.Body)
	defer func(decoder *eventsource.Decoder) {
		_, _ = decoder.Decode()
	}(decoder)
	// 响应id
	id := v1.GenerateID(29)
	for {
		event, err := decoder.Decode()
		if err != nil {
			v1.ErrorResponse(c, http.StatusInternalServerError, "", err)
			logger.Logger.Error(err.Error())
			break
		}
		name := event.Event()
		data := event.Data()
		// 空白数据不处理
		if data == "" {
			continue
		}
		// 结束标志
		if data == "[DONE]" {
			apiRespObj := &respmodel.StreamObj{}
			// id
			apiRespObj.ID = id
			// created
			apiRespObj.Created = time.Now().Unix()
			// object
			apiRespObj.Object = "chat.completion.chunk"
			// choices
			delta := respmodel.StreamDeltaObj{
				Content: "",
			}
			choices := respmodel.StreamChoiceObj{
				Delta:        delta,
				FinishReason: "stop",
			}
			apiRespObj.Choices = append(apiRespObj.Choices, choices)
			// model
			apiRespObj.Model = apiReq.Model
			// 生成响应
			bytes, err := v1.Obj2Bytes(apiRespObj)
			if err != nil {
				logger.Logger.Error(err.Error())
				continue
			}
			// 发送响应
			c.SSEvent(name, fmt.Sprint(" ", string(bytes)))
			// 结束
			c.SSEvent(name, " [DONE]")
			break
		}
		chatResp35 := &respmodel.ChatResp35{}
		err = json.Unmarshal([]byte(data), chatResp35)
		// 仅处理assistant的消息
		if chatResp35.Message.Author.Role == "assistant" {
			// 仅处理第一个part
			parts := chatResp35.Message.Content.Parts[0]
			// 去除重复数据
			content := strings.Replace(parts, messageTemp, "", 1)
			messageTemp = parts
			// 空白数据不处理
			if content == "" {
				continue
			}
			// 生成响应 model
			apiRespObj := &respmodel.StreamObj{}
			// id
			apiRespObj.ID = id
			// created
			apiRespObj.Created = time.Now().Unix()
			// object
			apiRespObj.Object = "chat.completion.chunk"
			// choices
			delta := respmodel.StreamDeltaObj{
				Content: content,
			}
			choices := respmodel.StreamChoiceObj{
				Delta: delta,
			}
			apiRespObj.Choices = append(apiRespObj.Choices, choices)
			// model
			apiRespObj.Model = apiReq.Model
			// 生成响应
			bytes, err := v1.Obj2Bytes(apiRespObj)
			if err != nil {
				logger.Logger.Error(err.Error())
				continue
			}
			// 发送响应
			c.SSEvent(name, fmt.Sprint(" ", string(bytes)))
			// 继续
			continue
		}
	}
}

func __CompletionsNoStream(c *gin.Context, apiReq *reqmodel.ApiReq, resp *rv2.Response) {
	content := ""
	decoder := eventsource.NewDecoder(resp.RawResponse.Body)
	defer func(decoder *eventsource.Decoder) {
		_, _ = decoder.Decode()
	}(decoder)
	for {
		event, err := decoder.Decode()
		if err != nil {
			v1.ErrorResponse(c, http.StatusInternalServerError, "", err)
			logger.Logger.Error(err.Error())
			break
		}
		data := event.Data()
		// 空白数据不处理
		if data == "" {
			continue
		}
		// 结束标志
		if data == "[DONE]" {
			apiRespObj := &respmodel.JsonObj{}
			// id
			apiRespObj.ID = v1.GenerateID(29)
			// created
			apiRespObj.Created = time.Now().Unix()
			// object
			apiRespObj.Object = "chat.completion"
			// model
			apiRespObj.Model = apiReq.Model
			// usage
			usage := respmodel.JsonUsageObj{
				PromptTokens:     0,
				CompletionTokens: 0,
				TotalTokens:      0,
			}
			apiRespObj.Usage = usage
			// choices
			message := respmodel.JsonMessageObj{
				Role:    "assistant",
				Content: content,
			}
			choice := respmodel.JsonChoiceObj{
				Message:      message,
				FinishReason: "stop",
				Index:        0,
			}
			apiRespObj.Choices = append(apiRespObj.Choices, choice)
			// 返回响应
			c.JSON(http.StatusOK, apiRespObj)
			break
		}
		chatResp35 := &respmodel.ChatResp35{}
		err = json.Unmarshal([]byte(data), chatResp35)
		// 仅处理assistant的消息
		if chatResp35.Message.Author.Role == "assistant" {
			// 仅处理第一个part
			content = chatResp35.Message.Content.Parts[0]
			// 空白数据不处理
			if content == "" {
				continue
			}
			continue
		}
	}
}
