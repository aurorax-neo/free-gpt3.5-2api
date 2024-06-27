package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"free-gpt3.5-2api/AccessTokenPool"
	"free-gpt3.5-2api/FreeChat"
	"free-gpt3.5-2api/common"
	"free-gpt3.5-2api/constant"
	"free-gpt3.5-2api/types"
	"github.com/aurorax-neo/tls_client_httpi"
	"github.com/donnie4w/go-logger/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func Completions(c *gin.Context) {
	// 从请求中获取参数
	apiReq := &types.ApiReq{}
	err := c.BindJSON(apiReq)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid parameter", nil)
		return
	}
	// 转换请求
	ChatReq35 := types.ApiReq2ChatReq35(apiReq)
	if ChatReq35.Model == "" {
		errStr := fmt.Sprint("Model is unsupported")
		logger.Error(errStr)
		common.ErrorResponse(c, http.StatusBadRequest, errStr, nil)
		return
	}
	// 请求参数
	body, err := common.Struct2BytesBuffer(ChatReq35)
	if err != nil {
		logger.Error(err.Error())
		common.ErrorResponse(c, http.StatusInternalServerError, "", err)
		return

	}
	token := c.Request.Header.Get("Authorization")
	freeChat, err := FreeChat.GetFreeChat(token, constant.ReTry)
	if err != nil {
		logger.Error(err)
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	if freeChat == nil {
		logger.Error(err)
		common.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}
	headers, cookies := freeChat.GetHC(freeChat.ChatUrl)
	// 设置请求头
	headers.Set(strings.ToLower("Accept-Encoding"), "gzip, deflate, br")
	headers.Set(strings.ToLower("Accept"), "text/event-stream")
	headers.Set(strings.ToLower("Content-Type"), "application/json")
	headers.Set(strings.ToLower("oai-device-id"), freeChat.FreeAuth.OaiDeviceId)
	headers.Set(strings.ToLower("openai-sentinel-chat-requirements-token"), freeChat.FreeAuth.Token)
	if freeChat.FreeAuth.ProofWork.Required {
		headers.Set(strings.ToLower("Openai-Sentinel-Proof-Token"), freeChat.FreeAuth.ProofWork.Ospt)
	}
	// 发送请求
	response, err := freeChat.Http.Request(tls_client_httpi.POST, freeChat.ChatUrl, headers, cookies, body)
	if err != nil {
		errStr := "Http Do error"
		logger.Error(fmt.Sprint(errStr, " ", err))
		common.ErrorResponse(c, http.StatusInternalServerError, errStr, err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode == 429 {
		AccessTokenPool.GetAccAuthPoolInstance().SetCanUseAt(headers.Get("Token"), common.GetTimestampSecond(3600))
	}
	if HandleResponseError(c, response) {
		return
	}
	// 处理响应
	content, err := HandlerResponse(c, apiReq, freeChat, response)
	if err != nil {
		logger.Error(err)
		common.ErrorResponse(c, http.StatusInternalServerError, "", err.Error())
		return
	}
	// 流式返回
	if apiReq.Stream {
		c.String(200, "content: [DONE]\n\n")
	} else { // 非流式回应
		apiRespObj := types.NewApiRespJson(apiReq.Model, content)
		c.JSON(http.StatusOK, apiRespObj)
	}
}

func HandleResponseError(c *gin.Context, response *http.Response) bool {
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var errorResponse map[string]interface{}
		err := json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			// Read response body
			body, _ := io.ReadAll(response.Body)
			common.ErrorResponse(c, response.StatusCode, "Unknown error", errors.New(string(body)))
			return true
		}
		common.ErrorResponse(c, response.StatusCode, errorResponse["detail"], nil)
		return true
	}
	return false
}

var urlAttrMap = make(map[string]string)

type urlAttr struct {
	Url         string `json:"url"`
	Attribution string `json:"attribution"`
}

type fileInfo struct {
	DownloadURL string `json:"download_url"`
	Status      string `json:"status"`
}

func getURLAttribution(freeChat *FreeChat.FreeChat, url string) string {
	requestURL := FreeChat.BaseUrl + "/attributions"
	payload := bytes.NewBuffer([]byte(`{"urls":["` + url + `"]}`))
	headers, cookies := freeChat.GetHC(requestURL)
	headers.Set("Content-Type", "application/json")
	response, err := freeChat.Http.Request(tls_client_httpi.POST, requestURL, headers, cookies, payload)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	var attr urlAttr
	err = json.NewDecoder(response.Body).Decode(&attr)
	if err != nil {
		return ""
	}
	return attr.Attribution
}

