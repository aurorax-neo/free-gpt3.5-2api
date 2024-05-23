package RequestClient

import (
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/bogdanfinn/tls-client/profiles"
	"math/rand"
	"time"
)

const maxForceLogin = 5

var (
	ClientProfile = randomClientProfile()
	Ua            = browser.Random()
	MaxForceLogin = maxForceLogin
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
	MaxForceLogin--
}

func GetClientProfile() profiles.ClientProfile {
	if MaxForceLogin > 0 {
		return ClientProfile
	}
	ClientProfile = randomClientProfile()
	Ua = browser.Random()
	MaxForceLogin = maxForceLogin
	return ClientProfile
}
