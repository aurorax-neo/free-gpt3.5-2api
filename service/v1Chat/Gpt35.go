package v1Chat

import (
	"encoding/json"
	"fmt"
	"free-gpt3.5-2api/FreeGpt35"
	"free-gpt3.5-2api/FreeGpt35Pool"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/service/v1"
	"free-gpt3.5-2api/service/v1Chat/reqmodel"
	"free-gpt3.5-2api/service/v1Chat/respmodel"
	"github.com/aurorax-neo/go-logger"
	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"github.com/launchdarkly/eventsource"
	"io"
	"net/http"
	"strings"
	"time"
)

func Gpt35Completions(c *gin.Context, apiReq *reqmodel.ApiReq) {
	// 获取 FreeGpt35 实例
	ChatGpt35 := FreeGpt35Pool.GetGpt35PoolInstance().GetGpt35(3)
	if ChatGpt35 == nil {
		errStr := "please restart the program、change the IP address、use a proxy to try again."
		logger.Logger.Error(errStr)
		common.ErrorResponse(c, http.StatusUnauthorized, errStr, nil)
		return
	}
	// 转换请求
	ChatReq35 := reqmodel.ApiReq2ChatReq35(apiReq)
	// 请求参数
	body, err := common.Struct2BytesBuffer(ChatReq35)
	if err != nil {
		logger.Logger.Error(err.Error())
		common.ErrorResponse(c, http.StatusInternalServerError, "", err)
		return

	}
	// 生成请求
	request, err := ChatGpt35.NewRequest(fhttp.MethodPost, FreeGpt35.ChatUrl, body)
	if err != nil || request == nil {
		errStr := "Request is nil or error"
		logger.Logger.Error("Request is nil or error")
		common.ErrorResponse(c, http.StatusInternalServerError, errStr, err)
		return
	}
	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("oai-device-id", ChatGpt35.FreeAuth.OaiDeviceId)
	request.Header.Set("openai-sentinel-chat-requirements-token", ChatGpt35.FreeAuth.Token)
	if ChatGpt35.FreeAuth.ProofWork.Required {
		request.Header.Set("Openai-Sentinel-Proof-Token", ChatGpt35.FreeAuth.ProofWork.Ospt)
	}
	// 发送请求
	response, err := ChatGpt35.RequestClient.Do(request)
	if err != nil {
		errStr := "RequestClient Do error"
		logger.Logger.Error(fmt.Sprint(errStr, " ", err))
		common.ErrorResponse(c, http.StatusInternalServerError, errStr, err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		errStr := "Request error"
		logger.Logger.Error(fmt.Sprint(errStr, " ", response.StatusCode))
		common.ErrorResponse(c, response.StatusCode, errStr, nil)
		return
	}
	// 流式返回
	if apiReq.Stream {
		__CompletionsStream(c, apiReq, response)
	} else { // 非流式回应
		__CompletionsNoStream(c, apiReq, response)
	}
}

func __CompletionsStream(c *gin.Context, apiReq *reqmodel.ApiReq, resp *fhttp.Response) {
	messageTemp := ""
	decoder := eventsource.NewDecoder(resp.Body)
	defer func(decoder *eventsource.Decoder) {
		_, _ = decoder.Decode()
	}(decoder)
	// 响应id
	id := v1.GenerateID(29)
	handlingSigns := false
	for {
		event, err := decoder.Decode()
		if err != nil {
			logger.Logger.Error(err.Error())
			common.ErrorResponse(c, http.StatusInternalServerError, "", err)
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
			bytes, err := common.Struct2Bytes(apiRespObj)
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
		if chatResp35.Error != nil && !handlingSigns {
			logger.Logger.Error(fmt.Sprint(chatResp35.Error))
			common.ErrorResponse(c, http.StatusInternalServerError, "", chatResp35.Error)
			return
		}
		// 脏数据不处理
		if err != nil {
			continue
		}
		// 仅处理assistant的消息
		if chatResp35.Message.Author.Role == "assistant" && (chatResp35.Message.Status == "in_progress" || handlingSigns) {
			// handlingSigns 置为 true
			handlingSigns = true
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
			bytes, err := common.Struct2Bytes(apiRespObj)
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

func __CompletionsNoStream(c *gin.Context, apiReq *reqmodel.ApiReq, resp *fhttp.Response) {
	content := ""
	decoder := eventsource.NewDecoder(resp.Body)
	defer func(decoder *eventsource.Decoder) {
		_, _ = decoder.Decode()
	}(decoder)
	handlingSigns := false
	for {
		event, err := decoder.Decode()
		if err != nil {
			logger.Logger.Error(err.Error())
			common.ErrorResponse(c, http.StatusInternalServerError, "", err)
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
		if chatResp35.Error != nil && !handlingSigns {
			logger.Logger.Error(fmt.Sprint(chatResp35.Error))
			common.ErrorResponse(c, http.StatusInternalServerError, "", chatResp35.Error)
			return
		}
		// 脏数据不处理
		if err != nil {
			continue
		}
		// 仅处理assistant的消息
		if chatResp35.Message.Author.Role == "assistant" && (chatResp35.Message.Status == "in_progress" || handlingSigns) {
			// handlingSigns 置为 true
			handlingSigns = true
			// 如果不包含上一次的数据则不处理
			if !strings.Contains(chatResp35.Message.Content.Parts[0], content) {
				continue
			}
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
