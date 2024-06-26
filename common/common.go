package common

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func ErrorResponse(c *gin.Context, code int, msg interface{}, err interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
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

// GetTimestampSecond 获取当前时间戳 + 指定 秒
func GetTimestampSecond(second int) int64 {
	return time.Now().Add(time.Second * time.Duration(second)).Unix()
}

func ParseUrl(link string) *url.URL {
	if link == "" {
		return &url.URL{}
	}
	u, err := url.Parse(link)
	if err != nil {
		return &url.URL{}
	}
	return u
}

func GetOrigin(link string) string {
	u := ParseUrl(link)
	if u == nil {
		return ""
	}
	return u.Scheme + "://" + u.Host
}

func Struct2BytesBuffer(v interface{}) (*bytes.Buffer, error) {
	data := new(bytes.Buffer)
	err := json.NewEncoder(data).Encode(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Struct2Bytes(v interface{}) ([]byte, error) {
	// 创建一个jsonIter的Encoder
	configCompatibleWithStandardLibrary := jsoniter.ConfigCompatibleWithStandardLibrary
	// 将结构体转换为JSON文本并保持顺序
	bytes_, err := configCompatibleWithStandardLibrary.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes_, nil
}

func SplitAndAddPre(pre string, str string, sep string) []string {
	var authTokenList []string
	for _, v := range strings.Split(str, sep) {
		authTokenList = append(authTokenList, pre+v)
	}
	return authTokenList
}

func GetRand() rand.Rand {
	// 初始化随机数生成器
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	return *rng
}

func RandomLanguage() string {
	// 初始化随机数生成器
	rng := GetRand()
	// 语言列表
	languages := []string{"af", "am", "ar-sa", "as", "az-Latn", "be", "bg", "bn-BD", "bn-IN", "bs", "ca", "ca-ES-valencia", "cs", "cy", "da", "de", "de-de", "el", "en-GB", "en-US", "es", "es-ES", "es-US", "es-MX", "et", "eu", "fa", "fi", "fil-Latn", "fr", "fr-FR", "fr-CA", "ga", "gd-Latn", "gl", "gu", "ha-Latn", "he", "hi", "hr", "hu", "hy", "id", "ig-Latn", "is", "it", "it-it", "ja", "ka", "kk", "km", "kn", "ko", "kok", "ku-Arab", "ky-Cyrl", "lb", "lt", "lv", "mi-Latn", "mk", "ml", "mn-Cyrl", "mr", "ms", "mt", "nb", "ne", "nl", "nl-BE", "nn", "nso", "or", "pa", "pa-Arab", "pl", "prs-Arab", "pt-BR", "pt-PT", "qut-Latn", "quz", "ro", "ru", "rw", "sd-Arab", "si", "sk", "sl", "sq", "sr-Cyrl-BA", "sr-Cyrl-RS", "sr-Latn-RS", "sv", "sw", "ta", "te", "tg-Cyrl", "th", "ti", "tk-Latn", "tn", "tr", "tt-Cyrl", "ug-Arab", "uk", "ur", "uz-Latn", "vi", "wo", "xh", "yo-Latn", "zh-Hans", "zh-Hant", "zu"}
	// 随机选择一个语言
	randomIndex := rng.Intn(len(languages))
	return languages[randomIndex]
}

// GetAbsPathAndGenerate 获取绝对路径并生成文件或文件夹
func GetAbsPathAndGenerate(path string, isFilePath bool, content string) string {
	// 获取绝对路径
	path = GetAbsPath(path)
	if isFilePath {
		//	判断文件是否存在
		if isExist := fileIsExistAndCreat(path, content); isExist {
			return path
		}
	} else {
		//	判断文件夹是否存在
		if isExist := dirIsExistAndMkdir(path, false); isExist {
			return path
		}
	}
	return path
}

// GetAbsPath 获取绝对路径
func GetAbsPath(path string) string {
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return ""
		}
		return absPath
	}
	return path
}

func dirIsExistAndMkdir(dirPath string, isFile bool) bool {
	// 判断路径是否存在
	_, err := os.Stat(dirPath)
	dir := dirPath
	if err != nil {
		if isFile {
			dir = filepath.Dir(dirPath)
		}
		// 创建路径
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return false
		}
	}
	return true
}

