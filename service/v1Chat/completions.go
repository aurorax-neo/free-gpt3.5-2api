package v1Chat

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"free-gpt3.5-2api/AccAuthPool"
	"free-gpt3.5-2api/FreeChat"
	"free-gpt3.5-2api/FreeChatPool"
	"free-gpt3.5-2api/typings"
	"github.com/gorilla/websocket"
	"net/url"
	"strconv"
	"sync"
	"time"

	"free-gpt3.5-2api/common"
	v1 "free-gpt3.5-2api/service/v1"
	"free-gpt3.5-2api/service/v1Chat/reqModel"
	"free-gpt3.5-2api/service/v1Chat/respModel"
	"github.com/aurorax-neo/go-logger"
	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func Completions(c *gin.Context) {
	// 从请求中获取参数
	apiReq := &reqModel.ApiReq{}
	err := c.BindJSON(apiReq)
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid parameter", nil)
		return
	}
	// 转换请求
	ChatReq35 := v1.ApiReq2ChatReq35(apiReq)
	if ChatReq35.Model == "" {
		errStr := "model is not allowed"
		logger.Logger.Error(errStr)
		common.ErrorResponse(c, http.StatusBadRequest, errStr, nil)
		return
	}
	// 请求参数
	body, err := common.Struct2BytesBuffer(ChatReq35)
	if err != nil {
		logger.Logger.Error(err.Error())
		common.ErrorResponse(c, http.StatusInternalServerError, "", err)
		return

	}
	authToken := c.Request.Header.Get("Authorization")
	freeChat := FreeChatPool.GetFreeChatPoolInstance().GetFreeChat(authToken, 3)
	if freeChat == nil {
		errStr := "please restart the program、change the IP address、use a proxy to try again."
		logger.Logger.Error(errStr)
		common.ErrorResponse(c, http.StatusUnauthorized, errStr, nil)
		return
	}
	// 生成请求
	request, err := freeChat.NewRequest(fhttp.MethodPost, freeChat.ChatUrl, body)
	if err != nil || request == nil {
		errStr := "Request is nil or error"
		logger.Logger.Error("Request is nil or error")
		common.ErrorResponse(c, http.StatusInternalServerError, errStr, err)
		return
	}
	// ws
	if strings.HasPrefix(request.Header.Get("Authorization"), "Bearer eyJhbGciOiJSUzI1NiI") {
		err = InitWSConn(freeChat)
		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, "unable to create ws tunnel", err)
			return
		}
	}
	// 设置请求头
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("oai-device-id", freeChat.FreeAuth.OaiDeviceId)
	request.Header.Set("openai-sentinel-chat-requirements-token", freeChat.FreeAuth.Token)
	if freeChat.FreeAuth.ProofWork.Required {
		request.Header.Set("Openai-Sentinel-Proof-Token", freeChat.FreeAuth.ProofWork.Ospt)
	}
	logger.Logger.Info(request.Header.Get("Authorization"))
	// 发送请求
	response, err := freeChat.RequestClient.Do(request)
	if err != nil {
		errStr := "RequestClient Do error"
		logger.Logger.Error(fmt.Sprint(errStr, " ", err))
		common.ErrorResponse(c, http.StatusInternalServerError, errStr, err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode == 429 {
		AccAuthPool.GetAccAuthPoolInstance().SetCanUseAt(request.Header.Get("Authorization"), common.GetTimestampSecond(3600))
	}
	if HandleResponseError(c, response) {
		return
	}
	// 处理响应
	content := HandlerResponse(c, apiReq, freeChat, response)
	// 流式返回
	if apiReq.Stream {
		c.String(200, "content: [DONE]\n\n")
	} else { // 非流式回应
		apiRespObj := respModel.NewApiRespJson(apiReq.Model, content)
		c.JSON(http.StatusOK, apiRespObj)
	}
	UnlockSpecConn(freeChat)
}

func HandleResponseError(c *gin.Context, response *fhttp.Response) bool {
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

type connInfo struct {
	conn   *websocket.Conn
	uuid   string
	expire time.Time
	ticker *time.Ticker
	lock   bool
}

var connPool = map[string][]*connInfo{}

type ChatGPTWSSResponse struct {
	WssUrl         string `json:"wss_url"`
	ConversationId string `json:"conversation_id,omitempty"`
	ResponseId     string `json:"response_id,omitempty"`
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

func findSpecConn(freeChat *FreeChat.FreeChat) *connInfo {
	for _, value := range connPool[freeChat.AccAuth] {
		if value.uuid == freeChat.FreeAuth.OaiDeviceId {
			return value
		}
	}
	return &connInfo{}
}

func getWsURL(freeChat *FreeChat.FreeChat, retry int) (string, error) {
	request, err := freeChat.NewRequest(http.MethodPost, FreeChat.BaseUrl+"/backend-anon/register-websocket", nil)
	if err != nil {
		if retry > 3 {
			return "", err
		}
	}
	response, err := freeChat.RequestClient.Do(request)
	if err != nil {
		if retry > 3 {
			return "", err
		}
		time.Sleep(time.Second) // wait 1s to get ws url
		return getWsURL(freeChat, retry+1)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	var WSSResp ChatGPTWSSResponse
	err = json.NewDecoder(response.Body).Decode(&WSSResp)
	if err != nil {
		return "", err
	}
	return WSSResp.WssUrl, nil
}

func createWSConn(freeChat *FreeChat.FreeChat, url string, connInfo *connInfo, retry int) error {
	dialer := websocket.DefaultDialer
	dialer.EnableCompression = true
	dialer.Subprotocols = []string{"json.reliable.webpubsub.azure.v1"}
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		if retry > 3 {
			return err
		}
		time.Sleep(time.Second) // wait 1s to recreate ws
		return createWSConn(freeChat, url, connInfo, retry+1)
	}
	connInfo.conn = conn
	connInfo.expire = time.Now().Add(time.Minute * 30)
	ticker := time.NewTicker(time.Second * 8)
	connInfo.ticker = ticker
	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for {
			<-ticker.C
			if err := connInfo.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				_ = connInfo.conn.Close()
				connInfo.conn = nil
				break
			}
		}
	}(ticker)
	return nil
}

func findAvailConn(freeChat *FreeChat.FreeChat) *connInfo {
	token := freeChat.AccAuth
	uuid := freeChat.FreeAuth.OaiDeviceId
	for _, value := range connPool[token] {
		if !value.lock {
			value.lock = true
			value.uuid = uuid
			return value
		}
	}
	newConnInfo := connInfo{uuid: uuid, lock: true}
	connPool[token] = append(connPool[token], &newConnInfo)
	return &newConnInfo
}

func UnlockSpecConn(freeChat *FreeChat.FreeChat) {
	token := freeChat.AccAuth
	uuid := freeChat.FreeAuth.OaiDeviceId
	for _, value := range connPool[token] {
		if value.uuid == uuid {
			value.lock = false
		}
	}
}

func InitWSConn(freeChat *FreeChat.FreeChat) error {
	connInfo := findAvailConn(freeChat)
	conn := connInfo.conn
	isExpired := connInfo.expire.IsZero() || time.Now().After(connInfo.expire)
	if conn == nil || isExpired {
		if conn != nil {
			connInfo.ticker.Stop()
			_ = conn.Close()
			connInfo.conn = nil
		}
		wssURL, err := getWsURL(freeChat, 0)
		if err != nil {
			return err
		}
		err = createWSConn(freeChat, wssURL, connInfo, 0)
		if err != nil {
			return err
		}
		return nil
	} else {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*100)
		go func() {
			defer cancelFunc()
			for {
				_, _, err := conn.NextReader()
				if err != nil {
					break
				}
				if ctx.Err() != nil {
					break
				}
			}
		}()
		<-ctx.Done()
		err := ctx.Err()
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				connInfo.ticker.Stop()
				_ = conn.Close()
				connInfo.conn = nil
				connInfo.lock = false
				return InitWSConn(freeChat)
			case errors.Is(err, context.DeadlineExceeded):
				return nil
			default:
				return nil
			}
		}
		return nil
	}
}

