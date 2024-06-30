package common

import (
	"free-gpt3.5-2api/constant"
	"github.com/bogdanfinn/tls-client/profiles"
	"math/rand"
	"time"
)

const retry = 5

var (
	clientProfile   = getRandomClientProfile()
	ua              = FakeUaAgent()
	updateThreshold = retry
)

func SubUpdateThreshold() {
	updateThreshold--
}
func getRandomClientProfile() profiles.ClientProfile {
	// 初始化随机数生成器
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	clientProfiles := []profiles.ClientProfile{
		profiles.Chrome_110,
		profiles.Okhttp4Android13,
		profiles.Opera_90,
	}
	// 随机选择一个
	randomIndex := rng.Intn(len(clientProfiles))
	return clientProfiles[randomIndex]
}

func GetClientProfile() profiles.ClientProfile {
	if updateThreshold < 0 {
		clientProfile = getRandomClientProfile()
		updateThreshold = retry
	}
	return clientProfile
}

func GetUa() string {
	if updateThreshold < 0 {
		ua = FakeUaAgent()
		updateThreshold = constant.ReTry
	}
	return ua
}
