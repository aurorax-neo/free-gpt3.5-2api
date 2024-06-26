package AccessTokenPool

import (
	"free-gpt3.5-2api/common"
	"sync"
)

var (
	instance *AccessTokenPool
	once     sync.Once
)

const AccAuthAuthorizationPre = "Bearer ac"

type AccessTokenPool struct {
	AccessTokens []*AccessToken
	index        int
}

type AccessToken struct {
	Token     string `yaml:"token,omitempty"`
	ExpiresAt int64  `yaml:"expires_at,omitempty"`
	CanUseAt  int64  `yaml:"-"`
}

func newAccAuthPool() *AccessTokenPool {
	return &AccessTokenPool{
		AccessTokens: make([]*AccessToken, 0),
		index:        0,
	}
}

func GetAccAuthPoolInstance() *AccessTokenPool {
	once.Do(func() {
		instance = newAccAuthPool()
	})
	return instance
}

func (a *AccessTokenPool) AddAccessToken(accessToken *AccessToken) {
	a.AccessTokens = append(a.AccessTokens, accessToken)
}

func (a *AccessTokenPool) AppendAccessTokens(accessTokens []*AccessToken) {
	a.AccessTokens = append(a.AccessTokens, accessTokens...)
}

func (a *AccessTokenPool) Size() int {
	return len(a.AccessTokens)
}

func (a *AccessTokenPool) IsEmpty() bool {
	return a.Size() == 0
}

func (a *AccessTokenPool) CanUseSize() int {
	var count int
	for _, v := range a.AccessTokens {
		if v.CanUseAt <= common.GetTimestampSecond(0) && v.ExpiresAt > common.GetTimestampSecond(0) {
			count++
		}
	}
	return count
}

func (a *AccessTokenPool) GetToken() string {
	if a.IsEmpty() || a.CanUseSize() == 0 {
		return ""
	}
	a.index = (a.index + 1) % len(a.AccessTokens)
	accessToken := a.AccessTokens[a.index]
	if accessToken.CanUseAt > common.GetTimestampSecond(0) || accessToken.ExpiresAt < common.GetTimestampSecond(0) {
		return a.GetToken()
	}
	return accessToken.Token
}

func (a *AccessTokenPool) SetCanUseAt(token string, canUseAt int64) {
	for _, v := range a.AccessTokens {
		if v.Token == token {
			v.CanUseAt = canUseAt
			break
		}
	}
}