func fileIsExistAndCreat(filePath string, content string) bool {
	//判断文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		// 判断文件夹是否存在
		if isExist := dirIsExistAndMkdir(filePath, true); !isExist {
			return false
		}
		// 创建文件
		_, err := os.Create(filePath)
		if err != nil {
			return false
		}
		if content != "" {
			//	写入content
			file, _ := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
			_, _ = file.Write([]byte(content))
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
		}
	}
	return true
}

// AsyncTimingTask 定时任务 参数含函数
func AsyncTimingTask(nanosecond time.Duration, fun func()) {
	go func() {
		timerChan := time.After(nanosecond)
		// 使用for循环阻塞等待定时器的信号
		for {
			// 通过select语句监听定时器通道和其他事件
			select {
			case <-timerChan:
				fun()
				// 重新设置定时器，以便下一次执行
				timerChan = time.After(nanosecond)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

// AsyncLoopTask AsyncTimingTask 定时任务 参数含函数
func AsyncLoopTask(sleep time.Duration, fun func()) {
	go func() {
		for {
			fun()
			time.Sleep(sleep)
		}
	}()
}

// DeepCopyStruct 深拷贝函数
func DeepCopyStruct(original interface{}) interface{} {
	// 使用 JSON 序列化和反序列化来实现深拷贝
	data, err := json.Marshal(original)
	if err != nil {
		return nil
	}
	target := reflect.New(reflect.TypeOf(original).Elem()).Interface()
	err = json.Unmarshal(data, target)
	return target
}

func RandomHexadecimalString() string {
	rng := GetRand()
	const charset = "0123456789abcdef"
	const length = 16 // The length of the string you want to generate
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

func IsStrInArray(str string, strS []string) bool {
	// 如果 strS 为空，直接返回 true
	if len(strS) == 0 {
		return true
	}
	for _, v := range strS {
		if v == str {
			return true
		}
	}
	return false
}

// IsFileExist 判断文件是否存在
func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}

// ReadFile 读取文件
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	all, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return all, nil
}

// GetStrPreOrSuf 获取 指定字符串 前或后的字符串 0-原 1-前 1-后，没有返回原
func GetStrPreOrSuf(str string, sep string, t int) string {
	switch t {
	case 0:
		return str
	case -1:
		index := strings.Index(str, sep)
		if index != -1 {
			before := str[:index]
			return before
		}
		return str
	case 1:
		index := strings.Index(str, sep)
		if index != -1 {
			after := str[index+len(sep):]
			return after
		}
		return str
	default:
		return str
	}
}

func FakeUaAgent() string {
	var userAgent = []string{
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 7.1.1; CPH1723 Build/N6F26Q) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1; A37f Build/LMY47V) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.93 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 7.0; SM-G610M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.80 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; SM-G532G Build/MMB29T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.83 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.101 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.81 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 8.0.0; WAS-LX3 Build/HUAWEIWAS-LX3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.10 (KHTML, like Gecko) Chrome/8.0.552.237 Safari/534.10",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.117 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.10 (KHTML, like Gecko) Chrome/8.0.552.237 Safari/534.10",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.167 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 7.1.2; Redmi Note 5A Build/N2G47H; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/63.0.3239.111 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1.1; A37fw Build/LMY47V) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.84 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.67 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1; A1601 Build/LMY47I) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.98 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 8.0.0; FIG-LX3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.80 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0; vivo 1610 Build/MMB29M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.124 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; SM-G532G Build/MMB29T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.83 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 4.4.2; SM-G7102 Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.84 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 7.0; SM-G570M Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; SM-G532M Build/MMB29T; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; SM-J700M Build/MMB29K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36",
		"Mozilla/5.0 (X11; Datanyze; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.10 (KHTML, like Gecko) Chrome/8.0.552.237 Safari/534.10",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 7.1.1; CPH1723 Build/N6F26Q) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36",
		"Mozilla/5.0 (X11; Linux i686) AppleWebKit/534.30 (KHTML, like Gecko) Slackware/Chrome/12.0.742.100 Safari/534.30",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; de-DE) AppleWebKit/534.17 (KHTML, like Gecko) Chrome/10.0.649.0 Safari/534.17",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.63 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.90 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.64 Safari/537.31",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; CPH1607 Build/MMB29M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/63.0.3239.111 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36",
	}
	rng := GetRand()
	randomIndex := rng.Intn(len(userAgent) - 1)
	return userAgent[randomIndex]
}
