package v1

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
)

func ErrorResponse(c *gin.Context, code int, msg interface{}, err error) {
	c.JSON(code, gin.H{
		"detail": struct {
			Code  int         `json:"code"`
			Msg   interface{} `json:"msg"`
			Error interface{} `json:"error"`
		}{
			Code:  code,
			Msg:   msg,
			Error: err,
		},
	})
	return
}

func MappingModel(model string) string {
	var modelMapping = map[string]string{
		"gpt-3.5-turbo":          "text-davinci-002-render-sha",
		"gpt-3.5-turbo-16k":      "text-davinci-002-render-sha",
		"gpt-3.5-turbo-16k-0613": "text-davinci-002-render-sha",
		"gpt-3.5-turbo-0301":     "text-davinci-002-render-sha",
		"gpt-3.5-turbo-0613":     "text-davinci-002-render-sha",
		"gpt-3.5-turbo-1106":     "text-davinci-002-render-sha",
	}
	if model == "" {
		return "text-davinci-002-render-sha"
	}
	if v, ok := modelMapping[model]; ok {
		return v
	}
	return "text-davinci-002-render-sha"
}
func GenerateID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := "chatcmpl-"
	for i := 0; i < length; i++ {
		id += string(charset[rand.Intn(len(charset))])
	}
	return id
}

func Obj2Bytes(obj interface{}) ([]byte, error) {
	// 创建一个jsonIter的Encoder
	configCompatibleWithStandardLibrary := jsoniter.ConfigCompatibleWithStandardLibrary
	// 将结构体转换为JSON文本并保持顺序
	bytes, err := configCompatibleWithStandardLibrary.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