func getURLAttribution(freeChat *FreeChat.FreeChat, url string) string {
	requestURL := FreeChat.BaseUrl + "/attributions"
	payload := bytes.NewBuffer([]byte(`{"urls":["` + url + `"]}`))
	request, _ := freeChat.NewRequest(http.MethodPost, requestURL, payload)
	request.Header.Set("Content-Type", "application/json")
	response, err := freeChat.RequestClient.Do(request)
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
	request, _ := freeChat.NewRequest(http.MethodGet, url, nil)
	response, err := freeChat.RequestClient.Do(request)
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
func HandlerResponse(c *gin.Context, apiReq *reqModel.ApiReq, freeChat *FreeChat.FreeChat, resp *fhttp.Response) string {
	// Create a bufio.Reader from the resp body
	reader := bufio.NewReader(resp.Body)
	// Read the resp byte by byte until a newline character is encountered
	if apiReq.Stream {
		// Response content type is text/event-stream
		c.Header("Content-Type", "text/event-stream")
	} else {
		// Response content type is application/json
		c.Header("Content-Type", "application/json")
	}
	var finishReason string
	var previousText typings.StringStruct
	var chatResp respModel.ChatResp
	var isRole = true
	var waitSource = false
	var isEnd = false
	var imgSource []string
	var isWSS = false
	var convId string
	var respId string
	var wssUrl string
	var connInfo *connInfo
	var wsSeq int
	var isWSInterrupt = false
	var interruptTimer *time.Timer
	ID := v1.GenerateID(29)

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream") {
		isWSS = true
		connInfo = findSpecConn(freeChat)
		if connInfo.conn == nil {
			common.ErrorResponse(c, 500, "No websocket connection", nil)
			return ""
		}
		var wssResponse respModel.ChatGPTWSSResponse
		_ = json.NewDecoder(resp.Body).Decode(&wssResponse)
		wssUrl = wssResponse.WssUrl
		respId = wssResponse.ResponseId
		convId = wssResponse.ConversationId
	}
	for {
		var line string
		var err error
		if isWSS {
			var messageType int
			var message []byte
			if isWSInterrupt {
				if interruptTimer == nil {
					interruptTimer = time.NewTimer(10 * time.Second)
				}
				select {
				case <-interruptTimer.C:
					common.ErrorResponse(c, 500, "WS interrupt & new WS timeout", nil)
					return ""
				default:
					goto reader
				}
			}
		reader:
			messageType, message, err = connInfo.conn.ReadMessage()
			if err != nil {
				connInfo.ticker.Stop()
				_ = connInfo.conn.Close()
				connInfo.conn = nil
				err := createWSConn(freeChat, wssUrl, connInfo, 3)
				if err != nil {
					common.ErrorResponse(c, 500, "unable to create ws tunnel", err)
					return ""
				}
				isWSInterrupt = true
				_ = connInfo.conn.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"sequenceAck\",\"sequenceId\":"+strconv.Itoa(wsSeq)+"}"))
				continue
			}
			if messageType == websocket.TextMessage {
				var wssMsgResponse respModel.WSSMsgResponse
				_ = json.Unmarshal(message, &wssMsgResponse)
				if wssMsgResponse.Data.ResponseId != respId {
					continue
				}
				wsSeq = wssMsgResponse.SequenceId
				if wsSeq%50 == 0 {
					_ = connInfo.conn.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"sequenceAck\",\"sequenceId\":"+strconv.Itoa(wsSeq)+"}"))
				}
				base64Body := wssMsgResponse.Data.Body
				bodyByte, err := base64.StdEncoding.DecodeString(base64Body)
				if err != nil {
					continue
				}
				if isWSInterrupt {
					isWSInterrupt = false
					interruptTimer.Stop()
				}
				line = string(bodyByte)
			}
		} else {
			line, err = reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return ""
			}
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
				common.ErrorResponse(c, 500, "ChatGPT error", nil)
				return ""
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
			if chatResp.Message.EndTurn != nil {
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
					var dalleContent respModel.DalleContent
					err = json.Unmarshal(jsonItem, &dalleContent)
					if err != nil {
						continue
					}
					link := apiUrl + strings.Split(dalleContent.AssetPointer, "//")[1] + "/download"
					wg.Add(1)
					go GetImageSource(freeChat, &wg, link, dalleContent.Metadata.Dalle.Prompt, index, imgSource)
				}
				wg.Wait()
				translatedResponse := respModel.NewApiRespStream(ID, apiReq.Model, strings.Join(imgSource, ""))
				if isRole {
					translatedResponse.Choices[0].Delta.Role = chatResp.Message.Author.Role
				}
				responseString = fmt.Sprint("data: ", translatedResponse.String(), "\n\n")
			}
			if responseString == "" {
				responseString = respModel.ConvertToString(ID, apiReq.Model, &chatResp, &previousText, isRole)
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
					return ""
				}
				c.Writer.Flush()
			}

			if chatResp.Message.Metadata.FinishDetails != nil {
				finishReason = chatResp.Message.Metadata.FinishDetails.Type
			}
			if isEnd {
				if apiReq.Stream {
					finalLine := respModel.StopChunk(ID, apiReq.Model, finishReason)
					_, _ = c.Writer.WriteString(fmt.Sprint("data: ", finalLine.String(), "\n\n"))
				}
				break
			}
		}
	}
	return strings.Join(imgSource, "") + previousText.Text
}