func GetImageSource(freeChat *FreeChat.FreeChat, wg *sync.WaitGroup, url string, prompt string, idx int, imgSource []string) {
	defer wg.Done()
	// 获取请求头和cookies
	headers, cookies := freeChat.GetHC(url)
	// 发送请求
	response, err := freeChat.Http.Request(tls_client_httpi.GET, url, headers, cookies, nil)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	var fInfo fileInfo
	err = json.NewDecoder(response.Body).Decode(&fInfo)
	if err != nil || fInfo.Status != "success" {
		return
	}
	imgSource[idx] = "[![image](" + fInfo.DownloadURL + " \"" + prompt + "\")](" + fInfo.DownloadURL + ")"
}
func HandlerResponse(c *gin.Context, apiReq *types.ApiReq, freeChat *FreeChat.FreeChat, resp *http.Response) (string, error) {
	// Create a bufio.Reader from the resp body
	reader := bufio.NewReader(resp.Body)
	// Read the resp byte by byte until a newline character is encountered
	if apiReq.Stream {
		// Response content types is text/event-stream
		c.Header("Content-Type", "text/event-stream")
	} else {
		// Response content types is application/json
		c.Header("Content-Type", "application/json")
	}
	var finishReason string
	var previousText types.StringStruct
	var chatResp types.ChatResp
	var isRole = true
	var waitSource = false
	var isEnd = false
	var imgSource []string
	var convId string
	ID := types.GenerateID(29)

	for {
		var line string
		var err error
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if len(line) < 6 {
			continue
		}
		// Remove "data: " from the beginning of the line
		line = line[6:]
		// Check if line starts with [DONE]
		if !strings.HasPrefix(line, "[DONE]") {
			// Parse the line as JSON
			err = json.Unmarshal([]byte(line), &chatResp)
			if err != nil {
				continue
			}
			if chatResp.Error != nil {
				return "", fmt.Errorf("ChatGPT error: %v", chatResp.Error)
			}
			if chatResp.ConversationId != convId {
				if convId == "" {
					convId = chatResp.ConversationId
				} else {
					continue
				}
			}
			if !(chatResp.Message.Author.Role == "assistant" || (chatResp.Message.Author.Role == "tool" && chatResp.Message.Content.ContentType != "text")) || chatResp.Message.Content.Parts == nil {
				continue
			}
			if chatResp.Message.Metadata.MessageType == "" {
				continue
			}
			if chatResp.Message.Metadata.MessageType != "next" && chatResp.Message.Metadata.MessageType != "continue" || !strings.HasSuffix(chatResp.Message.Content.ContentType, "text") {
				continue
			}
			if chatResp.Message.EndTurn == true {
				if waitSource {
					waitSource = false
				}
				isEnd = true
			}
			if len(chatResp.Message.Metadata.Citations) != 0 {
				r := []rune(chatResp.Message.Content.Parts[0].(string))
				if waitSource {
					if string(r[len(r)-1:]) == "】" {
						waitSource = false
					} else {
						continue
					}
				}
				offset := 0
				for _, citation := range chatResp.Message.Metadata.Citations {
					rl := len(r)
					attr := urlAttrMap[citation.Metadata.URL]
					if attr == "" {
						u, _ := url.Parse(citation.Metadata.URL)
						BaseURL := u.Scheme + "://" + u.Host + "/"
						attr = getURLAttribution(freeChat, BaseURL)
						if attr != "" {
							urlAttrMap[citation.Metadata.URL] = attr
						}
					}
					chatResp.Message.Content.Parts[0] = string(r[:citation.StartIx+offset]) + " ([" + attr + "](" + citation.Metadata.URL + " \"" + citation.Metadata.Title + "\"))" + string(r[citation.EndIx+offset:])
					r = []rune(chatResp.Message.Content.Parts[0].(string))
					offset += len(r) - rl
				}
			} else if waitSource {
				continue
			}
			responseString := ""
			if chatResp.Message.Recipient != "all" {
				continue
			}
			if chatResp.Message.Content.ContentType == "multimodal_text" {
				apiUrl := FreeChat.BaseUrl + "/files/"
				imgSource = make([]string, len(chatResp.Message.Content.Parts))
				var wg sync.WaitGroup
				for index, part := range chatResp.Message.Content.Parts {
					jsonItem, _ := json.Marshal(part)
					var dalleContent types.DalleContent
					err = json.Unmarshal(jsonItem, &dalleContent)
					if err != nil {
						continue
					}
					link := apiUrl + strings.Split(dalleContent.AssetPointer, "//")[1] + "/download"
					wg.Add(1)
					go GetImageSource(freeChat, &wg, link, dalleContent.Metadata.Dalle.Prompt, index, imgSource)
				}
				wg.Wait()
				translatedResponse := types.NewApiRespStream(ID, apiReq.Model, strings.Join(imgSource, ""))
				if isRole {
					translatedResponse.Choices[0].Delta.Role = chatResp.Message.Author.Role
				}
				responseString = fmt.Sprint("data: ", translatedResponse.String(), "\n\n")
			}
			if responseString == "" {
				responseString = types.ConvertToString(ID, apiReq.Model, &chatResp, &previousText, isRole)
			}
			if responseString == "" {
				if isEnd {
					goto endProcess
				} else {
					continue
				}
			}
			if responseString == "【" {
				waitSource = true
				continue
			}
		endProcess:
			isRole = false
			if apiReq.Stream {
				_, err = c.Writer.WriteString(responseString)
				if err != nil {
					return "", err
				}
				c.Writer.Flush()
			}

			if chatResp.Message.Metadata.FinishDetails != nil {
				finishReason = chatResp.Message.Metadata.FinishDetails.Type
			}
			if isEnd {
				if apiReq.Stream {
					finalLine := types.StopChunk(ID, apiReq.Model, finishReason)
					_, _ = c.Writer.WriteString(fmt.Sprint("data: ", finalLine.String(), "\n\n"))
				}
				break
			}
		}
	}
	return strings.Join(imgSource, "") + previousText.Text, nil
}
