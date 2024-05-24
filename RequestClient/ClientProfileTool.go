package RequestClient

import (
	"free-gpt3.5-2api/constant"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/bogdanfinn/tls-client/profiles"
	"math/rand"
	"time"
)

var (
	clientProfile = randomClientProfile()
	ua            = browser.Random()
	maxForceLogin = constant.ReTry
)

func randomClientProfile() profiles.ClientProfile {
	// 初始化随机数生成器
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	clientProfiles := []profiles.ClientProfile{
		profiles.Chrome_110,
		profiles.Okhttp4Android13,
		profiles.CloudflareCustom,
		profiles.Opera_90,
	}
	// 随机选择一个
	randomIndex := rng.Intn(len(clientProfiles))
	return clientProfiles[randomIndex]
}

func SubMaxForceLogin() {
	maxForceLogin--
	if maxForceLogin < 0 {
		clientProfile = randomClientProfile()
		ua = browser.Random()
		maxForceLogin = constant.ReTry
	}
}

func GetClientProfile() profiles.ClientProfile {
	return clientProfile
}

func GetUa() string {
	return ua
}
